package main

import (
	"context"
	"distributed/log"
	"distributed/registry"
	"distributed/service"
	"fmt"
)

func main() {
	log.Run("./logs/distributed.log")
	defer log.Close()

	host := "localhost"
	port := "8091"
	reg := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceUrl:       "http://" + host + ":" + port,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateUrl: "http://" + host + ":" + port + "/services",
	}

	ctx, err := service.Start(context.Background(), port, reg, log.RegisterHandlers)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
	fmt.Printf("%s stopped.\n", reg.ServiceName)
}
