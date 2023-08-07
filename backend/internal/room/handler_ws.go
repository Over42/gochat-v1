package room

import (
	"gochat/config"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service Service
	config  *config.Config
}

func NewHandler(svc Service, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		config:  cfg,
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.service.CreateRoom(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) DeleteRoom(c *gin.Context) {
	var req DeleteRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.DeleteRoom(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) GetRooms(c *gin.Context) {
	res, err := h.service.GetRooms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) JoinRoom(c *gin.Context) {
	// CSRF protection
	upgrader.CheckOrigin = func(_ *http.Request) bool {
		return c.Request.Header.Get("Origin") == h.config.OriginHost
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &JoinRoomReq{
		Conn:     conn,
		RoomID:   c.Param("roomId"),
		UserID:   c.Query("userId"),
		Username: c.Query("username"),
	}

	err = h.service.JoinRoom(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) GetClients(c *gin.Context) {
	req := GetClientsReq{
		RoomID: c.Param("roomId"),
	}

	res, err := h.service.GetClients(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
