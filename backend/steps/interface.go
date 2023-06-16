package steps

type Step interface {
	GetID() int
	Execute() error
}

type Edge interface {
	GetID() int
	GetNextStep() *Step
	GetPreviousStep() *Step
}
