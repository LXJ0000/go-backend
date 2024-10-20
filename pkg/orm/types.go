package orm

import "database/sql/driver"

const (
	defaultPage = 1
	defaultSize = 10
)

type JSON[T any] struct {
	JSONValue T
	JSONValid bool
}

func (j JSON[T]) Value() (driver.Value, error) {
	return nil, nil
}

func (j *JSON[T]) Scan(src interface{}) error {
	return nil
}
