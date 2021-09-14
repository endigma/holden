package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	termmd "github.com/MichaelMure/go-term-markdown"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"gitcat.ca/endigma/holden/utils"
)

// Handler is the main response handler for requests
var Handler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

	if viper.GetBool("general.debug") {
		utils.DebugReq(req)
	}

	newPath := strings.TrimPrefix(req.URL.Path, viper.GetString("general.prefix"))

	log.Debug().Msgf("Request for markdown file on %s", viper.GetString("general.docroot")+newPath)

	var path string
	if strings.HasSuffix(newPath, "/") {
		log.Debug().Msgf("Detected trailing / on %s, checking for index", newPath)
		path = newPath + "_index.md"
	} else if strings.HasSuffix(newPath, ".md") {
		log.Debug().Msgf("Detected trailing .md on %s, sending raw markdown", newPath)
		if utils.FileExists(viper.GetString("general.docroot") + newPath) {
			sourcefile, err := os.Open(viper.GetString("general.docroot") + newPath)
			utils.CheckErr(err)
			defer sourcefile.Close()

			source, err := ioutil.ReadAll(sourcefile)
			if err != nil {
				log.Error().Err(err).Msg("Failed to read file")
			}

			if viper.GetBool("general.fancycurl") && !utils.IsInArr("true", req.Header.Values("RawPlease")) {
				if strings.HasPrefix(req.UserAgent(), "curl") {
					rw.Write(termmd.Render(string(source), 80, 4))
				} else {
					rw.Write(source)
				}
			} else {
				rw.Write(source)
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
		return
	} else {
		path = newPath + ".md"
	}

	page := render(viper.GetString("general.docroot") + path)

	if utils.FileExists(viper.GetString("general.docroot") + "/_sidebar.md") {
		log.Debug().Msg("_sidebar found")
		page.SidebarContents = render(viper.GetString("general.docroot") + "/_sidebar.md").Contents
	} else {
		sidebarContent := renderSidebar(enumerateDir(viper.GetString("general.docroot")), "/")
		if sidebarContent == "" {
			page.DisplaySidebar = false
		}
		page.SidebarContents = fmt.Sprintf("<h3><a class='home' href='%s/'><i class='fas fa-home'></i><span>%s</span></a></h3><ul>", viper.GetString("general.prefix"), viper.GetString("website.sitename")) + sidebarContent + "</ul>"
	}

	cacheSince := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(0, 0, 1).Format(http.TimeFormat)

	rw.Header().Set("Cache-Control", "max-age:290304000, public")
	rw.Header().Set("Last-Modified", cacheSince)
	rw.Header().Set("Expires", cacheUntil)

	page.Raw = path

	tmpl := template.Must(template.ParseFiles(viper.GetString("workdir") + "assets/page.html"))
	tmpl.Execute(rw, page)
})
