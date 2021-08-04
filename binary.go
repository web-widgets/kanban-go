package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type DataManager struct{}

func (m *DataManager) FromRequest(r *http.Request, field string) (BinaryData, error) {
	// 10 MB
	r.ParseMultipartForm(10 << 20)
	rec := BinaryData{}

	file, handler, err := r.FormFile(field)
	if err != nil {
		return rec, err
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile(Config.BinaryData, "u*")
	if err != nil {
		return rec, err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return rec, err
	}

	rec.Name = handler.Filename
	rec.Path = filepath.Base(tempFile.Name())
	err = db.Save(&rec).Error
	if err != nil {
		return rec, err
	}

	//FIXME - use GUID for ID
	rec.URL = fmt.Sprintf("%s/data/%d/%s", Config.Server.URL, rec.ID, rec.Name)
	err = db.Save(&rec).Error

	return rec, err
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
