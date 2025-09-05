package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"distributed/student"
	"fmt"
)

func main() {
	host := "localhost"
	port := "8092"
	reg := registry.Registration{
		ServiceName: registry.StudentService,
		ServiceUrl:  "http://" + host + ":" + port,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
		},
		ServiceUpdateUrl: "http://" + host + ":" + port + "/services",
	}

	ctx, err := service.Start(context.Background(), port, reg, student.RegisterHandlers)
	if err != nil {
		panic(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, reg.ServiceName)
	}

	<-ctx.Done()
	fmt.Printf("%s stopped.\n", reg.ServiceName)
}
