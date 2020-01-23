package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func init() {

	configRoot, _ := os.Getwd()

	// setup config file
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath(configRoot) // path to look for the config file in
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Die Konfigurationsdatei config.toml konnte nicht gefunden werden!")
		log.Fatal(err)
	}

}

func main() {

	// create db instance
	db, err := sql.Open("mysql", viper.GetString("connection"))
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
	log.Fatal(http.ListenAndServe(":1316", router))
}
