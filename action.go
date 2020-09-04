package main

type Action string

const (
	ActionLint  Action = "lint"
	ActionDiff  Action = "diff"
	ActionApply Action = "apply"
)
