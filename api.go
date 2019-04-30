package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunServer() {
	r := gin.Default()

	r.POST("/build_route", func(c *gin.Context) {
		routeReq := new(RouteRequirements)
		err := c.ShouldBindJSON(routeReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		route := buildRoute(routeReq)
		if route == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "nil"})
			return
		}
		jsn, err := json.Marshal(route)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, string(jsn))
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
