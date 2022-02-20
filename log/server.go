package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type filelog string

func (fl filelog) Write(date []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if err != nil {
		return 0, err
	}

	defer f.Close()

	return f.Write(date)
}

func Run(destination string) {
	log = stlog.New(filelog(destination), "go: ", stlog.LstdFlags)

}

func RegisterHandlers() {
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(request.Body)
			if err != nil || len(msg) == 0 {
				writer.WriteHeader(http.StatusBadRequest)
			}
			write(string(msg))

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
