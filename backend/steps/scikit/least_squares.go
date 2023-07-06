package scikit

import (
	"di/model"
	"os"
)

type LeastSquares struct {
}

func (step *LeastSquares) SetConfig(stepConfig model.StepDataConfig) error {
	// TODO

	return nil
}

func (step LeastSquares) Execute(logFile *os.File) error {
	// TODO

	// python.Initialize()
	// defer python.Finalize()

	// // Import Python code (foo.py)
	// foo, _ := python.Import("foo")
	// defer foo.Release()

	// // Get access to a Python function
	// hello, _ := foo.GetAttr("hello")
	// defer hello.Release()

	return nil
}
