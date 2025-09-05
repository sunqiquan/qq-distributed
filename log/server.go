package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
	"path/filepath"
)

var log *stlog.Logger
var logfd *os.File

func RunWithOption(dest string, toStderr bool) {
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("Create log dir failed: " + err.Error())
	}
	fd, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Open log file failed: " + err.Error())
	}

	logfd = fd
	var writer io.Writer = fd
	if toStderr {
		writer = io.MultiWriter(os.Stdout, logfd)
	}
	log = stlog.New(writer, "[detributed] - ", stlog.LstdFlags)
}

func Run(dest string) {
	RunWithOption(dest, false)
}

func RunWithWirteToStderr(dest string) {
	RunWithOption(dest, true)
}

func Close() {
	if logfd != nil {
		logfd.Close()
		logfd = nil
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func write(msg string) {
	log.Printf("%v\n", msg)
}
