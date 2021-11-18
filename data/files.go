package data

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"gorm.io/gorm"
)

func NewFilesDAO(db *gorm.DB, httpPath, storePath string) *FilesDAO {
	return &FilesDAO{db, storePath, httpPath}
}

type FilesDAO struct {
	db        *gorm.DB
	storePath string
	httpPath  string
}

func (m *FilesDAO) FromRequest(r *http.Request, field string) (BinaryData, error) {
	// 10 MB
	r.ParseMultipartForm(10 << 20)
	rec := BinaryData{}

	file, handler, err := r.FormFile(field)
	if err != nil {
		return rec, err
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile(m.storePath, "u*")
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
	err = m.db.Save(&rec).Error
	if err != nil {
		return rec, err
	}

	//FIXME - use GUID for ID
	rec.URL = fmt.Sprintf("%s/uploads/%d/%s", m.httpPath, rec.ID, rec.Name)
	err = m.db.Save(&rec).Error

	return rec, err
}

func (m *FilesDAO) ToResponse(w http.ResponseWriter, id int) (bool, error) {
	data := BinaryData{}
	err := m.db.Find(&data, id).Error
	if err != nil || data.ID == 0 {
		return false, err
	}

	file, err := os.Open(path.Join(m.storePath, data.Path))
	if err != nil || data.ID == 0 {
		return false, err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return true, err
}
