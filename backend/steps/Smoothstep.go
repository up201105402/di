package steps

import (
	"di/model"
	"strconv"
)

type Smoothstep struct {
	SourceID int `json:"source"`
	TargetID int `json:"target"`
}

func (step *Smoothstep) SetData(stepDescription model.NodeDescription) {
	step.SourceID, _ = strconv.Atoi(stepDescription.SourceID)
	step.TargetID, _ = strconv.Atoi(stepDescription.TargetID)
}

func (step *Smoothstep) GetSourceID() int {
	return step.SourceID
}

func (step *Smoothstep) GetTargetID() int {
	return step.TargetID
}
