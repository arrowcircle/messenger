package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// APIIndex is status function
func APIIndex(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("Welcome!\n")
}

// APIDialogIndex is used to get dialogs index
func (i *Impl) APIDialogIndex(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	page, err := strconv.Atoi(r.FormValue("page"))
	offset := 0
	if err == nil {
		offset = (page - 1) * 10
	}

	dialogs := i.GetDialogs(userID, offset)
	w.WriteJson(&dialogs)
}

// APIDialogShow is used to show dialog for RAILS
func (i *Impl) APIDialogShow(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))

	dialog := i.ShowDialog(userID, dialogID)

	w.WriteJson(&dialog)
}

// APIMessageIndex is used to show dialog messages
func (i *Impl) APIMessageIndex(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))
	page, err := strconv.Atoi(r.FormValue("page"))
	offset := 0
	if err == nil {
		offset = (page - 1) * 10
	}

	messages := i.IndexMessages(userID, dialogID, offset)

	w.WriteJson(&messages)
}

// APIUserShow returns number of unread dialogs
func (i *Impl) APIUserShow(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	user := i.ShowUser(userID)
	w.WriteJson(&user)
}

// APIDialogCreate is used to create dialog and message
func (i *Impl) APIDialogCreate(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	params := DialogCreateJSON{}
	if err := r.DecodeJsonPayload(&params); err != nil {
		fmt.Println("error decoding json: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dialog, err := i.CreateDialog(userID, params)

	if err != nil {
		fmt.Println("error creating dialog: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&dialog)
}

// APIMessageCreate creates message for dialog
func (i *Impl) APIMessageCreate(w rest.ResponseWriter, r *rest.Request) {
	userID := r.PathParam("user_id")
	dialogID, _ := strconv.Atoi(r.PathParam("dialog_id"))
	message := Message{}
	if err := r.DecodeJsonPayload(&message); err != nil {
		fmt.Println("error decoding json: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message, err := i.CreateMessage(userID, dialogID, message)
	if err != nil {
		fmt.Println("error creating message: ", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&message)
}
