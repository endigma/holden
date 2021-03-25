package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	termmd "github.com/MichaelMure/go-term-markdown"
	log "github.com/sirupsen/logrus"

	"git.cya.cx/endigma/holden/structure"
	"git.cya.cx/endigma/holden/utils"
)

// Handler is the main response handler for requests
func Handler(rw http.ResponseWriter, req *http.Request) {
	if structure.Conf.General.Debug {
		utils.DebugReq(req)
	}

	newPath := strings.TrimPrefix(req.URL.Path, structure.Conf.General.Prefix)

	log.Debugf("Request for markdown file on %s", structure.Conf.General.Root+newPath)

	var path string
	if strings.HasSuffix(newPath, "/") {
		log.Debug("Detected trailing / on %s, checking for index", newPath)
		path = newPath + "_index.md"
	} else if strings.HasSuffix(newPath, ".md") {
		log.Debug("Detected trailing .md on %s, sending raw markdown", newPath)
		if utils.FileExists(structure.Conf.General.Root + newPath) {
			sourcefile, err := os.Open(structure.Conf.General.Root + newPath)
			utils.CheckErr(err)
			defer sourcefile.Close()

			source, err := ioutil.ReadAll(sourcefile)

			if structure.Conf.General.FancyCurl && !utils.IsInArr("true", req.Header.Values("RawPlease")) {
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

	page := render(structure.Conf.General.Root + path)

	if utils.FileExists(structure.Conf.General.Root + "/_sidebar.md") {
		log.Debug("_sidebar found")
		page.SidebarContents = render(structure.Conf.General.Root + "/_sidebar.md").Contents
	} else {
		page.SidebarContents = fmt.Sprintf("<h3><a class='home' href='%s/'><i class='fas fa-home'></i><span>%s</span></a></h3><ul>", structure.Conf.General.Prefix, structure.Conf.Website.SiteName) + renderSidebar(enumerateDir(structure.Conf.General.Root), "/") + "</ul>"
	}

	page.Raw = path

	tmpl := template.Must(template.ParseFiles(structure.Conf.General.WorkDir + "assets/static/page.html"))
	tmpl.Execute(rw, page)
}
