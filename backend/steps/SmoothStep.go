package steps

type SmoothStepStep struct {
	Source uint `json:"source"`
	Target uint `json:"target"`
}

func (step *SmoothStepStep) Execute() {

}
