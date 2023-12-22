package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"splitwise/domains"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	contextTimeout = 30 * time.Second
)

var (
	internalError = errors.New("internal server error")
)

const (
	addUserQuery     = "insert into Users (username,email,phone,addedOn) values (?,?,?,?)"
	addExpenseQuery  = "insert into Expenses (name,split,createdBy,amount,paidBy,paidFor,createdOn) values (?,?,?,?,?,?,?)"
	addBalanceEntry  = "insert into Balances (paidBy,paidFor,amount,expenseId,addedOn) values "
	getOwesBalances  = "select paidBy,sum(amount) as amount from Balances where paidFor=? group by paidBy"
	getGetsBalances  = "select paidFor,sum(amount) as amount from Balances where paidBy=? group by paidFor"
	getExpensesQuery = "select * from Expenses order by id desc"
)

var r *SplitWiseRepository

type SplitWiseRepository struct {
	Connection *sqlx.DB
}

func StartDB() error {
	conn, err := DbConnection()
	if err != nil {
		return err
	}
	r = &SplitWiseRepository{
		Connection: conn,
	}

	return nil
}

//AddUser is for adding a new user
func AddUser(user domains.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	_, err := r.Connection.ExecContext(ctx, addUserQuery, user.Username, user.Email, user.Phone, time.Now())
	if err != nil {
		log.Printf("error in adding user : %s", err.Error())
		return internalError
	}

	return nil
}

//AddExpense peforms a transaction to add an expense and create entries in Balances table
func AddExpense(expense domains.Expense, balances map[string]map[string]float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	// Begin transaction
	tx, err := r.Connection.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("error in beginning transaction %s", err.Error())
		tx.Rollback()
		return internalError
	}

	paidBy, _ := json.Marshal(expense.PaidBy)
	paidFor, _ := json.Marshal(expense.PaidFor)

	// adding an expense
	res, err := tx.ExecContext(ctx, addExpenseQuery, expense.ExpenseName, expense.SplitType, expense.CreatedBy, expense.Amount, paidBy, paidFor, time.Now())
	if err != nil {
		log.Printf("error in adding expense %s", err.Error())
		tx.Rollback()
		return internalError
	}

	expenseID, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in fetching last inserted id %s", err.Error())
		tx.Rollback()
		return internalError
	}

	query := addBalanceEntry
	args := []interface{}{}

	for paidBy, balance := range balances {
		for paidFor, amount := range balance {
			query = query + " (?,?,?,?,?),"
			args = append(args, paidBy, paidFor, amount, expenseID, time.Now())
		}
	}

	query = query[:len(query)-1]

	// running bulk insert query
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("error in adding balance entries : %s", err.Error())
		tx.Rollback()
		return internalError
	}

	if tx.Commit() != nil {
		tx.Rollback()
		log.Printf("error in fetching last inserted id %s", err.Error())
		return internalError
	}

	return nil
}

//GetBalances gives the amount a user owes to others or gets from others
func GetBalances(userId int) (*domains.Balances, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	owes := &[]domains.OwesBalances{}
	gets := &[]domains.GetsBalances{}

	args := []interface{}{userId}

	err := r.Connection.SelectContext(ctx, owes, getOwesBalances, args...)
	if err != nil {
		log.Println("error in getting balances", err.Error())
		return nil, internalError
	}

	err = r.Connection.SelectContext(ctx, gets, getGetsBalances, args...)
	if err != nil {
		log.Println("error in getting balances", err.Error())
		return nil, internalError
	}

	balances := &domains.Balances{
		Owes: make(map[int]float64),
		Gets: make(map[int]float64),
	}

	for _, entry := range *owes {
		balances.Owes[entry.OwesTo] = entry.Amount
	}

	for _, entry := range *gets {
		balances.Gets[entry.GetsFrom] = entry.Amount
	}

	return balances, nil
}

func GetExpenses() ([]domains.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	response := []domains.ExpenseDBEntry{}

	err := r.Connection.SelectContext(ctx, &response, getExpensesQuery)
	if err != nil {
		log.Println("error in getting expense", err.Error())
		return nil, internalError
	}

	result := []domains.Expense{}

	for _, entry := range response {
		row := domains.Expense{
			CreatedBy:   entry.CreatedBy,
			SplitType:   entry.SplitType,
			Amount:      entry.Amount,
			CreatedOn:   entry.CreatedOn,
			ExpenseName: entry.ExpenseName,
			ExpenseID:   entry.ExpenseID,
		}
		paidBy := make(map[string]float64)
		paidFor := make(map[string]float64)
		json.Unmarshal([]byte(entry.PaidBy), &paidBy)
		json.Unmarshal([]byte(entry.PaidFor), &paidFor)
		row.PaidBy = paidBy
		row.PaidFor = paidFor
		result = append(result, row)
	}

	return result, nil
}
