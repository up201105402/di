package datasets

import (
	"di/model"
	"os"
)

type LoadFromFile struct {
}

func (step *LoadFromFile) SetConfig(stepConfig model.StepDataConfig) error {
	// TODO

	return nil
}

func (step LoadFromFile) Execute(logFile *os.File) error {
	// TODO

}
