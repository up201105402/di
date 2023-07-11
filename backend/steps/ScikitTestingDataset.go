package steps

import (
	"bytes"
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

	if step.DataConfig.LowerXRangeIndex.Valid {
		args = append(args, "-l1")
		args = append(args, string(step.DataConfig.LowerXRangeIndex.Int64))
	}

	if step.DataConfig.UpperXRangeIndex.Valid {
		args = append(args, "-u1")
		args = append(args, string(step.DataConfig.UpperXRangeIndex.Int64))
	}

	if step.DataConfig.LowerYRangeIndex.Valid {
		args = append(args, "-l2")
		args = append(args, string(step.DataConfig.LowerYRangeIndex.Int64))
	}

	if step.DataConfig.UpperYRangeIndex.Valid {
		args = append(args, "-u2")
		args = append(args, string(step.DataConfig.UpperYRangeIndex.Int64))
	}

	if step.DataConfig.DataFilePath.Valid {
		args = append(args, "-f1")
		args = append(args, currentPipelineWorkDir+string(step.DataConfig.DataFilePath.String))
	}

	if step.DataConfig.TargetFilePath.Valid {
		args = append(args, "-f2")
		args = append(args, currentPipelineWorkDir+string(step.DataConfig.TargetFilePath.String))
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	runLogger.Println(stderr.String())
	runLogger.Println(stdout.String())

	return err
}
