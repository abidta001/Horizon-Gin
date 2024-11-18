package main

import (
	db "admin/DB"
	"admin/route"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	db.InitDatabase()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	router := gin.Default()
	route.RegisterURL(router)
	router.Run(":3000")
}
