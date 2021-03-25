package handler

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"git.cya.cx/endigma/holden/structure"
	"git.cya.cx/endigma/holden/utils"
	"github.com/alecthomas/chroma/formatters/html"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var md goldmark.Markdown

func render(reqPath string) structure.Page {
	if !utils.FileExists(reqPath) {
		log.Errorf("No file found at: %s", reqPath)
		if utils.FileExists(structure.Conf.General.Root + "/_404.md") {
			return render(structure.Conf.General.Root + "/_404.md")
		} else {
			return structure.Page{
				Prefix:          structure.Conf.General.Prefix,
				Raw:             "_404",
				Contents:        "<p>404</p>",
				SidebarContents: "You shouldn't be seeing this!",
			}
		}

	}

	sourcefile, err := os.Open(reqPath)
	utils.CheckErr(err)
	defer sourcefile.Close()

	source, err := ioutil.ReadAll(sourcefile)
	utils.CheckErr(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	var page structure.Page = structure.Page{
		Prefix:           structure.Conf.General.Prefix,
		Contents:         buf.String(),
		Meta:             meta.Get(context),
		SidebarContents:  "You shouldn't be seeing this!",
		DisplayBackToTop: structure.Conf.Website.DisplayBackToTop,
	}

	return page
}

func renderSidebar(dirInfo structure.Directory, prefix string) string {
	sidebarContent := strings.Builder{}
	if len(dirInfo.Directories) > 0 {
		for _, directory := range dirInfo.Directories {
			if utils.IsInArr("_index.md", directory.Files) {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder-plus'></i> <a href=\"" + structure.Conf.General.Prefix + "/" + directory.Name + "/" + "\">" + directory.Name + "</a></li>")
			} else {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder'></i> " + directory.Name + "</li>")
			}
			sidebarContent.WriteString("<ul>" + renderSidebar(directory, prefix+directory.Name+"/") + "</ul>")
		}
	}
	if len(dirInfo.Files) > 0 {
		for _, file := range dirInfo.Files {
			if !strings.HasPrefix(file, "_") {
				sidebarContent.WriteString("<li class=\"file\"><i class='fas fa-file-alt'></i> <a href=\"" + structure.Conf.General.Prefix + prefix + strings.TrimSuffix(file, ".md") + "\">" + strings.TrimSuffix(file, ".md") + "</a></li>")
			}
		}
	}
	result := sidebarContent.String()
	return result
}

func enumerateDir(path string) structure.Directory {
	dirInfo := new(structure.Directory)
	dircontents, err := ioutil.ReadDir(path)
	utils.CheckErr(err)

	currentDirInfo, err := os.Stat(path)
	utils.CheckErr(err)

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

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle(structure.Conf.Aesthetic.HighlightStyle),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(structure.Conf.Aesthetic.LineNumbers),
					html.TabWidth(structure.Conf.Aesthetic.TabWidth),
					html.LineNumbersInTable(structure.Conf.Aesthetic.LineNumbersInTable),
					html.WithClasses(structure.Conf.Aesthetic.UseClasses),
					html.LinkableLineNumbers(structure.Conf.General.LinkableLines, "l"),
				),
			),
		),
	)
	if structure.Conf.General.AllowHtml {
		md.Renderer().AddOptions(goldmarkhtml.WithUnsafe())
	}
}
