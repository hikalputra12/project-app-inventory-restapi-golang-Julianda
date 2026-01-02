package main

import (
	"app-inventory/database"
	"app-inventory/handler"
	"app-inventory/repository"
	"app-inventory/router"
	"app-inventory/service"
	"app-inventory/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal("error file configration")
	}
	fmt.Println(config)
	db, err := database.InitDB(config.DB)
	if err != nil {
		panic(err)
	}
	logger, err := utils.InitLogger("./logs/app-", true)

	repo := repository.AllRepo(db, logger)
	service := service.AllService(repo, logger)
	handler := handler.AllHandler(service, logger)

	r := router.NewRouter(handler, service, logger)

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error server")
	}
}
