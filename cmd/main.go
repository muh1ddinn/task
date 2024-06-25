package main

import (
	"context"
	"fmt"
	"task/api"
	"task/config"
	"task/pkg/logger"
	"task/service"
	"task/storage/postgres"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		fmt.Println("error while connecting db, err: ", err)
		return
	}
	defer store.CloseDB()

	service := service.New(store, log)

	c := api.New(service, log)

	fmt.Println("programm is running on localhost:9090...")
	c.Run(":9090")
}
