package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"splitwise/domains"
	"splitwise/repository"
	"splitwise/service"
	"splitwise/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userinfo, _ := ioutil.ReadAll(r.Body)
	request := domains.User{}
	json.Unmarshal(userinfo, &request)
	message := domains.Message{}

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := utils.GetErrorCode(validationErrors)
		w.WriteHeader(http.StatusBadRequest)
		message.Message = errorMessage
		json.NewEncoder(w).Encode(message)
		return
	}

	err = service.RegisterUser(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Message = err.Error()
		json.NewEncoder(w).Encode(message)
		return
	}

	message.Message = "User registered successfully"

	json.NewEncoder(w).Encode(message)
}

func addExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	reqinfo, _ := ioutil.ReadAll(r.Body)
	request := domains.Expense{}
	json.Unmarshal(reqinfo, &request)
	message := domains.Message{}
	code, err := service.AddExpense(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if code == "invalid_request" {
			w.WriteHeader(http.StatusBadRequest)
		}
		message.Message = err.Error()
		json.NewEncoder(w).Encode(message)
		return
	}

	message.Message = "Expense added successfully"

	json.NewEncoder(w).Encode(message)
}

func getBalances(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId := mux.Vars(r)["userId"]
	message := domains.Message{}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message.Message = "Invalid userId passed"
		json.NewEncoder(w).Encode(message)
		return
	}

	balances, err := service.GetBalances(userIdInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Message = err.Error()
		json.NewEncoder(w).Encode(message)
		return
	}

	json.NewEncoder(w).Encode(balances)

}

func getExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	expenses, err := service.GetExpenses()
	message := domains.Message{}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message.Message = err.Error()
		json.NewEncoder(w).Encode(message)
		return
	}

	json.NewEncoder(w).Encode(expenses)

}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/register", registerUser).Methods("POST")
	myRouter.HandleFunc("/expense", addExpense).Methods("POST")
	myRouter.HandleFunc("/expenses", getExpenses).Methods("GET")
	myRouter.HandleFunc("/{userId}/balances", getBalances).Methods("GET")
	fmt.Println("Started Server")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

func main() {
	validate = validator.New()
	err := repository.StartDB()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error in connecting with DB")
	}
	HandleRequests()

}
