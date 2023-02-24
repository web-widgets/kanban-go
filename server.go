package main

import (
	"context"
	"log"
	"net/http"
	"web-widgets/kanban-go/api"
	"web-widgets/kanban-go/data"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jinzhu/configor"
	"github.com/unrolled/render"
)

var format = render.New()

// Config is the structure that stores the settings for this backend app
var Config AppConfig

func main() {
	configor.New(&configor.Config{ENVPrefix: "APP", Silent: true}).Load(&Config, "config.yml")

	dao := data.NewDAO(Config.DB, Config.Server.URL, Config.BinaryData, Config.Features)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if len(Config.Server.Cors) > 0 {

		c := cors.New(cors.Options{
			AllowedOrigins:   Config.Server.Cors,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Remote-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(c.Handler)
	}

	// auth
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Remote-Token")
			if token == "" {
				if r.Method == http.MethodGet {
					token = r.URL.Query().Get("token")
				}
			}

			if token != "" {
				id, device, err := verifyUserToken([]byte(token))
				if err != nil {
					log.Println("[token]", err.Error())
				} else {
					r = r.WithContext(context.WithValue(context.WithValue(r.Context(), "user_id", id), "device_id", device))
				}
			}
			next.ServeHTTP(w, r)
		})
	})

	apiServer := api.BuildAPI(dao)
	r.Get("/api/v1", apiServer.ServeHTTP)
	r.Post("/api/v1", apiServer.ServeHTTP)

	initRoutes(r, dao, apiServer.Events)

	log.Printf("Starting webserver at port " + Config.Server.Port)
	err := http.ListenAndServe(Config.Server.Port, r)
	if err != nil {
		log.Println(err.Error())
	}
}
