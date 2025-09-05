package log

import (
	"bytes"
	registry "distributed/registry"
	"fmt"
	stlog "log"
	"net/http"
)

func SetClientLogger(serviceUrl string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{serviceUrl: serviceUrl})
}

type clientLogger struct {
	serviceUrl string
}

func (c *clientLogger) Write(data []byte) (n int, err error) {
	b := bytes.NewBuffer(data)
	res, err := http.Post(c.serviceUrl+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("could not send log with response: %s", res.Status)
	}
	return len(data), nil
}
