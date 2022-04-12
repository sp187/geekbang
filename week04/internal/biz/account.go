package biz

import "github.com/sp187/geekbang/week04/internal/data"

func GetAccount(id string) (data.Account, error) {
	return data.GetAccountById(id)
}

func UpdateAccount(account data.Account) (data.Account, error) {
	return data.UpdateAccountById(account)
}
