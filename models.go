package main

import "time"
import "github.com/ernsheong/grand"

type Statistics struct {
	Tokens       int `json:"tokens"`
	Reservations int `json:"reservations"`
	Accounts     int `json:"accounts"`
}

type Error struct {
	Message string `json:"error"`
}

func NewError(message string) Error {
	return Error{
		Message: message,
	}
}

func NewId() string {
	generator := grand.NewGenerator(grand.CharSetBase62)
	return generator.GenerateRandomString(8)
}

type Token struct {
	Id string `json:"id"`
}

func NewToken() Token {
	return AddToken(NewId())
}

func AddToken(tokenId string) Token {
	token := Token{
		Id: tokenId}
	tokens[token.Id] = &token
	return token
}

type Reservation struct {
	Id         string    `json:"id"`
	Expiration time.Time `json:"expiration"`
	Token      *Token    `json:"token"`
	Account    *Account  `json:"credentials"`
}

func NewReservation(token *Token) Reservation {
	reservation := Reservation{
		Id:    NewId(),
		Token: token,
	}
	reservations[reservation.Id] = &reservation
	return reservation
}

func (r *Reservation) Delete() {
	delete(reservations, r.Id)
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AddAccount(username string, password string) {
	account := Account{
		Username: username,
		Password: password,
	}
	accounts[username] = &account
}

var tokens = make(map[string]*Token)
var reservations = make(map[string]*Reservation)
var accounts = make(map[string]*Account)
