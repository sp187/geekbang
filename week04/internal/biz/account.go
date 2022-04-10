package biz

import "gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/data"

func GetAccount(id string) (data.Account, error) {
	return data.GetAccountById(id)
}

func UpdateAccount(account data.Account) (data.Account, error) {
	return data.UpdateAccountById(account)
}
