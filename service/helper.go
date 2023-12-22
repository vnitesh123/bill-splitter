package service

import (
	"splitwise/domains"
)

func createBalances(expense domains.Expense) map[string]map[string]float64 {
	proportions := make(map[string]float64)
	balances := make(map[string]map[string]float64)
	var total float64
	total = 0
	for _, val := range expense.PaidBy {
		total = total + val
	}

	for key, val := range expense.PaidBy {
		proportions[key] = val / total
	}

	for user := range expense.PaidBy {
		perUserBalance := make(map[string]float64)
		for payee, val := range expense.PaidFor {
			perUserBalance[payee] = proportions[user] * val
		}
		balances[user] = perUserBalance
	}

	return balances
}
