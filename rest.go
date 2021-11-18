package main

import (
	"net/http"
	"web-widgets/kanban-go/data"

	"github.com/go-chi/chi"
)

func initRoutes(r chi.Router, dao *data.DAO) {

	r.Get("/cards", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Cards.GetAll()
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
			id, err = dao.Cards.Add(info)
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
			err = dao.Cards.Update(id, info)
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
			err = dao.Cards.Move(id, info.Card, int(info.Before))
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := dao.Cards.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})

	r.Get("/uploads/{id}/{name}", func(w http.ResponseWriter, r *http.Request) {
		res, err := dao.Files.ToResponse(w, NumberParam(r, "id"))

		if err != nil {
			format.Text(w, 500, err.Error())
		} else if !res {
			format.Text(w, 500, "")
		}
	})

	r.Post("/uploads", func(w http.ResponseWriter, r *http.Request) {
		rec, err := dao.Files.FromRequest(r, "upload")
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, rec)
		}
	})

	r.Get("/columns", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Columns.GetAll()
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
			id, err = dao.Columns.Add(info)
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
			err = dao.Columns.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := dao.Columns.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})

	r.Get("/rows", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Rows.GetAll()
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
			id, err = dao.Rows.Add(info)
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
			err = dao.Rows.Update(id, info)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}
	})

	r.Delete("/rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := dao.Rows.Delete(NumberParam(r, "id"))
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}
	})
}
