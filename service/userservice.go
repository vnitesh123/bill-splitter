package service

import (
	"splitwise/domains"
	"splitwise/repository"
)

func RegisterUser(user domains.User) error {
	return repository.AddUser(user)
}

func GetBalances(userId int) (*domains.Balances, error) {
	balances, err := repository.GetBalances(userId)
	if err != nil {
		return nil, err
	}

	for key, val := range balances.Owes {
		gets, ok := balances.Gets[key]
		if ok {
			if val > gets {
				balances.Owes[key] = val - gets
				delete(balances.Gets, key)
			} else if gets > val {
				balances.Gets[key] = gets - val
				delete(balances.Owes, key)
			} else {
				delete(balances.Owes, key)
				delete(balances.Gets, key)
			}
		}
	}
	return balances, nil
}
