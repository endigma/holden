package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	"gitcat.ca/endigma/holden/handler"
	"gitcat.ca/endigma/holden/utils"
	"github.com/NYTimes/gziphandler"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.SetEnvPrefix("HLDN_")
	viper.AutomaticEnv()

	viper.AddConfigPath("/etc/holden/")
	viper.AddConfigPath("$HOME/.holden")
	viper.AddConfigPath(".")

	viper.SetDefault("general.docroot", "docs")
	viper.SetDefault("general.port", "11011")
	viper.SetDefault("general.prefix", "")
	viper.SetDefault("general.workdir", "")
	viper.SetDefault("general.allowhtml", false)
	viper.SetDefault("general.linkablelines", true)
	viper.SetDefault("general.fancystdout", true)
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

	if viper.GetBool("general.fancystdout") {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.With().Caller().Logger()
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if viper.GetBool("general.debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	handler.InitializeRenderer()

	log.Info().Msg("Initialized")
	log.Debug().Msg("Debugging enabled")
}

func main() {
	log.Info().Msgf("Serving markdown files in %s", viper.GetString("general.docroot"))

	if viper.GetString("general.workdir") == "_binary" {
		ex, err := os.Executable()
		utils.CheckErr(err)

		viper.Set("general.workdir", filepath.Dir(ex)+"/")
	}

	log.Debug().Msg(viper.GetString("general.workdir"))

	fs := http.FileServer(http.Dir(viper.GetString("general.workdir") + "assets/static"))
	http.Handle(viper.GetString("general.prefix")+"/static/", http.StripPrefix(viper.GetString("general.prefix")+"/static/", fs))

	fs2 := http.FileServer(http.Dir(viper.GetString("general.workdir") + "assets/public"))
	http.Handle(viper.GetString("general.prefix")+"/public/", http.StripPrefix(viper.GetString("general.prefix")+"/public/", fs2))
	http.Handle(viper.GetString("general.prefix")+"/favicon.ico", http.StripPrefix(viper.GetString("general.prefix")+"/public/", fs2))
	http.Handle(viper.GetString("general.prefix")+"/", gziphandler.GzipHandler(handler.Handler))

	log.Info().Msgf("Starting http server on :%s", viper.GetString("general.port"))
	log.Fatal().Err(http.ListenAndServe(":"+viper.GetString("general.port"), nil)).Msg("Fatal")
}
