package controller

import "net/http"

func SendFollowingRequestHandler(w http.ResponseWriter, r *http.Request)

func CancelFollowingRequestHandler(w http.ResponseWriter, r *http.Request)

func SendFollowingResponseHandler(w http.ResponseWriter, r *http.Request)

func UnfollowHandler(w http.ResponseWriter, r *http.Request)

func ViewFollowingRequestsHandler(w http.ResponseWriter, r *http.Request)

func ViewFollowersHandler(w http.ResponseWriter, r *http.Request)

func ViewLeadersHandler(w http.ResponseWriter, r *http.Request)
