package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	maxSize int

	server  *echo.Echo
	address string

	uc Usecase
}

func NewServer(ip string, port int, maxSize int, uc Usecase) *Server {
	api := Server{
		maxSize: maxSize,
		uc:      uc,
	}

	api.server = echo.New()

	api.server.POST("/api/register", api.PostCreateUser)
	api.server.POST("/api/login", api.PostLogin)
	api.server.GET("/api/profile/settings", api.GetSettings)
	api.server.POST("/api/profile/settings", api.PostSettings)
	api.server.GET("/api/profile/tests/:id", api.GetTests)
	api.server.GET("/api/profile/tests", api.GetTests)
	api.server.DELETE("/api/profile/tests/:id", api.DeleteGroup)
	api.server.POST("/api/profile/code", api.GetTestGroupRez)
	api.server.POST("/api/profile/tests", api.PostTests)
	api.server.POST("/api/profile/logout", api.PostLogout)
	api.server.GET("/api/profile/results", api.GetResults)
	api.server.File("*/", "./build/index.html")
	api.server.Static("/static/*", "./build/static")
	api.address = fmt.Sprintf("%s:%d", ip, port)
	return &api
}

func (api *Server) Run() {
	api.server.Logger.Fatal(api.server.Start(api.address))
}
