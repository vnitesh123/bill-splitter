package domains

type User struct {
	UserId   int
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	AddedOn  string
}

type Message struct {
	Message string `json:"message"`
}

type Balances struct {
	Owes map[int]float64 `json:"owes"`
	Gets map[int]float64 `json:"gets"`
}

type OwesBalances struct {
	OwesTo int     `db:"paidBy"`
	Amount float64 `db:"amount"`
}

type GetsBalances struct {
	GetsFrom int     `db:"paidFor"`
	Amount   float64 `db:"amount"`
}
