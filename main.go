package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/didi/gendry/scanner"
	"github.com/spf13/viper"
)

func init() {

	configRoot, _ := os.Getwd()
	var configPath = flag.String("config", configRoot, "defines path to config file")

	flag.Parse()

	// setup config file
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.AddConfigPath(*configPath) // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Couldn't find config.toml!")
		log.Fatal(err)
	}

}

func main() {

	// create db instance
	db, err := sql.Open("mysql", viper.GetString("connection"))
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// set tag name for sql scanner
	scanner.SetTagName("json")

	// create router with db instance
	router := APIRouter(&RouterConfig{
		Handlers: &Handlers{
			DB: db,
		},
	})

	fmt.Println("biokiste backend listen on localhost:1316")
	log.Fatal(http.ListenAndServe(":1316", router))
}
