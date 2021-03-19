package utils

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.Fatal(e)
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
	log.WithFields(log.Fields{
		"agent":  req.UserAgent(),
		"method": req.Method,
		"addr":   exposeIP(req),
	}).Debug("Request")
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
