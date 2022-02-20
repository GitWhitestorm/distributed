package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}
func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)

}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct {
}

func (s RegistryService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log.Println("Request received")

	switch request.Method {
	case http.MethodPost:
		dec := json.NewDecoder(request.Body)
		var r Registration
		err := dec.Decode(&r)

		if err != nil {
			log.Println("err")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Addming service: %v with URL:%s\n", r.ServiceName, r.ServiceURL)

		err = reg.add(r)

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		fmt.Println("url:", url)
		log.Printf("Removing service at URL:%s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
