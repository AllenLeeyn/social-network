package controller

import "net/http"

func CreateEventHandler(w http.ResponseWriter, r *http.Request)

func UpdateEventHandler(w http.ResponseWriter, r *http.Request)

func ViewEventsHandler(w http.ResponseWriter, r *http.Request)

func ViewEventHandler(w http.ResponseWriter, r *http.Request)

func SendEventResponseHandler(w http.ResponseWriter, r *http.Request)

func ViewEventResponsesHandler(w http.ResponseWriter, r *http.Request)
