package runtime

import (
	"errors"

	"gorm.io/gorm"
)

func (e *Environment) MainDB() (*gorm.DB, error) {
	if e == nil || e.dataSources == nil {
		return nil, errors.New("trun: 运行环境尚未初始化")
	}
	return e.dataSources.Main()
}

func (e *Environment) DB(name string) (*gorm.DB, error) {
	if e == nil || e.dataSources == nil {
		return nil, errors.New("trun: 运行环境尚未初始化")
	}
	return e.dataSources.Resolve(name)
}
