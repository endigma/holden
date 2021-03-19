package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	termmd "github.com/MichaelMure/go-term-markdown"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/naoina/toml"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

var conf Config

var md goldmark.Markdown

type Config struct {
	General struct {
		Root          string `toml:"root"`
		Port          string `toml:"port"`
		Debug         bool   `toml:"debug"`
		LinkableLines bool   `toml:"linkablelines"`
	} `toml:"general"`
	Aesthetic struct {
		HighlightStyle     string `toml:"highlightstyle"`
		LineNumbers        bool   `toml:"linenumbers"`
		LineNumbersInTable bool   `toml:"linenumbersintable"`
		TabWidth           int    `toml:"tabwidth"`
		UseClasses         bool   `toml:"useclasses"`
	} `toml:"aesthetic"`
}

type Page struct {
	Path            string
	Contents        string
	Meta            map[string]interface{}
	SidebarContents string
	Raw             string
}

type Directory struct {
	Directories []Directory
	Files       []string
	Name        string
}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func isInArr(s string, a []string) bool {
	for _, val := range a {
		if s == val {
			return true
		}
	}
	return false
}

func debugReq(req *http.Request) {
	log.WithFields(log.Fields{
		"agent":  req.UserAgent(),
		"method": req.Method,
		"addr":   exposeIP(req),
	}).Debug("Request")
}

func exposeIP(req *http.Request) string {
	addr := req.Header.Get("X-Real-Ip")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = req.RemoteAddr
	}
	return addr
}

func handler(rw http.ResponseWriter, req *http.Request) {
	if conf.General.Debug {
		debugReq(req)
	}

	log.Debugf("Request for markdown file on %s", conf.General.Root+req.URL.Path)

	var path string
	if strings.HasSuffix(req.URL.Path, "/") {
		log.Debug("Detected trailing / on %s, checking for index", req.URL.Path)
		path = req.URL.Path + "_index.md"
	} else if strings.HasSuffix(req.URL.Path, ".md") {
		log.Debug("Detected trailing .md on %s, sending raw markdown", req.URL.Path)
		if fileExists(conf.General.Root + req.URL.Path) {
			sourcefile, err := os.Open(conf.General.Root + req.URL.Path)
			checkErr(err)
			defer sourcefile.Close()

			source, err := ioutil.ReadAll(sourcefile)
			checkErr(err)

			if strings.HasPrefix(req.UserAgent(), "curl") {
				rw.Write(termmd.Render(string(source), 80, 4))
			} else {
				rw.Write(source)
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
		return
	} else {
		path = req.URL.Path + ".md"
	}

	page := render(conf.General.Root + path)

	if fileExists(conf.General.Root + "/_sidebar.md") {
		log.Debug("_sidebar found")
		page.SidebarContents = render(conf.General.Root + "/_sidebar.md").Contents
	} else {
		page.SidebarContents = "<h3><a class='home' href='/'><i class='fas fa-home'></i></a></h3><ul>" + renderSidebar(enumerateDir(conf.General.Root), "/") + "</ul>"
	}

	page.Raw = path

	tmpl := template.Must(template.ParseFiles("assets/static/page.html"))
	tmpl.Execute(rw, page)
}

func render(reqPath string) Page {
	if !fileExists(reqPath) {
		log.Errorf("No file found at: %s", reqPath)
		if fileExists(conf.General.Root + "/_404.md") {
			return render(conf.General.Root + "/_404.md")
		} else {
			return Page{
				Path:            "404",
				Contents:        "<p>404</p>",
				SidebarContents: "You shouldn't be seeing this!",
			}
		}

	}

	sourcefile, err := os.Open(reqPath)
	checkErr(err)
	defer sourcefile.Close()

	source, err := ioutil.ReadAll(sourcefile)
	checkErr(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	var page Page = Page{
		Path:            reqPath,
		Contents:        buf.String(),
		Meta:            meta.Get(context),
		SidebarContents: "You shouldn't be seeing this!",
	}

	return page
}

func renderSidebar(dirInfo Directory, prefix string) string {
	sidebarContent := strings.Builder{}
	if len(dirInfo.Directories) > 0 {
		for _, directory := range dirInfo.Directories {
			if isInArr("_index.md", directory.Files) {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder-plus'></i> <a href=\"" + prefix + directory.Name + "/" + "\">" + directory.Name + "</a></li>")
			} else {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder'></i> " + directory.Name + "</li>")
			}
			sidebarContent.WriteString("<ul>" + renderSidebar(directory, prefix+directory.Name+"/") + "</ul>")
		}
	}
	if len(dirInfo.Files) > 0 {
		for _, file := range dirInfo.Files {
			if !strings.HasPrefix(file, "_") {
				sidebarContent.WriteString("<li class=\"file\"><i class='fas fa-file-alt'></i> <a href=\"" + prefix + strings.TrimSuffix(file, ".md") + "\">" + strings.TrimSuffix(file, ".md") + "</a></li>")
			}
		}
	}
	result := sidebarContent.String()
	return result
}

func enumerateDir(path string) Directory {
	dirInfo := new(Directory)
	dircontents, err := ioutil.ReadDir(path)
	checkErr(err)

	currentDirInfo, err := os.Stat(path)
	checkErr(err)

	dirInfo.Name = currentDirInfo.Name()

	for _, file := range dircontents {
		if !strings.HasPrefix(file.Name(), "_") || file.Name() == "_index.md" {
			if file.IsDir() {
				dirInfo.Directories = append(dirInfo.Directories, enumerateDir(path+"/"+file.Name()))
			} else {
				dirInfo.Files = append(dirInfo.Files, file.Name())
			}
		}
	}

	sort.Strings(dirInfo.Files)
	sort.Slice(dirInfo.Directories, func(i int, j int) bool {
		return []byte(dirInfo.Directories[i].Name)[0] < []byte(dirInfo.Directories[j].Name)[0]
	})

	return *dirInfo
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	logfmt := new(log.TextFormatter)
	logfmt.TimestampFormat = "2006-01-02 15:04:05"
	logfmt.FullTimestamp = true
	log.SetFormatter(logfmt)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if conf.General.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(os.Args) == 1 {
		log.Fatal("Please provide a valid config file")
	}

	if !fileExists(os.Args[1]) {
		log.Fatal("Please provide a valid config file")
	}

	configFile, err := os.Open(os.Args[1])
	checkErr(err)

	defer configFile.Close()

	err = toml.NewDecoder(configFile).Decode(&conf)
	checkErr(err)

	md = goldmark.New(
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle(conf.Aesthetic.HighlightStyle),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(conf.Aesthetic.LineNumbers),
					html.TabWidth(conf.Aesthetic.TabWidth),
					html.LineNumbersInTable(conf.Aesthetic.LineNumbersInTable),
					html.WithClasses(conf.Aesthetic.UseClasses),
					html.LinkableLineNumbers(conf.General.LinkableLines, "l"),
				),
			),
		),
	)

	log.Info("Initialized")
	log.Debug("Debugging enabled")
}

func main() {
	log.Infof("Serving markdown files in %s", conf.General.Root)

	fs := http.FileServer(http.Dir("assets/serve"))
	http.Handle("/serve/", http.StripPrefix("/serve/", fs))

	fs2 := http.FileServer(http.Dir("assets/public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs2))
	http.Handle("/favicon.ico", http.StripPrefix("/public/", fs2))

	http.HandleFunc("/", handler)

	log.Infof("Starting http server on :%s", conf.General.Port)
	log.Fatal(http.ListenAndServe(":"+conf.General.Port, nil))
}
