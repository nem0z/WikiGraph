package database

type Graph interface {
	CreateEdge(name string, fromKey string, toKey string)
}
