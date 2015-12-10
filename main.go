package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate/migrate"
	"github.com/spf13/viper"
)
import _ "github.com/mattes/migrate/driver/postgres"

// UserJSON is used for empty requests
type UserJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DialogJSON used for index action of API
type DialogJSON struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UserIds           string    `json:"user_ids"`
	LastMessage       string    `json:"last_message"`
	LastMessageID     int       `json:"last_message_id"`
	LastMessageUserID int       `json:"last_message_user_id"`
	LastSeenMessageID int       `json:"last_seen_message_id"`
}

// DialogCreateJSON is used for dialogs creation
type DialogCreateJSON struct {
	Name    string `json:"name"`
	UserIds []int  `json:"user_ids"`
	Message string `json:"message"`
}

// Dialog is used to save dialogs into DB via GORM
type Dialog struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	LastMessageID int       `json:"last_message_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DialogShowJSON is used to form json
type DialogShowJSON struct {
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Messages []MessageJSON `json:"messages"`
}

// MessageJSON is used to response message in JSON format
type MessageJSON struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
}

// Message is used to put messages in DB via GORM
type Message struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	UserID   int    `json:"user_id"`
	DialogID int    `json:"dialog_id"`
}

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

// Impl used to provide handler to DB
type Impl struct {
	DB gorm.DB
}

func (i *Impl) startChat() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", Index),
		rest.Get("/users/:user_id.json", UserShow),
		rest.Get("/users/:user_id/dialogs.json", i.DialogIndex),
		rest.Get("/users/:user_id/dialogs/:dialog_id/messages.json", i.MessageIndex),
		rest.Get("/users/:user_id/dialogs/:dialog_id.json", i.DialogShow),
		rest.Post("/users/:user_id/dialogs.json", i.DialogCreate),
		rest.Post("/users/:user_id/dialogs/:dialog_id/messages.json", i.MessageCreate),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	fmt.Println("address: ", viper.GetString("bind_address"))
	log.Fatal(http.ListenAndServe(viper.GetString("bind_address"), api.MakeHandler()))
}

// Index is status function
func Index(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("Welcome!\n")
}

// DialogIndex is used to get dialogs index
func (i *Impl) DialogIndex(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	page, err := strconv.Atoi(r.FormValue("page"))
	offset := 0
	if err == nil {
		offset = (page - 1) * 10
	}

	dialogs := []DialogJSON{}
	i.DB.Raw(`
    SELECT c.*, array_agg(du.user_id) AS user_ids
    FROM
      (SELECT
        dialogs.id AS id,
        dialogs.name AS name,
        dialogs.created_at AS created_at,
        dialogs.updated_at AS updated_at,
        dialogs.last_message_id AS last_message_id,
        messages.text AS last_message,
        messages.user_id AS last_message_user_id,
    	  dialog_users.last_seen_message_id AS last_seen_message_id
      FROM dialogs
      JOIN messages ON messages.id = dialogs.last_message_id
      JOIN dialog_users ON dialog_users.dialog_id = dialogs.id
      WHERE dialog_users.user_id = ?
      ORDER BY dialogs.last_message_id DESC
      ) c
    JOIN dialog_users du ON c.id = du.dialog_id
    GROUP BY
      c.id,
      c.name,
      c.created_at,
      c.updated_at,
      c.last_message_id,
      c.last_message,
      c.last_message_user_id,
      c.last_seen_message_id
    LIMIT 10
    OFFSET ?
  `, userID, offset).Find(&dialogs)
	w.WriteJson(&dialogs)
}

// DialogShow is used to show dialog for RAILS
func (i *Impl) DialogShow(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))

	dialog := DialogJSON{}
	i.DB.Raw(`
    SELECT c.*, array_agg(du.user_id) AS user_ids
    FROM
      (SELECT
        dialogs.id AS id,
        dialogs.name AS name,
        dialogs.created_at AS created_at,
        dialogs.updated_at AS updated_at,
        dialogs.last_message_id AS last_message_id,
        messages.text AS last_message,
        messages.user_id AS last_message_user_id,
    	  dialog_users.last_seen_message_id AS last_seen_message_id
      FROM dialogs
      JOIN messages ON messages.id = dialogs.last_message_id
      JOIN dialog_users ON dialog_users.dialog_id = dialogs.id
      WHERE dialog_users.user_id = ?
      ORDER BY dialogs.last_message_id DESC
      ) c
    JOIN dialog_users du ON c.id = du.dialog_id
    WHERE c.id = ?
    GROUP BY
      c.id,
      c.name,
      c.created_at,
      c.updated_at,
      c.last_message_id,
      c.last_message,
      c.last_message_user_id,
      c.last_seen_message_id
  `, userID, dialogID).Find(&dialog)

	lastMessageID := 0
	i.DB.Raw("SELECT last_message_id FROM dialogs WHERE id = ?", dialogID).Row().Scan(&lastMessageID)
	i.DB.Exec("UPDATE dialog_users SET last_seen_message_id = ? WHERE dialog_id = ? AND user_id = ?", lastMessageID, dialogID, userID)

	w.WriteJson(&dialog)
}

// MessageIndex is used to show dialog messages
func (i *Impl) MessageIndex(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))
	page, err := strconv.Atoi(r.FormValue("page"))
	offset := 0
	if err == nil {
		offset = (page - 1) * 10
	}
	messages := []MessageJSON{}
	i.DB.Raw(`
    SELECT * FROM messages
    WHERE messages.dialog_id = ?
    ORDER BY messages.id DESC
    LIMIT 10
    OFFSET ?
  `, dialogID, offset).Find(&messages)

	lastMessageID := 0
	i.DB.Raw("SELECT last_message_id FROM dialogs WHERE id = ?", dialogID).Row().Scan(&lastMessageID)
	i.DB.Exec("UPDATE dialog_users SET last_seen_message_id = ? WHERE dialog_id = ? AND user_id = ?", lastMessageID, dialogID, userID)

	w.WriteJson(&messages)
}

// UserShow if fake method for rails
func UserShow(w rest.ResponseWriter, r *rest.Request) {
	userID, _ := strconv.Atoi(r.PathParam("user_id"))
	user := UserJSON{}
	user.ID = userID
	w.WriteJson(&user)
}

// DialogCreate is used to create dialog and message
func (i *Impl) DialogCreate(w rest.ResponseWriter, r *rest.Request) {
	userID, _ := strconv.Atoi(r.PathParam("user_id"))
	dialogJSON := DialogCreateJSON{}
	if err := r.DecodeJsonPayload(&dialogJSON); err != nil {
		fmt.Println("error decoding json: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dialog := Dialog{}
	dialog.Name = dialogJSON.Name

	if err := i.DB.Save(&dialog).Error; err != nil {
		fmt.Println("error saving message: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := Message{}
	message.DialogID = dialog.ID
	message.Text = dialogJSON.Message
	message.UserID = userID

	if err := i.DB.Save(&message).Error; err != nil {
		fmt.Println("error saving message: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, element := range dialogJSON.UserIds {
		i.DB.Exec("INSERT INTO dialog_users (dialog_id, user_id, last_seen_message_id) VALUES (?, ?, 0)", dialog.ID, element)
	}

	i.DB.Exec("UPDATE dialogs SET last_message_id = ? WHERE id = ?", message.ID, dialog.ID)

	dialog.LastMessageID = message.ID

	fmt.Println("dialog json: ", dialog.ID)

	w.WriteJson(&dialog)
}

// MessageCreate creates message for dialog
func (i *Impl) MessageCreate(w rest.ResponseWriter, r *rest.Request) {
	userID, _ := strconv.Atoi(r.PathParam("user_id"))
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))
	message := Message{}
	if err := r.DecodeJsonPayload(&message); err != nil {
		fmt.Println("error decoding json: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message.DialogID = dialogID
	message.UserID = userID
	if err := i.DB.Save(&message).Error; err != nil {
		fmt.Println("error saving message: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i.DB.Exec("UPDATE dialogs SET last_message_id = ? WHERE dialogs.id = ?", message.ID, message.DialogID)
	i.DB.Exec("UPDATE dialog_users SET last_seen_message_id = ? WHERE dialog_id = ? AND user_id = ?", message.ID, message.DialogID, message.UserID)

	w.WriteJson(&message)
}
