package main

import (
	"gochat/config"
	"gochat/db"
	"gochat/internal/room"
	"gochat/internal/user"
	"gochat/router"

	"log"

	"github.com/go-playground/validator/v10"
)

func main() {
	cfg := config.New()

	dbConn, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err)
	}
	defer dbConn.Close()

	val := validator.New()
	userHdl := user.Init(cfg, val, dbConn.GetDB())
	roomHdl := room.Init(cfg, val)

	r := router.InitRouter(cfg, userHdl, roomHdl)
	router.Start(r, cfg.ServerHost)
}
