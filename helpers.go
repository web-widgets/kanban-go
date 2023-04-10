package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"web-widgets/kanban-go/data"

	"github.com/go-chi/chi"
)

type Response struct {
	ID int `json:"id"`
}

func NumberParam(r *http.Request, key string) int {
	value := chi.URLParam(r, key)
	num, _ := strconv.Atoi(value)

	return num
}

func ParseFormCard(w http.ResponseWriter, r *http.Request) (data.CardUpdate, error) {
	c := data.CardUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseFormMoveCard(w http.ResponseWriter, r *http.Request) (data.CardPosUpdate, error) {
	c := data.CardPosUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseFormColumn(w http.ResponseWriter, r *http.Request) (data.ColumnUpdate, error) {
	c := data.ColumnUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseFormColumnMove(w http.ResponseWriter, r *http.Request) (data.ColumnMove, error) {
	c := data.ColumnMove{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseFormRowMove(w http.ResponseWriter, r *http.Request) (data.RowMove, error) {
	row := data.RowMove{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&row)

	return row, err
}

func ParseFormRow(w http.ResponseWriter, r *http.Request) (data.RowUpdate, error) {
	c := data.RowUpdate{}

	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&c)

	return c, err
}

func ParseForm(w http.ResponseWriter, r *http.Request, o interface{}) error {
	body := http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(body)
	err := dec.Decode(&o)

	return err
}

func geDeviceID(r *http.Request) int {
	v := r.Context().Value("device_id")
	asInt, _ := v.(int)
	return asInt
}

func getUserID(r *http.Request) int {
	v := r.Context().Value("user_id")
	asInt, _ := v.(int)
	return asInt
}
