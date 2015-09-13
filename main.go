package main

import (
  "github.com/spf13/viper"
  "github.com/mattes/migrate/migrate"
  "github.com/julienschmidt/httprouter"
  _ "github.com/lib/pq"
  "database/sql"
  "fmt"
  "net/http"
  "log"
)

func main() {
  readConfig()
  applyMigrations()
  connectToDb()
  startChat()
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

func connectToDb() {
  var db *sql.DB
  var err error

  db, err = sql.Open("postgres", viper.GetString("database_url"))

  if err != nil {
    fmt.Printf("sql.Open error: %v\n",err)
    return
  }

  defer db.Close()
}

func startChat() {
  router := httprouter.New()

  router.GET("/", Index)
  router.GET("/users/:user_id/dialogs", DialogIndex)
  router.GET("/users/:user_id/dialogs/:dialog_id", DialogShow)
  router.POST("/users/:user_id/dialogs", DialogCreate)
  router.POST("/users/:user_id/dialogs/:dialog_id/messages", MessageCreate)
  fmt.Println("address: ", viper.GetString("bind_address"))
  log.Fatal(http.ListenAndServe(viper.GetString("bind_address"), router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

func DialogIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "hello, %s!\n", ps.ByName("user_id"))
}

func DialogShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "hello, %s!\n", ps.ByName("dialog_id"))
}

func DialogCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "create, %s!\n", ps.ByName("user_id"))
}

func MessageCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  fmt.Fprintf(w, "create, %s!\n", ps.ByName("dialog_id"))
}
