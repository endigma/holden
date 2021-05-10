package handler

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"gitcat.ca/endigma/holden/structure"
	"gitcat.ca/endigma/holden/utils"
	"github.com/alecthomas/chroma/formatters/html"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var md goldmark.Markdown
var metareader goldmark.Markdown

func render(reqPath string) structure.Page {
	if !utils.FileExists(reqPath) {
		log.Error().Msgf("No file found at: %s", reqPath)
		if utils.FileExists(viper.GetString("general.docroot") + "/_404.md") {
			return render(viper.GetString("general.docroot") + "/_404.md")
		} else {
			return structure.Page{
				Prefix:          viper.GetString("general.prefix"),
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
		Prefix:           viper.GetString("general.prefix"),
		Contents:         buf.String(),
		Meta:             meta.Get(context),
		SidebarContents:  "You shouldn't be seeing this!",
		DisplayBackToTop: viper.GetBool("website.backtotop"),
		DisplaySidebar:   viper.GetBool("website.sidebar"),
	}

	return page
}

func renderSidebar(dirInfo structure.Directory, prefix string) string {
	sidebarContent := strings.Builder{}
	if len(dirInfo.Directories) > 0 {
		for _, directory := range dirInfo.Directories {
			if utils.IsInArr("_index.md", directory.Files) {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder-plus'></i> <a href=\"" + viper.GetString("general.prefix") + "/" + directory.Name + "/" + "\">" + directory.Name + "</a></li>")
			} else {
				sidebarContent.WriteString("<li class=\"folder\"><i class='fas fa-folder'></i> " + directory.Name + "</li>")
			}
			sidebarContent.WriteString("<ul>" + renderSidebar(directory, prefix+directory.Name+"/") + "</ul>")
		}
	}
	if len(dirInfo.Files) > 0 {
		for _, file := range dirInfo.Files {
			if !strings.HasPrefix(file, "_") {
				filecontent, err := os.ReadFile(viper.GetString("general.docroot") + prefix + file)
				utils.CheckErr(err)

				context := parser.NewContext()
				if err := metareader.Convert(filecontent, ioutil.Discard, parser.WithContext(context)); err != nil {
					log.Panic().Err(err).Msg("Failure in sidebar generation")
				}

				fileMeta, err := meta.TryGet(context)
				utils.CheckErr(err)

				if fileMeta["Short"] != nil && fileMeta["Short"].(string) != "" {
					sidebarContent.WriteString("<li class=\"file\"><i class='fas fa-file-alt'></i> <a href=\"" + viper.GetString("general.prefix") + prefix + strings.TrimSuffix(file, ".md") + "\">" + fileMeta["Short"].(string) + "</a></li>")
				} else {
					sidebarContent.WriteString("<li class=\"file\"><i class='fas fa-file-alt'></i> <a href=\"" + viper.GetString("general.prefix") + prefix + strings.TrimSuffix(file, ".md") + "\">" + strings.TrimSuffix(file, ".md") + "</a></li>")
				}
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
				highlighting.WithStyle(viper.GetString("aesthetic.highlightstyle")),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(viper.GetBool("aesthetic.linenumbers")),
					html.TabWidth(viper.GetInt("aesthetic.tabwidth")),
					html.LineNumbersInTable(viper.GetBool("aesthetic.linenumbersintable")),
					html.WithClasses(viper.GetBool("aesthetic.tabwidth")),
					html.LinkableLineNumbers(viper.GetBool("general.linkablelines"), "l"),
				),
			),
		),
	)
	metareader = goldmark.New(goldmark.WithExtensions(meta.Meta))
	if viper.GetBool("general.allowhtml") {
		md.Renderer().AddOptions(goldmarkhtml.WithUnsafe())
	}
}
