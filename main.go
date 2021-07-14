package main

import (
	"encoding/json"
	"howanski/transmission-gonductor/backend"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ConfigStorage struct {
	gorm.Model
	ConfigKey   string
	ConfigValue string
}

var router *gin.Engine

func main() {

	db, err := gorm.Open(sqlite.Open("database/gonductor.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to open database in database/gonductor.sqlite")
	}

	db.AutoMigrate(&ConfigStorage{})

	router = gin.Default()

	router.LoadHTMLGlob("front-end/html/*")

	router.Static("/assets/images", "./front-end/images")
	router.Static("/assets/js", "./front-end/js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "GOnductor",
			},
		)
	})

	router.GET("/gonductor-stats", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"connectionStatus": "Ruined",
				"lastPing":         "Never",
			},
		)
	})

	router.GET("/settings", func(c *gin.Context) {
		var configs []ConfigStorage
		db.Find(&configs)
		c.JSON(
			http.StatusOK,
			configs,
		)
	})

	router.POST("/settings", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)
		var f interface{}
		err := json.Unmarshal(body, &f)
		if err == nil {
			mapmap := f.(map[string]interface{})
			for dbKey, valInterface := range mapmap {
				if dbKey != "" {
					if valInterface != "" {
						value := valInterface.(string)
						var dbConfig ConfigStorage
						db.First(&dbConfig, "config_key = ?", dbKey)
						if dbConfig.ConfigKey != "" {
							db.Model(&dbConfig).Update("config_value", value)
						} else {
							db.Create(&ConfigStorage{ConfigKey: dbKey, ConfigValue: value})
						}
					}
				}
			}
		}
	})

	go router.Run()
	backend.TransmissionDaemon()
}
