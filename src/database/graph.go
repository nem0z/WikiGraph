package database

type Graph interface {
	CreateEdge(parent string, child string) error
}
