package rest

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	db *sql.DB
}

func NewFile(db *sql.DB) *File {
	return &File{db}
}

type Entity struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
}

func (f *File) CreateFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	if err := c.SaveUploadedFile(file, "./media/"+newFileName); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	stmt, err := f.db.Prepare(`INSERT INTO medias(name, extension) VALUES (?, ?)`)
	if err != nil {
		fmt.Println(err)
	}
	response, err := stmt.Exec(newFileName, extension)
	_ = response
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	data, _ := f.db.Query(`SELECT * FROM medias WHERE name = ?`, newFileName)
	var rows []Entity
	for data.Next() {
		var id, name, extension sql.NullString
		data.Scan(&id, &name, &extension)
		rows = append(rows, Entity{ID: id.String, Name: name.String, Extension: extension.String})
	}

	res := map[string]interface{}{
		"message": fmt.Sprintf("File %s uploaded successfully", file.Filename),
		"data":    rows,
	}

	c.JSON(http.StatusOK, res)
}

func (f *File) DeleteAllFile(c *gin.Context) {
	stmt, err := f.db.Prepare(`DELETE FROM medias`)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

	res := map[string]interface{}{
		"message": fmt.Sprintf("All file deleted"),
	}

	c.JSON(http.StatusOK, res)
}

func (f *File) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	stmt, err := f.db.Prepare(`DELETE FROM medias WHERE id = ?`)
	if err != nil {
		fmt.Println(id)
	}
	stmt.Exec(id)

	res := map[string]interface{}{
		"message": fmt.Sprintf("All file with id => " + id),
	}

	c.JSON(http.StatusOK, res)
}
