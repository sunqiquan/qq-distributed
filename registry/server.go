package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ServerPort = "8090"
const ServicesUrl = "http://localhost:" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	err := r.sendRequiredServices(reg)
	if err != nil {
		return err
	}
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				Url:  reg.ServiceUrl,
			},
		},
	})
	return nil
}

func (r *registry) notify(fullPath patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{
					Added:   []patchEntry{},
					Removed: []patchEntry{},
				}
				sendUpdate := false
				for _, added := range fullPath.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPath.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}

				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateUrl)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}(reg)
	}
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, regService := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if regService.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: reg.ServiceName,
					Url:  reg.ServiceUrl,
				})
			}
		}
	}

	err := r.sendPatch(p, reg.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for i := range r.registrations {
		if r.registrations[i].ServiceUrl == url {
			r.mutex.Lock()
			defer r.mutex.Unlock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: r.registrations[i].ServiceName,
						Url:  r.registrations[i].ServiceUrl,
					},
				},
			})
			return nil
		}
	}
	return fmt.Errorf("service %s not found", url)
}

var reg = &registry{
	registrations: make([]Registration, 0),
	mutex:         &sync.RWMutex{},
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := json.Marshal(reg.registrations)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with url %s\n", r.ServiceName, r.ServiceUrl)
		reg.add(r)
	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Panicln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)

		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("removing service with url %s\n", url)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
