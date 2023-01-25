package rest

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	port := ":9000"
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// sts := `
	// DROP TABLE IF EXISTS medias;
	// CREATE TABLE medias(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, extension TEXT);
	// `
	// _, err = db.Exec(sts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	route := gin.Default()
	filesHandler := NewFile(db)
	route.LoadHTMLFiles("index.html")
	route.Static("/media", "./media")
	route.GET("/", func(c *gin.Context) {
		data, _ := db.Query(`SELECT * FROM medias`)
		var rows []Entity
		for data.Next() {
			var id, name, extension sql.NullString
			data.Scan(&id, &name, &extension)
			rows = append(rows, Entity{ID: id.String, Name: name.String, Extension: extension.String})
		}
		// isImage := rows[0].Extension == ".aac" || rows[0].Extension == ".wav" || rows[0].Extension == ".mp3"

		c.HTML(http.StatusOK, "index.html", gin.H{
			"data": rows,
		})
	})
	route.POST("/create", filesHandler.createFile)

	fmt.Println("Server running on PORT =>", port)
	route.Run(port)
}
