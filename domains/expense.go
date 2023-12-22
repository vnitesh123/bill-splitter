package domains

type Expense struct {
	CreatedBy   int     `json:"createdBy"`
	SplitType   string  `json:"splitType"`
	Amount      float64 `json:"amount"`
	CreatedOn   string
	ExpenseName string             `json:"name"`
	PaidBy      map[string]float64 `json:"paidBy"`
	PaidFor     map[string]float64 `json:"paidFor"`
	ExpenseID   int                `json:"expenseId"`
}

type ExpenseDBEntry struct {
	ExpenseID   int     `db:"id"`
	SplitType   string  `db:"split"`
	Amount      float64 `db:"amount"`
	CreatedOn   string  `db:"createdOn"`
	ExpenseName string  `db:"name"`
	PaidBy      string  `db:"paidBy"`
	PaidFor     string  `db:"paidFor"`
	CreatedBy   int     `db:"createdBy"`
}
