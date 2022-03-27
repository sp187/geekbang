package week02

import (
	"database/sql"
	"errors"
	xerrors "github.com/pkg/errors"
)

// 1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？


// 个人认为需要根据业务的场景进行区分，例如想获取某个账户的信息，但账号刚被号主删除了，此时进行降级处理，不将其视为错误。
func SituationOne() error {
	_, err := DummyQueryUserInfo("geekbang")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// do something，返回账号不存在，或者已删除之类的信息
			return nil
		} else {
			return xerrors.Wrap(err, "query user information fail")
		}
	}
	return nil
}


// 如果sql.ErrNoRows完全属于意外情况，无法进行降级处理，则需要将错误打包抛回上层
func SituationTwo() error {
	_, err := DummyQueryUserInfo("geekbang")
	if err != nil {
		return xerrors.Wrap(err, "query user information fail")
	}
	return nil
}

func DummyQueryUserInfo(id string) (interface{}, error) {
	return nil, nil
}
