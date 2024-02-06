package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	UpdateAccountBalance(string, int64) error
}

type postgredb struct {
	db *sql.DB
}

func NewPosgreDB() (*postgredb, error) {
	connectionString := fmt.Sprintf("user=postgres dbname=postgres password=%s sslmode=disable", os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgredb{
		db: db,
	}, nil
}

func (s *postgredb) CreateAccount(account *Account) error {
	queryStr := `insert into account 
	(first_name, last_name, email, encrypted_password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		queryStr,
		account.FirstName,
		account.LastName,
		account.Email,
		account.EncryptedPassword,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *postgredb) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		email varchar(50) unique,
		encrypted_password varchar(100),
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *postgredb) DeleteAccount(id int) error {

	_, err := s.db.Query("delete from account where id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

func (s *postgredb) UpdateAccount(account *Account) error {
	return nil
}

func (s *postgredb) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")

	s.db.Begin()

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		acc, err := scanAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)

	}

	return accounts, nil

}

func (s *postgredb) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)

	if rows.Next() {
		return scanAccount(rows)
	}

	return nil, err

}

func (s *postgredb) GetAccountByEmail(email string) (*Account, error) {
	rows, err := s.db.Query("select * from account where email = $1", email)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanAccount(rows)
	}

	return nil, fmt.Errorf("no such user")

}

func (s *postgredb) UpdateAccountBalance(userEmail string, newBalance int64) error {
	stmt, err := s.db.Prepare("UPDATE account SET balance = $1 WHERE email = $2")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newBalance, userEmail)

	if err != nil {
		return err
	}

	return nil

}

func scanAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)

	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)

	return acc, err

}
