package main

import (
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"git.cya.cx/endigma/holden/handler"
	"git.cya.cx/endigma/holden/structure"
	"git.cya.cx/endigma/holden/utils"
)

func init() {
	logfmt := new(log.TextFormatter)
	logfmt.TimestampFormat = "2006-01-02 15:04:05"
	logfmt.FullTimestamp = true
	log.SetFormatter(logfmt)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if structure.Conf.General.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Initialized")
	log.Debug("Debugging enabled")
}

func main() {
	log.Infof("Serving markdown files in %s", structure.Conf.General.Root)

	if structure.Conf.General.WorkDir == "_binary" {
		ex, err := os.Executable()
		utils.CheckErr(err)

		structure.Conf.General.WorkDir = filepath.Dir(ex) + "/"
	}

	log.Debug(structure.Conf.General.WorkDir)

	fs := http.FileServer(http.Dir(structure.Conf.General.WorkDir + "assets/static"))
	http.Handle(structure.Conf.General.Prefix+"/static/", http.StripPrefix(structure.Conf.General.Prefix+"/static/", fs))

	fs2 := http.FileServer(http.Dir(structure.Conf.General.WorkDir + "assets/public"))
	http.Handle(structure.Conf.General.Prefix+"/public/", http.StripPrefix(structure.Conf.General.Prefix+"/public/", fs2))
	http.Handle(structure.Conf.General.Prefix+"/favicon.ico", http.StripPrefix(structure.Conf.General.Prefix+"/public/", fs2))

	http.HandleFunc(structure.Conf.General.Prefix+"/", handler.Handler)

	log.Infof("Starting http server on :%s", structure.Conf.General.Port)
	log.Fatal(http.ListenAndServe(":"+structure.Conf.General.Port, nil))
}
