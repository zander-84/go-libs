package mysql

type Mysql interface {
	Engine() interface{}
	Transaction(f func(tx interface{}) (int, error)) (int, error)
}
