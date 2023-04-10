package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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
