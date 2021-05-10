package structure

import (
	"os"
	"path/filepath"
	"strings"

	"gitcat.ca/endigma/holden/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

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
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath("/etc/holden/")
	viper.AddConfigPath("$HOME/.holden")
	viper.AddConfigPath(".")

	viper.SetDefault("general.docroot", "docs")
	viper.SetDefault("general.port", "11011")
	viper.SetDefault("general.prefix", "")
	viper.SetDefault("general.workdir", "")
	viper.SetDefault("general.allowhtml", false)
	viper.SetDefault("general.linkablelines", true)
	viper.SetDefault("general.fancycurl", true)
	viper.SetDefault("general.debug", false)
	viper.SetDefault("website.sitename", "")
	viper.SetDefault("website.backtotop", true)
	viper.SetDefault("website.sidebar", true)
	viper.SetDefault("aesthetic.highlightstyle", "solarized-dark256")
	viper.SetDefault("aesthetic.linenumbers", true)
	viper.SetDefault("aesthetic.linenumbersintable", true)
	viper.SetDefault("aesthetic.tabwidth", 4)
	viper.SetDefault("aesthetic.useclasses", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal().Err(err).Str("cfgfile", viper.ConfigFileUsed()).Msg("Failure to load config")
		} else {
			log.Info().Msg("Creating a new config file!")
			_, err = os.Create("config.toml")
			utils.CheckErr(err)
			viper.WriteConfig()
		}
	} else {
		viper.WatchConfig()
	}

	viper.WriteConfig()

	ex, err := os.Executable()
	utils.CheckErr(err)

	binPath := filepath.Dir(ex) + "/"

	viper.Set("docroot", strings.ReplaceAll(viper.GetString("docroot"), "_binary", strings.TrimSuffix(binPath, "/")))
}
