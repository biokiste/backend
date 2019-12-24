package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// create db instance
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/foodkoop_biokiste")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create router with db instance
	router := APIRouter(&RouterConfig{
		Handlers: &Handlers{
			DB: db,
		},
	})

	fmt.Println("biokiste backend listen on localhost:1316")
	log.Fatal(http.ListenAndServe("localhost:1316", router))
}
