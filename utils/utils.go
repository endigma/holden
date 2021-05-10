package utils

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// FileExists checks if a file or folder
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// CheckErr checks and logs errors
func CheckErr(e error) {
	if e != nil {
		log.Fatal().Err(e).Msg(e.Error())
	}
}

// IsInArr checks if a string is part of a string slice
func IsInArr(s string, a []string) bool {
	for _, val := range a {
		if s == val {
			return true
		}
	}
	return false
}

// DebugReq prints extra debug info about a request.
func DebugReq(req *http.Request) {
	log.Debug().
		Str("agent", req.UserAgent()).
		Str("method", req.Method).
		Str("addr", exposeIP(req)).
		Msg("Request")
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
