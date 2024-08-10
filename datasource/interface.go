package datasource

// if exist, return value
// if not-exist, return nil
type DataSource interface {
	Value(key string) (interface{}, error)
}

type nonExist int

type Cache interface {
	// if not found, return nil and nil
	// if found, return value and nil
	// if already know not exist, return nonExist type and nil
	Value(key string) (interface{}, error)
	Store(key string, value interface{}) error
}
