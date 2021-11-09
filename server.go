package main

import (
	"log"
	"net/http"
	"web-widgets/kanban-go/data"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jinzhu/configor"
	"github.com/unrolled/render"
)

var format = render.New()

var cards *data.CardsManager
var bdata *data.DataManager
var columns *data.ColumnsManager
var rows *data.RowsManager

// Config is the structure that stores the settings for this backend app
var Config AppConfig

func main() {
	configor.New(&configor.Config{ENVPrefix: "APP", Silent: true}).Load(&Config, "config.yml")

	bdata = data.NewDataManager(Config.Server.URL, Config.BinaryData)
	columns = &data.ColumnsManager{}
	rows = &data.RowsManager{}

	data.Init(Config.DB)

	cards = &data.CardsManager{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if Config.Server.Cors {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(c.Handler)
	}

	initRoutes(r)

	log.Printf("Starting webserver at port " + Config.Server.Port)
	err := http.ListenAndServe(Config.Server.Port, r)
	if err != nil {
		log.Println(err.Error())
	}
}
