package main

import (
	"net/http"

	"github.com/go-chi/chi"
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

	r.Put("/cards/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormMoveCard(w, r)
		if err == nil {
			id = NumberParam(r, "id")
			err = cards.Move(id, info.Card, int(info.Before))
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

	r.Get("/uploads/{id}/{name}", func(w http.ResponseWriter, r *http.Request) {
		res, err := bdata.ToResponse(w, NumberParam(r, "id"))

		if err != nil {
			format.Text(w, 500, err.Error())
		} else if !res {
			format.Text(w, 500, "")
		}
	})

	r.Post("/uploads", func(w http.ResponseWriter, r *http.Request) {
		rec, err := bdata.FromRequest(r, "upload")
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, rec)
		}
	})

	r.Get("/columns", func(w http.ResponseWriter, r *http.Request) {
		data, err := columns.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/columns", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormColumn(w, r)
		if err == nil {
			id, err = columns.Add(info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Put("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormColumn(w, r)
		if err == nil {
			id = NumberParam(r, "id")
			err = columns.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := columns.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})

	r.Get("/rows", func(w http.ResponseWriter, r *http.Request) {
		data, err := rows.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/rows", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormRow(w, r)
		if err == nil {
			id, err = rows.Add(info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Put("/rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		info, err := ParseFormRow(w, r)
		if err == nil {
			id = NumberParam(r, "id")
			err = rows.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := rows.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})
}
