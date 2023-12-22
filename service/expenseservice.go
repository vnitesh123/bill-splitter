package service

import (
	"errors"
	"log"
	"splitwise/domains"
	"splitwise/repository"
)

func AddExpense(expense domains.Expense) (string, error) {
	code := ""
	var total float64
	total = 0
	for _, val := range expense.PaidBy {
		total = total + val
	}

	if total != expense.Amount {
		code = "inval_request"
		log.Println("Invalid amounts received in the request")
		return code, errors.New("amount paid by users not matching the total amount")
	}

	paidBy := make(map[string]float64)
	paidFor := make(map[string]float64)
	for key, val := range expense.PaidBy {
		paidBy[key] = val
	}

	for key, val := range expense.PaidFor {
		paidFor[key] = val
	}

	var balances map[string]map[string]float64

	switch expense.SplitType {
	case "EQUAL":
		paidForLen := len(expense.PaidFor)
		for key := range expense.PaidFor {
			expense.PaidFor[key] = expense.Amount / float64(paidForLen)
		}

	case "PERCENTAGE":
		for key, val := range expense.PaidFor {
			expense.PaidFor[key] = expense.Amount * val / 100
		}

	}
	for key, val := range expense.PaidBy {
		owes, ok := expense.PaidFor[key]
		if ok {
			if val > owes {
				expense.PaidBy[key] = expense.PaidBy[key] - owes
				delete(expense.PaidFor, key)
			} else if owes > val {
				expense.PaidFor[key] = expense.PaidFor[key] - val
				delete(expense.PaidBy, key)
			} else {
				delete(expense.PaidFor, key)
				delete(expense.PaidBy, key)
			}
		}
	}

	balances = createBalances(expense)

	expense.PaidFor = paidFor
	expense.PaidBy = paidBy
	err := repository.AddExpense(expense, balances)
	if err != nil {
		return code, err
	}

	return code, nil
}

func GetExpenses() ([]domains.Expense, error) {

	return repository.GetExpenses()
}
