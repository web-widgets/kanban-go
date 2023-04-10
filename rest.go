package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"web-widgets/kanban-go/api"
	"web-widgets/kanban-go/data"

	"github.com/go-chi/chi"
	remote "github.com/mkozhukh/go-remote"
)

var errFeatureDisabled error = errors.New("feature disabled")

func initRoutes(r chi.Router, dao *data.DAO, hub *remote.Hub) {

	r.Get("/cards", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Cards.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Get("/cards/column/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := NumberParam(r, "id")
		data, err := dao.Cards.GetColumn(id)
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/cards", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var upd data.CardUpdate
		err := ParseForm(w, r, &upd)
		if err == nil {
			id, err = dao.Cards.Add(upd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			card, _ := dao.Cards.GetOne(id)
			hub.Publish("cards", api.CardEvent{
				Type: "add-card",
				From: geDeviceID(r),
				Card: card,
			})
		}
	})

	r.Put("/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var cardUpd data.CardUpdate
		err := ParseForm(w, r, &cardUpd)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Cards.Update(id, cardUpd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			card, _ := dao.Cards.GetOne(id)
			hub.Publish("cards", api.CardEvent{
				Type: "update-card",
				From: geDeviceID(r),
				Card: card,
			})
		}
	})

	r.Put("/cards/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var moveInfo data.CardPosUpdate
		err := ParseForm(w, r, &moveInfo)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Cards.Move(id, moveInfo)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})
		}

		hub.Publish("cards", api.CardEvent{
			Type: "move-card",
			From: geDeviceID(r),
			Card: &data.Card{
				ID:       id,
				ColumnID: int(moveInfo.ColumnID),
				RowID:    int(moveInfo.RowID),
			},
			Before: int(moveInfo.Before),
		})
	})

	r.Delete("/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := NumberParam(r, "id")
		err := dao.Cards.Delete(id)
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})
		}

		hub.Publish("cards", api.CardEvent{
			Type: "delete-card",
			From: geDeviceID(r),
			Card: &data.Card{ID: id},
		})
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
		var columnUpd data.ColumnUpdate
		err := ParseForm(w, r, &columnUpd)
		if err == nil {
			id, err = dao.Columns.Add(columnUpd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			column, _ := dao.Columns.GetOne(id)
			hub.Publish("columns", api.ColumnEvent{
				Type:   "add-column",
				From:   geDeviceID(r),
				Column: column,
			})
		}
	})

	r.Put("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var columnUpd data.ColumnUpdate
		err := ParseForm(w, r, &columnUpd)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Columns.Update(id, columnUpd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			column, _ := dao.Columns.GetOne(id)
			hub.Publish("columns", api.ColumnEvent{
				Type:   "update-column",
				From:   geDeviceID(r),
				Column: column,
			})
		}
	})

	r.Put("/columns/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var moveInfo data.ColumnMove
		err := ParseForm(w, r, &moveInfo)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Columns.Move(id, int(moveInfo.Before))
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			column, _ := dao.Columns.GetOne(id)
			hub.Publish("columns", api.ColumnEvent{
				Type:   "move-column",
				From:   geDeviceID(r),
				Column: column,
				Before: int(moveInfo.Before),
			})
		}
	})

	r.Delete("/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := NumberParam(r, "id")
		column, _ := dao.Columns.GetOne(id)
		err := dao.Columns.Delete(id)
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})

			hub.Publish("columns", api.ColumnEvent{
				Type:   "delete-column",
				From:   geDeviceID(r),
				Column: column,
			})
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
		var rowUpd data.RowUpdate
		err := ParseForm(w, r, &rowUpd)
		if err == nil {
			id, err = dao.Rows.Add(rowUpd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			row, _ := dao.Rows.GetOne(id)
			hub.Publish("rows", api.RowEvent{
				Type: "add-row",
				From: geDeviceID(r),
				Row:  row,
			})
		}
	})

	r.Put("/rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var rowUpd data.RowUpdate
		err := ParseForm(w, r, &rowUpd)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Rows.Update(id, rowUpd)
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			row, _ := dao.Rows.GetOne(id)
			hub.Publish("rows", api.RowEvent{
				Type: "update-row",
				From: geDeviceID(r),
				Row:  row,
			})
		}
	})

	r.Put("/rows/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		var id int
		var moveInfo data.RowMove
		err := ParseForm(w, r, &moveInfo)
		if err == nil {
			id = NumberParam(r, "id")
			err = dao.Rows.Move(id, int(moveInfo.Before))
		}
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{id})

			row, _ := dao.Rows.GetOne(id)
			hub.Publish("rows", api.RowEvent{
				Type:   "move-row",
				From:   geDeviceID(r),
				Row:    row,
				Before: int(moveInfo.Before),
			})
		}
	})

	r.Delete("/rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := NumberParam(r, "id")
		row, _ := dao.Rows.GetOne(id)
		err := dao.Rows.Delete(id)
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{})

			hub.Publish("rows", api.RowEvent{
				Type: "delete-row",
				From: geDeviceID(r),
				Row:  row,
			})
		}
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		data, err := dao.Users.GetAll()
		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, data)
		}
	})

	r.Post("/cards/{id}/vote", func(w http.ResponseWriter, r *http.Request) {
		if !Config.Features.WithVotes {
			format.Text(w, 500, errFeatureDisabled.Error())
			return
		}

		cid := NumberParam(r, "id")
		userId := getUserID(r)
		err := dao.Cards.SetVote(cid, userId)

		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{userId})
		}
	})

	r.Delete("/cards/{id}/vote", func(w http.ResponseWriter, r *http.Request) {
		if !Config.Features.WithVotes {
			format.Text(w, 500, errFeatureDisabled.Error())
			return
		}

		cid := NumberParam(r, "id")
		userId := getUserID(r)
		err := dao.Cards.RemoveVote(cid, userId)

		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{userId})
		}
	})

	r.Post("/comments", func(w http.ResponseWriter, r *http.Request) {
		if !Config.Features.WithComments {
			format.Text(w, 500, errFeatureDisabled.Error())
			return
		}

		comment := data.CommentUpdate{}
		err := ParseForm(w, r, &comment)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		userId := getUserID(r)
		commentId, err := dao.Comments.Add(userId, comment)

		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{commentId})
		}
	})

	r.Put("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !Config.Features.WithComments {
			format.Text(w, 500, errFeatureDisabled.Error())
			return
		}

		commentData := data.CommentUpdate{}
		err := ParseForm(w, r, &commentData)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		commentId := NumberParam(r, "id")
		userId := getUserID(r)
		err = dao.Comments.Update(commentId, userId, commentData)

		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{commentId})
		}
	})

	r.Delete("/comments/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !Config.Features.WithComments {
			format.Text(w, 500, errFeatureDisabled.Error())
			return
		}

		commentId := NumberParam(r, "id")
		userId := getUserID(r)
		err := dao.Comments.Delete(commentId, userId)

		if err != nil {
			format.Text(w, 500, err.Error())
		} else {
			format.JSON(w, 200, Response{userId})
		}
	})

	// DEMO ONLY, imitate login
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(r.URL.Query().Get("id"))
		device := newDeviceID()
		token, err := createUserToken(userId, device)
		if err != nil {
			log.Println("[token]", err.Error())
		}
		w.Write(token)
	})
}

var dID int

func init() {
	dID = int(time.Now().Unix())
}

func newDeviceID() int {
	dID += 1
	return dID
}
