package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type serviceUpdateHandler struct{}

func (s *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var p patch
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Update received: %v", p)
	prov.Update(p)
	w.WriteHeader(http.StatusOK)
}

func RegisterService(mux *http.ServeMux, r Registration) error {
	serviceUpdateUrl, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	mux.Handle(serviceUpdateUrl.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesUrl, "application/json", buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not register service with response: %s", res.Status)
	}
	return nil
}

func DeregisterService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not deregister service with response: %s", res.Status)
	}
	return nil
}

type providers struct {
	mutex    *sync.RWMutex
	services map[ServiceName][]string
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, entry := range pat.Added {
		if _, ok := p.services[entry.Name]; !ok {
			p.services[entry.Name] = make([]string, 0)
		}
		p.services[entry.Name] = append(p.services[entry.Name], entry.Url)
	}

	for _, entry := range pat.Removed {
		if providerUrls, ok := p.services[entry.Name]; ok {
			for i := range providerUrls {
				if providerUrls[i] == entry.Url {
					p.services[entry.Name] = append(p.services[entry.Name][:i], p.services[entry.Name][i+1:]...)
				}
			}
		}
	}
}

// there is only one provider
func (p *providers) get(serviceName ServiceName) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	providerUrls, ok := p.services[serviceName]
	if !ok {
		return "", fmt.Errorf("service %s not found", serviceName)
	}
	return providerUrls[0], nil
}

func GetProvider(serviceName ServiceName) (string, error) {
	return prov.get(serviceName)
}

var prov = &providers{
	mutex:    new(sync.RWMutex),
	services: make(map[ServiceName][]string),
}
