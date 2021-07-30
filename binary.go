package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type DataManager struct{}

func (m *DataManager) FromRequest(r *http.Request, field string) (int, error) {
	// 10 MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile(field)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile(Config.BinaryData, "u*")
	if err != nil {
		return 0, err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return 0, err
	}

	data := BinaryData{
		Name: handler.Filename,
		Path: tempFile.Name(),
	}
	db.Save(&data)

	return data.ID, err
}

func (m *DataManager) ToResponse(w http.ResponseWriter, id int) (bool, error) {
	data := BinaryData{}
	err := db.Find(&data, id).Error
	if err != nil || data.ID == 0 {
		return false, err
	}

	file, err := os.Open(path.Join(Config.BinaryData, data.Path))
	if err != nil || data.ID == 0 {
		return false, err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return true, err
}
