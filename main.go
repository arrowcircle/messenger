package main

import (
  "github.com/spf13/viper"
  "github.com/mattes/migrate/migrate"
  "github.com/ant0ine/go-json-rest/rest"
  _ "github.com/lib/pq"
  "github.com/jinzhu/gorm"
  "fmt"
  "net/http"
  "log"
)

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

type Impl struct {
  DB gorm.DB
}

func (i *Impl) startChat() {
  api := rest.NewApi()
  api.Use(rest.DefaultDevStack...)
  router, err := rest.MakeRouter(
    rest.Get("/", Index),
    rest.Get("/users/:user_id/dialogs", i.DialogIndex),
    rest.Get("/users/:user_id/dialogs/:dialog_id", DialogShow),
    rest.Post("/users/:user_id/dialogs", DialogCreate),
    rest.Post("/users/:user_id/dialogs/:dialog_id/messages", MessageCreate),
  )
  if err != nil {
    log.Fatal(err)
  }
  api.SetApp(router)
  fmt.Println("address: ", viper.GetString("bind_address"))
  log.Fatal(http.ListenAndServe(viper.GetString("bind_address"), api.MakeHandler()))
}

func Index(w rest.ResponseWriter, r *rest.Request) {
  w.WriteJson("Welcome!\n")
}

func (i *Impl) DialogIndex(w rest.ResponseWriter, r *rest.Request) {
  dialogs := []Dialog{}
  userId := r.PathParam("user_id")
  i.DB.Table("dialogs").Joins("INNER JOIN dialog_users on dialog_users.dialog_id = dialogs.id").Where("dialog_users.user_id = ?", userId).Find(&dialogs)
  w.WriteJson(&dialogs)
}

func DialogShow(w rest.ResponseWriter, r *rest.Request) {
  // fmt.Fprintf(w, "hello, %s!\n", ps.ByName("dialog_id"))
}

func DialogCreate(w rest.ResponseWriter, r *rest.Request) {
  // fmt.Fprintf(w, "create, %s!\n", ps.ByName("user_id"))
}

func MessageCreate(w rest.ResponseWriter, r *rest.Request) {
  // fmt.Fprintf(w, "create, %s!\n", ps.ByName("dialog_id"))
}
