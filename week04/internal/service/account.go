package service

import (
	"encoding/json"
	"errors"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/biz"
	"gitlab.bj.sensetime.com/sense-remote/project/geekbang/week04/internal/data"
	"net/http"
)

type AccountReq struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}

func (aq *AccountReq) ToDO() data.Account {
	return data.Account{
		Id:          aq.Id,
		Name:        aq.Name,
		Age:         aq.Age,
		Description: aq.Description,
	}
}

type AccountResp = AccountReq

func ToDTO(account data.Account) AccountResp {
	return AccountResp{
		Id:          account.Id,
		Name:        account.Name,
		Age:         account.Age,
		Description: account.Description,
	}
}

func BuildErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func BuildOKResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&data)
}

func GetAccount(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	id := queryParams.Get("id")
	account, err := biz.GetAccount(id)
	if err != nil {
		if err.Error() == "account does not exit" {
			BuildErrorResponse(w, err, http.StatusNotFound)
			return
		} else {
			BuildErrorResponse(w, errors.New("服务器内部错误"), http.StatusInternalServerError)
			return
		}
	}
	BuildOKResponse(w, ToDTO(account))
	return
}

func UpdateAccount(w http.ResponseWriter, req *http.Request) {
	request := AccountReq{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		BuildErrorResponse(w, errors.New("参数非法"), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	account, err := biz.UpdateAccount(request.ToDO())
	if err != nil {
		if err.Error() == "account does not exit" {
			BuildErrorResponse(w, err, http.StatusNotFound)
			return
		} else {
			BuildErrorResponse(w, errors.New("服务器内部错误"), http.StatusInternalServerError)
			return
		}
	}
	BuildOKResponse(w, ToDTO(account))
	return
}
