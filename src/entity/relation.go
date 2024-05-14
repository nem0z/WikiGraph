package entity

type Relation struct {
	ParentLink string     `json:"parentLink"`
	Childs     []*Article `json:"childs"`
}

func NewRelation(parentLink string, childs ...*Article) *Relation {
	return &Relation{ParentLink: parentLink, Childs: childs}
}
