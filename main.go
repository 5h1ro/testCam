package main

import (
	"testcam/handler/rest"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rest.StartApp()
}
