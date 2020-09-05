package extension

type Action string

const (
	ActionLint  Action = "lint"
	ActionDiff  Action = "diff"
	ActionApply Action = "apply"
)

func (a Action) String() string {
	return string(a)
}

func AllowedActions() []string {
	return []string{
		ActionLint.String(),
		ActionDiff.String(),
		ActionApply.String(),
	}
}
