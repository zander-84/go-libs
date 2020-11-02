package mongo

type Mongo interface {
	Engine() interface{}
	DB() interface{} //*Database
}
