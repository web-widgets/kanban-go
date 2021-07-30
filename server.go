package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jinzhu/configor"
	"github.com/unrolled/render"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var format = render.New()
var db *gorm.DB

var cards *CardsManager
var bdata *DataManager
var stages *StagesManager

// Config is the structure that stores the settings for this backend app
var Config AppConfig

func main() {
	configor.New(&configor.Config{ENVPrefix: "APP", Silent: true}).Load(&Config, "config.yml")

	var err error
	db, err = gorm.Open("sqlite3", Config.DB.Path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Stage{})
	db.AutoMigrate(&BinaryData{})
	if Config.DB.ResetOnStart {
		dataDown()
		dataUp()
	}

	cards = &CardsManager{}

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
	err = http.ListenAndServe(Config.Server.Port, r)
	if err != nil {
		log.Println(err.Error())
	}
}

func NumberParam(r *http.Request, key string) int {
	value := chi.URLParam(r, key)
	num, _ := strconv.Atoi(value)

	return num
}

func ParseFormCard(w http.ResponseWriter, r *http.Request) (CardUpdate, error) {
	c := CardUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseFormStage(w http.ResponseWriter, r *http.Request) (StageUpdate, error) {
	c := StageUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

type Response struct {
	ID int `json:"id"`
}
