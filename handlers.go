package main

import "encoding/json"
import "net/http"

import "github.com/gorilla/mux"

func SendJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func Index(w http.ResponseWriter, r *http.Request) {
	statistics := Statistics{
		Tokens:       len(tokens),
		Accounts:     len(accounts),
		Reservations: len(reservations),
	}
	SendJSON(w, statistics)
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	SendJSON(w, NewError("Token not authorized."))
}

func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	if token, ok := tokens[r.FormValue("token")]; ok {
		reservation := NewReservation(token)
		http.Redirect(w, r, "/reservations/"+reservation.Id+"?token="+token.Id, http.StatusFound)
	} else {
		Unauthorized(w)
	}
}

func GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if token, ok := tokens[r.FormValue("token")]; ok {
		if reservation := reservations[vars["reservationId"]]; reservation.Token == token {
			SendJSON(w, reservation)
		} else {
			Unauthorized(w)
		}
	} else {
		Unauthorized(w)
	}
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if token, ok := tokens[r.FormValue("token")]; ok {
		if reservation := reservations[vars["reservationId"]]; reservation.Token == token {
			reservation.Delete()
			SendNoContent(w)
		} else {
			Unauthorized(w)
		}
	} else {
		Unauthorized(w)
	}
}
