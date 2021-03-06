package xorm

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

const (
	version string = "0.1.9"
)

// new a db manager according to the parameter. Currently support three
// driver
func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	engine := &Engine{ShowSQL: false, DriverName: driverName, Mapper: SnakeMapper{},
		DataSourceName: dataSourceName}

	engine.Tables = make(map[reflect.Type]*Table)
	engine.mutex = &sync.Mutex{}
	engine.TagIdentifier = "xorm"
	engine.Filters = make([]Filter, 0)
	if driverName == SQLITE {
		engine.Dialect = &sqlite3{}
	} else if driverName == MYSQL {
		engine.Dialect = &mysql{}
	} else if driverName == POSTGRES {
		engine.Dialect = &postgres{}
		engine.Filters = append(engine.Filters, &PgSeqFilter{})
		engine.Filters = append(engine.Filters, &PgQuoteFilter{})
	} else if driverName == MYMYSQL {
		engine.Dialect = &mysql{}
	} else {
		return nil, errors.New(fmt.Sprintf("Unsupported driver name: %v", driverName))
	}
	engine.Filters = append(engine.Filters, &IdFilter{})
	engine.Logger = os.Stdout

	//engine.Pool = NewSimpleConnectPool()
	//engine.Pool = NewNoneConnectPool()
	//engine.Cacher = NewLRUCacher()
	err := engine.SetPool(NewSysConnectPool())

	return engine, err
}
