package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", jwtMiddleware(makeHTTPHandler(s.handleAccountById)))
	router.HandleFunc("/login", makeHTTPHandler(s.handleLogin))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("Method: %s not alslowed.", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("Method: %s not allowed.", r.Method)
	}

	loginReq := new(LoginAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	if !IsValidLoginRequest(loginReq) {
		return fmt.Errorf("Invalid request.")
	}

	acc, err := s.store.GetAccountByEmail(loginReq.Email)

	if err != nil {
		return err
	}

	if !acc.IsValidPassword(loginReq.Password) {
		return fmt.Errorf("Not authenticated.")
	}

	token, err := CreateToken(acc.Email)

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, LoginResponse{
		Token: token,
		Email: acc.Email,
	})

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	registerReq := new(RegisterAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(registerReq); err != nil {
		return err
	}

	defer r.Body.Close()

	if !IsValidRegisterRequest(registerReq) {
		return fmt.Errorf("Invalid request!")
	}

	acc, err := CreateNewAccount(registerReq.FirstName, registerReq.LastName, registerReq.Email, registerReq.Password)

	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(acc); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, acc)

}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountById(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccountById(w, r)
	}

	return fmt.Errorf("Method: %s not allowed.", r.Method)

}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	id, err := intoInt(r)

	if err != nil {
		return err
	}

	acc, err := s.store.GetAccountByID(id)

	fmt.Println(acc)

	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, acc)

}

func (s *APIServer) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := intoInt(r)

	if err != nil {
		return err
	}

	return s.store.DeleteAccount(id)
}

func makeHTTPHandler(f apiFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func jwtMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		claims, err := validateToken(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		_ = claims

		f(w, r)

	}
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ""
	}

	return tokenParts[1]
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiError struct {
	Error string `json:"error"`
}

func intoInt(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

type apiFunc func(http.ResponseWriter, *http.Request) error