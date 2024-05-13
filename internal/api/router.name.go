package api

type Action string

const (
	CreateACLAction Action = "action:CreateACL"
	ListACLAction   Action = "action:ListACL"
)

func (a Action) String() string {
	return string(a)
}
