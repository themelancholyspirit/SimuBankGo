package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func CreateNewAccount(firstName string, lastName string, email string, password string) (*Account, error) {
	EncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Email:             email,
		EncryptedPassword: string(EncryptedPassword),
		Balance:           1000,
		CreatedAt:         time.Now().UTC(),
	}, nil

}

type RegisterAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type TransferMoneyRequest struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type UnauthorizedResponse struct {
	Error string `json:"error"`
}

type TransactionHistoryRequest struct {
	From              string    `json:"from"`
	To                string    `json:"to"`
	Amount            int       `json:"amount"`
	TransactionMadeAt time.Time `json:"transcationmadeat"`
}

func (a *Account) IsValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}
