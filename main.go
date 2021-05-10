package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	"gitcat.ca/endigma/holden/handler"
	"gitcat.ca/endigma/holden/utils"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if viper.GetBool("general.debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

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

	http.HandleFunc(viper.GetString("general.prefix")+"/", handler.Handler)

	log.Info().Msgf("Starting http server on :%s", viper.GetString("general.port"))
	log.Fatal().Err(http.ListenAndServe(":"+viper.GetString("general.port"), nil)).Msg("Fatal")
}
