package main

import "time"

type Statistics struct {
	TotalReservations   int `json:"total_reservations"`
	PendingReservations int `json:"pending_reservations"`
	TotalAccounts       int `json:"total_accounts"`
	DirtyAccounts       int `json:"dirty_accounts"`
	ReservedAccounts    int `json:"reserved_accounts"`
	CleaningAccounts    int `json:"cleaning_accounts"`
}

type Token struct {
	Id string
}

type Reservation struct {
	Id         string    `json:"id"`
	Duration   time.Time `json:"duration"`
	Expiration time.Time `json:"expiration"`
	Token      *Token
	Account    *Account
}

type Account struct {
	Username string
	Password string
}

var tokens = make(map[string]*Token)
var reservations = make(map[string]*Reservation)
var accounts = make(map[string]*Account)
