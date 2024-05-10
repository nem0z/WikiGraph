package entity

type Link struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
}

func NewLink(parent, child string) *Link {
	return &Link{Parent: parent, Child: child}
}
