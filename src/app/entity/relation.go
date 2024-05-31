package entity

type Relation struct {
	ParentLink string     `json:"parentLink"`
	Childs     []*Article `json:"childs"`
}

func NewRelation(parentLink string, childs ...*Article) *Relation {
	return &Relation{ParentLink: parentLink, Childs: childs}
}

type ResolvedRelation struct {
	ParentId int64   `json:"parent_id"`
	ChildIds []int64 `json:"child_ids"`
}

func NewResolvedRelation(parentId int64, childIds ...int64) *ResolvedRelation {
	return &ResolvedRelation{ParentId: parentId, ChildIds: childIds}
}
