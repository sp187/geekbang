package data

import (
	"database/sql"
	"errors"

	xerrors "github.com/pkg/errors"
)

const (
	qAccountById       = "select id, name, age, description from account where id = $1;"
	qUpdateAccountById = "update account set name = $2, age = $3, description = $4 where id = $1;"
)

type Account struct {
	Id          string
	Name        string
	Age         int
	Description string
}

func GetAccountById(id string) (Account, error) {
	account := Account{}
	err := GetDB().QueryRow(qAccountById, id).Scan(
		&account.Id,
		&account.Name,
		&account.Age,
		&account.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, xerrors.New("account does not exit")
		} else {
			return account, xerrors.Wrap(err, "query account information fail")
		}
	}
	return account, nil
}

func UpdateAccountById(account Account) (Account, error) {
	results, err := GetDB().Exec(qUpdateAccountById, account.Id, account.Name, account.Age, account.Description)
	if err != nil {
		return account, xerrors.Wrap(err, "update account fail")
	}
	if affectRow, _ := results.RowsAffected(); affectRow != 1 {
		return account, xerrors.New("account does not exit")
	}
	return account, nil
}
