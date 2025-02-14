package entity

// Interface is a sets of function need to be implemented by each entity.
type Interface interface {
	TableName() string
	FilterableFields() []interface{}
	TimeFields() []interface{}
}
