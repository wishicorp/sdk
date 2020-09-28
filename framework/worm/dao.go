package worm

import (
	"github.com/go-xorm/xorm"
)

type Dao interface {
	Count(bean interface{}) (int64, error)
	Exists(bean interface{}) (bool, error)

	InsertOne(entry interface{}) (int64, error)
	InsertMany(entries ...interface{}) (int64, error)

	Update(bean interface{}, where ...interface{}) (int64, error)
	UpdateById(id interface{}, bean interface{}) (int64, error)

	Delete(bean interface{}) (int64, error)
	DeleteById(id interface{}, bean interface{}) (int64, error)

	FindById(id interface{}, bean interface{}) (bool, error)
	FindOne(bean interface{}) (bool, error)
	FindMany(rowsSlicePtr interface{}, orderBy string, condiBean ...interface{}) error
	FindAndCount(rowsSlicePtr interface{},
		pageable Pageable, condiBean ...interface{}) (int64, error)

	Query(rowsSlicePtr interface{}, sql string, Args ...interface{}) error

	Native() DBInterface
}

type SessionDao interface {
	Dao
	Session() *xorm.Session
	Begin() error
	Commit() error
	Rollback() error
	Close()
}

type BaseDao interface {
	Dao
	NewSession() SessionDao
}

type OrmBaseDao struct {
	conn    DBInterface
	session *xorm.Session
}

func NewBaseDao(conn DBInterface) BaseDao {
	return &OrmBaseDao{conn: conn}
}

func (m *OrmBaseDao) Session() *xorm.Session {
	return m.session
}

//创建一个session
func (m *OrmBaseDao) NewSession() SessionDao {
	m.session = m.conn.NewSession()
	return m
}

//开启事务
func (m *OrmBaseDao) Begin() error {
	return m.session.Begin()
}

func (m *OrmBaseDao) Close() {
	m.session.Close()
}

func (m *OrmBaseDao) Commit() error {
	return m.session.Commit()
}

func (m *OrmBaseDao) Rollback() error {
	return m.session.Rollback()
}

func (m *OrmBaseDao) InsertOne(entry interface{}) (int64, error) {
	if m.session != nil {
		return m.session.InsertOne(entry)
	}
	return m.conn.InsertOne(entry)
}
func (m *OrmBaseDao) InsertMany(entries ...interface{}) (int64, error) {
	if m.session != nil {
		return m.session.Insert(entries...)
	}
	return m.conn.Insert(entries...)
}

func (m *OrmBaseDao) Update(bean interface{}, where ...interface{}) (int64, error) {
	if m.session != nil {
		return m.session.Update(bean, where...)
	}
	return m.conn.Update(bean, where...)
}
func (m *OrmBaseDao) UpdateById(id interface{}, bean interface{}) (int64, error) {
	if m.session != nil {
		return m.session.ID(id).Update(bean)
	}
	return m.conn.ID(id).Update(bean)
}

func (m *OrmBaseDao) Delete(bean interface{}) (int64, error) {
	if m.session != nil {
		return m.session.Delete(bean)
	}
	return m.conn.Delete(bean)
}
func (m *OrmBaseDao) DeleteById(id interface{}, bean interface{}) (int64, error) {
	if m.session != nil {
		return m.session.ID(id).Delete(bean)
	}
	return m.conn.ID(id).Delete(bean)
}

func (m *OrmBaseDao) Query(rowsSlicePtr interface{}, sql string, Args ...interface{}) error {

	if m.session != nil {
		return m.session.SQL(sql, Args...).Find(rowsSlicePtr)
	}
	return m.conn.SQL(sql, Args...).Find(rowsSlicePtr)
}

func (m *OrmBaseDao) FindById(id interface{}, bean interface{}) (bool, error) {
	if m.session != nil {
		return m.session.ID(id).Get(bean)
	}
	return m.conn.ID(id).Get(bean)
}

func (m *OrmBaseDao) FindOne(bean interface{}) (bool, error) {
	if m.session != nil {
		return m.session.Get(bean)
	}
	return m.conn.Get(bean)
}

func (m *OrmBaseDao) Count(bean interface{}) (int64, error) {
	if m.session != nil {
		return m.session.Count(bean)
	}
	return m.conn.Count(bean)
}

func (m *OrmBaseDao) Exists(bean interface{}) (bool, error) {
	if m.session != nil {
		return m.session.Exist(bean)
	}
	return m.conn.Exist(bean)
}
func (m *OrmBaseDao) FindMany(rowsSlicePtr interface{}, sort string, condiBean ...interface{}) error {
	if m.session != nil {
		return m.session.OrderBy(sort).Find(rowsSlicePtr, condiBean...)
	}
	return m.conn.OrderBy(sort).Find(rowsSlicePtr, condiBean...)
}

func (m *OrmBaseDao) FindAndCount(rowsSlicePtr interface{},
	pageable Pageable, condiBean ...interface{}) (int64, error) {

	if m.session != nil {
		return m.session.
			Limit(pageable.Limit(), pageable.Skip()).
			OrderBy(pageable.Sort()).
			FindAndCount(rowsSlicePtr, condiBean...)
	}
	return m.conn.
		Limit(pageable.Limit(), pageable.Skip()).
		OrderBy(pageable.Sort()).
		FindAndCount(rowsSlicePtr, condiBean...)
}

func (m *OrmBaseDao) Native() DBInterface {
	return m.conn
}
