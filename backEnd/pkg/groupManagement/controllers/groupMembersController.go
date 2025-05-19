package controller

import "net/http"

func ViewGroupMembersHandle(w http.ResponseWriter, r *http.Request)

func ViewGroupMembershipHandle(w http.ResponseWriter, r *http.Request)

func SendGroupInvitationHandler(w http.ResponseWriter, r *http.Request)

func CancelGroupInvitationHandler(w http.ResponseWriter, r *http.Request)

func SendGroupJoinRequestHandler(w http.ResponseWriter, r *http.Request)

func CancelGroupJoinRequestHandler(w http.ResponseWriter, r *http.Request)

func SendGroupMembershipResponseHandler(w http.ResponseWriter, r *http.Request)
