package router

import (
	"gochat/config"
	"gochat/internal/room"
	"gochat/internal/user"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(cfg *config.Config, userHandler *user.Handler, roomHandler *room.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.OriginHost},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == cfg.OriginHost
		},
		MaxAge: 24 * time.Hour,
	}))

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.POST("/rooms", roomHandler.CreateRoom)
	r.DELETE("/rooms", roomHandler.DeleteRoom)
	r.GET("/rooms", roomHandler.GetRooms)
	r.GET("/rooms/:roomId", roomHandler.JoinRoom)
	r.GET("/rooms/:roomId/clients", roomHandler.GetClients)

	return r
}

func Start(r *gin.Engine, addr string) {
	r.Run(addr)
}
