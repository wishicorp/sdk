package worm

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"xorm.io/core"
)

type SQLLogger struct {
	log     hclog.Logger
	showSql bool
}

func NewSQLLogger(log hclog.Logger, showSql bool) *SQLLogger {
	namedLog := log.Named("mysql")
	return &SQLLogger{
		log:     namedLog,
		showSql: showSql,
	}
}

func (s *SQLLogger) Debug(v ...interface{}) {
	s.log.Debug("", v...)
}

func (s *SQLLogger) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log.Debug(msg)
}

func (s *SQLLogger) Error(v ...interface{}) {
	s.log.Error("", v...)
}

func (s *SQLLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log.Error(msg)
}

func (s *SQLLogger) Info(v ...interface{}) {
	s.log.Error("", v...)
}

func (s *SQLLogger) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log.Info(msg)
}

func (s *SQLLogger) Warn(v ...interface{}) {
	s.log.Error("", v...)
}

func (s *SQLLogger) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	s.log.Warn(msg)
}

func (s *SQLLogger) Level() core.LogLevel {
	if s.log.IsDebug() || s.log.IsTrace() {
		return core.LOG_DEBUG
	}
	if s.log.IsInfo() {
		return core.LOG_INFO
	}
	if s.log.IsError() {
		return core.LOG_ERR
	}
	if s.log.IsWarn() {
		return core.LOG_WARNING
	}
	return core.LOG_UNKNOWN
}

func (s *SQLLogger) SetLevel(l core.LogLevel) {
	switch l {
	case core.LOG_INFO:
		s.log.SetLevel(hclog.Info)
	case core.LOG_WARNING:
		s.log.SetLevel(hclog.Warn)
	case core.LOG_DEBUG:
		s.log.SetLevel(hclog.Trace)
	case core.LOG_ERR:
		s.log.SetLevel(hclog.Error)
	default:
		s.log.SetLevel(hclog.NoLevel)
	}

}

func (s *SQLLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		s.showSql = show[0]
	}
}

func (s *SQLLogger) IsShowSQL() bool {
	return s.showSql
}
