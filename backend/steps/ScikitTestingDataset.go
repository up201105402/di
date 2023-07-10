package steps

import (
	"di/model"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type ScikitTestingDataset struct {
	ID          int
	PipelineID  uint
	RunID       uint
	Dataset     string
	Name        string
	IsFirstStep bool
	DataConfig  model.StepDataConfig
}

func (step ScikitTestingDataset) GetID() int {
	return int(step.ID)
}

func (step ScikitTestingDataset) GetName() string {
	return step.Name
}

func (step *ScikitTestingDataset) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *ScikitTestingDataset) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Dataset = stepDescription.Data.NameAndType.Dataset
	step.DataConfig = stepDescription.Data.StepConfig
	return nil
}

func (step *ScikitTestingDataset) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ScikitTestingDataset) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitTestingDataset) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitTestingDataset) GetRunID() uint {
	return step.RunID
}

func (step ScikitTestingDataset) Execute(logFile *os.File) error {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"

	scikitSnippetsDir := os.Getenv("SCIKIT_SNIPPETS_DIR") + "datasets/"

	var args []string

	switch step.Dataset {
	case "scikitBreastCancer":
		args = append(args, scikitSnippetsDir+"load_breast_cancer.py")
	case "scikitDiabetes":
		args = append(args, scikitSnippetsDir+"load_diabetes.py")
	case "scikitDigits":
		args = append(args, scikitSnippetsDir+"load_digits.py")
	case "scikitIris":
		args = append(args, scikitSnippetsDir+"load_iris.py")
	case "scikitLinerrud":
		args = append(args, scikitSnippetsDir+"load_linnerrud.py")
	case "scikitWine":
		args = append(args, scikitSnippetsDir+"load_wine.py")
	case "scikitLoadFile":
		args = append(args, scikitSnippetsDir+"load_dataset_from_csv.py")
	}

	args = append(args, "-d")
	args = append(args, currentPipelineWorkDir+"filtered_testing_data.csv")

	args = append(args, "-t")
	args = append(args, currentPipelineWorkDir+"filtered_testing_target.csv")

	if step.DataConfig.LowerXRangeIndex != 0 {
		args = append(args, "-l1")
		args = append(args, string(step.DataConfig.LowerXRangeIndex))
	}

	if step.DataConfig.UpperXRangeIndex != 0 {
		args = append(args, "-u1")
		args = append(args, string(step.DataConfig.UpperXRangeIndex))
	}

	if step.DataConfig.LowerYRangeIndex != 0 {
		args = append(args, "-l2")
		args = append(args, string(step.DataConfig.LowerYRangeIndex))
	}

	if step.DataConfig.UpperYRangeIndex != 0 {
		args = append(args, "-u2")
		args = append(args, string(step.DataConfig.UpperYRangeIndex))
	}

	if step.DataConfig.DataFilePath != "" {
		args = append(args, "-f1")
		args = append(args, currentPipelineWorkDir+string(step.DataConfig.DataFilePath))
	}

	if step.DataConfig.TargetFilePath != "" {
		args = append(args, "-f2")
		args = append(args, currentPipelineWorkDir+string(step.DataConfig.TargetFilePath))
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir

	out, err := cmd.Output()
	runLogger.Println(out)

	return err
}
