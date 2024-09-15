package main

import (
	"gochatv1/config"
	"gochatv1/db"
	"gochatv1/internal/room"
	"gochatv1/internal/user"
	"gochatv1/router"

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
