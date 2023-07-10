package service

import (
	"di/model"
	"di/steps"
	"di/util/errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type nodeServiceImpl struct {
	StepTypeRegistry map[string]reflect.Type
	EdgeTypeRegistry map[string]reflect.Type
}

func NewNodeService() NodeTypeService {
	return &nodeServiceImpl{
		StepTypeRegistry: initStepTypeRegistry(),
		EdgeTypeRegistry: initEdgeTypeRegistry(),
	}
}

func (nodeService *nodeServiceImpl) NewStepInstance(pipelineID uint, runID uint, stepType string, stepData model.StepData) (*steps.Step, error) {
	stepTypeStructName := nodeService.StepTypeRegistry[stepType]

	if stepTypeStructName == nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, errors.NewNotFound("stepType", stepType)
	}

	stepPtr := reflect.New(stepTypeStructName)

	setupStep := stepPtr.Interface().(steps.Step)
	setupStep.SetData(stepData)
	setupStep.SetPipelineID(pipelineID)
	setupStep.SetRunID(runID)

	return &setupStep, nil
}

func (nodeService *nodeServiceImpl) NewEdgeInstance(pipelineID uint, runID uint, edgeType string, stepConfig model.StepData) (*steps.Edge, error) {
	edgeTypeStructName := nodeService.EdgeTypeRegistry[edgeType]

	if edgeTypeStructName == nil {
		log.Printf("Unable to create new edge of type %v\n", edgeType)
		return nil, errors.NewNotFound("edgeType", edgeType)
	}

	step := reflect.New(edgeTypeStructName)
	setupEdge := step.Interface().(steps.Edge)

	return &setupEdge, nil
}

func initStepTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepo{}, steps.ScikitTestingDataset{}, steps.ScikitTrainingDataset{}}

	for _, v := range stepTypes {
		splitString := strings.SplitAfter(fmt.Sprintf("%T", v), ".")
		camelCased := strcase.ToLowerCamel(splitString[len(splitString)-1])
		stepTypeRegistry[camelCased] = reflect.TypeOf(v)
	}

	for _, v := range steps.ScikitUnsupervisedModelTypes {
		camelCased := strcase.ToLowerCamel(v)
		stepTypeRegistry[camelCased] = reflect.TypeOf(steps.ScikitUnsupervisedModel{})
	}

	return stepTypeRegistry
}

func initEdgeTypeRegistry() map[string]reflect.Type {
	var edgeTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.SmoothStep{}}

	for _, v := range stepTypes {
		edgeTypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}

	return edgeTypeRegistry
}
