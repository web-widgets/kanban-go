package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func initRoutes(r chi.Router) {

	r.Get("/cards", func(w http.ResponseWriter, r *http.Request) {
		data, err := cards.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/cards", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormCard(w, r)
		if err == nil {
			id, err = cards.Add(info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Put("/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormCard(w, r)
		if err == nil {
			id = NumberParam(r, "id")
			err = cards.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := cards.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})

	r.Get("/data/{id}/{name}", func(w http.ResponseWriter, r *http.Request) {
		res, err := bdata.ToResponse(w, NumberParam(r, "id"))

		if err != nil {
			format.Text(w, 500, err.Error())
		} else if !res {
			format.Text(w, 500, "")
		}
	})

	r.Post("/data", func(w http.ResponseWriter, r *http.Request) {
		rec, err := bdata.FromRequest(r, "upload")
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, rec)
		}
	})

	r.Get("/columns", func(w http.ResponseWriter, r *http.Request) {
		data, err := stages.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/columns", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormStage(w, r)
		if err == nil {
			id, err = stages.Add(info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Put("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormStage(w, r)
		if err == nil {
			id = NumberParam(r, "id")
			err = stages.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := stages.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})
}
