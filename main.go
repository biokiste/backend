package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {	
	router := NewRouter()

	fmt.Println("biokiste backend listen on localhost:1316")
	log.Fatal(http.ListenAndServe("localhost:1316", router))
}
