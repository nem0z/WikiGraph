package entity

type Link struct {
	Parent string
	Child  string
}

func NewLink(parent, child string) *Link {
	return &Link{Parent: parent, Child: child}
}
