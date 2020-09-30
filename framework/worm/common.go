package worm

import (
	"errors"
	"github.com/go-xorm/xorm"
)

var ErrGetEmpty = errors.New("found 0 rows")
var ErrUpdatedEmpty = errors.New("update affected 0 rows")
var ErrDeletedEmpty = errors.New("update affected 0 rows")
var ErrInsertedEmpty = errors.New("update affected 0 rows")

func DoGet(call func() (bool, error)) error {
	has, err := call()
	if nil != err {
		return err
	}
	if !has {
		return ErrGetEmpty
	}
	return nil
}
func DoUpdate(call func() (int64, error)) error {
	ret, err := call()
	if nil != err {
		return err
	}
	if 0 == ret {
		return ErrUpdatedEmpty
	}
	return nil
}

func DoInsert(call func() (int64, error)) error {
	ret, err := call()
	if nil != err {
		return err
	}
	if 0 == ret {
		return ErrInsertedEmpty
	}
	return nil
}

func DoDelete(call func() (int64, error)) error {
	ret, err := call()
	if nil != err {
		return err
	}
	if 0 == ret {
		return ErrDeletedEmpty
	}
	return nil
}

type SessionDoctor int

const (
	SessionDoctorCommit   SessionDoctor = 0
	SessionDoctorRollback SessionDoctor = 1
)

type SessionHandler func(session *xorm.Session) (interface{}, SessionDoctor, error)
type SessionWrapper struct {
	Session *xorm.Session
	Handler SessionHandler
}

//事务装饰器
//error 为真事务回滚
//SessionDoctor 决定回滚还是提交
func WarpSession(wrapper SessionWrapper) (interface{}, int, error) {
	defer wrapper.Session.Close()
	if err := wrapper.Session.Begin(); err != nil {
		return nil, 0, err
	}
	value, doctor, err := wrapper.Handler(wrapper.Session)
	if err != nil {
		return nil, 0, err
	}

	if doctor == SessionDoctorCommit {
		if err2 := wrapper.Session.Commit(); err2 != nil {
			return nil, 0, err2
		}
	} else {
		if err2 := wrapper.Session.Rollback(); err2 != nil {
			return nil, 0, err2
		}
	}

	return value, 0, nil
}
