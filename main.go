package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate/migrate"
	"github.com/spf13/viper"
)
import _ "github.com/mattes/migrate/driver/postgres"

func main() {
	readConfig()
	applyMigrations()
	i := Impl{}
	i.connectToDb()
	i.startChat()
}

func readConfig() {
	viper.SetEnvPrefix("chat")
	viper.SetDefault("database_url", "postgres:///chat_development?sslmode=disable")
	viper.SetDefault("bind_address", "localhost:8080")
	viper.AutomaticEnv()
	viper.BindEnv("database_url")
	viper.BindEnv("bind_address")
}

func applyMigrations() {
	allErrors, ok := migrate.UpSync(viper.GetString("database_url"), "./migrations")
	fmt.Println("Database: ", viper.GetString("database_url"))
	if !ok {
		fmt.Println("Migratinos failed")
		fmt.Println("driver: ", viper.GetString("database_url"))
		fmt.Println("Errors: ", allErrors)
	}
}

func (i *Impl) connectToDb() {
	var err error
	i.DB, err = gorm.Open("postgres", viper.GetString("database_url"))
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

func (i *Impl) startChat() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", APIIndex),
		rest.Get("/users/:user_id.json", i.APIUserShow),
		rest.Get("/users/:user_id/dialogs.json", i.APIDialogIndex),
		rest.Get("/users/:user_id/dialogs/:dialog_id/messages.json", i.APIMessageIndex),
		rest.Get("/users/:user_id/dialogs/:dialog_id.json", i.APIDialogShow),
		rest.Post("/users/:user_id/dialogs.json", i.APIDialogCreate),
		rest.Post("/users/:user_id/dialogs/:dialog_id/messages.json", i.APIMessageCreate),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	fmt.Println("address: ", viper.GetString("bind_address"))
	log.Fatal(http.ListenAndServe(viper.GetString("bind_address"), api.MakeHandler()))
}
