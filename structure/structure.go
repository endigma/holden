package structure

import (
	"os"
	"path/filepath"
	"strings"

	"git.cya.cx/endigma/holden/utils"
	"github.com/naoina/toml"
	log "github.com/sirupsen/logrus"
)

// Conf is a variable where all the config options are loaded
var Conf Config

// Config is a struct that holds config info
type Config struct {
	General struct {
		Root          string `toml:"docroot"`
		Port          string `toml:"port"`
		Prefix        string `toml:"prefix"`
		WorkDir       string `toml:"workdir"`
		AllowHtml     bool   `toml:"allowhtml"`
		LinkableLines bool   `toml:"linkablelines"`
		FancyCurl     bool   `toml:"fancycurl"`
		Debug         bool   `toml:"debug"`
	} `toml:"general"`
	Website struct {
		SiteName         string `toml:"sitename"`
		DisplayBackToTop bool   `toml:"backtotop"`
		DisplaySidebar   bool   `toml:"sidebar"`
	} `toml:"website"`
	Aesthetic struct {
		HighlightStyle     string `toml:"highlightstyle"`
		TabWidth           int    `toml:"tabwidth"`
		LineNumbers        bool   `toml:"linenumbers"`
		LineNumbersInTable bool   `toml:"linenumbersintable"`
		UseClasses         bool   `toml:"useclasses"`
	} `toml:"aesthetic"`
}

// Page is a struct that holds webpage info for the template
type Page struct {
	Prefix           string
	Contents         string
	Meta             map[string]interface{}
	SidebarContents  string
	Raw              string
	DisplayBackToTop bool
	DisplaySidebar   bool
}

// Directory is a struct that holds information
// about the directory that holds the served markdown files
type Directory struct {
	Directories []Directory
	Files       []string
	Name        string
}

func init() {
	var configPath string
	ex, err := os.Executable()
	utils.CheckErr(err)

	binPath := filepath.Dir(ex) + "/"

	if len(os.Args) == 1 { // no arguments
		if utils.FileExists(binPath + "config.toml") {
			configPath = binPath + "config.toml"
		} else {
			log.Fatal("Please provide a valid config file or put config.toml in " + binPath)
		}
	} else { // there are arguments
		if utils.FileExists(os.Args[1]) {
			configPath = os.Args[1]
		} else {
			log.Fatal("Please provide a valid config file or put config.toml in " + binPath)
		}
	}

	configFile, err := os.Open(configPath)
	utils.CheckErr(err)
	defer configFile.Close()
	err = toml.NewDecoder(configFile).Decode(&Conf)
	utils.CheckErr(err)

	Conf.General.Root = strings.ReplaceAll(Conf.General.Root, "_binary", strings.TrimSuffix(binPath, "/"))

	log.Debug("Structure initialized")
}
