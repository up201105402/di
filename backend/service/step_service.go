package service

import (
	"di/model"
	"di/steps"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type nodeServiceImpl struct {
	I18n             *i18n.Localizer
	StepTypeRegistry map[string]reflect.Type
	EdgeTypeRegistry map[string]reflect.Type
}

func NewNodeService(i18n *i18n.Localizer) StepService {
	return &nodeServiceImpl{
		I18n:             i18n,
		StepTypeRegistry: initStepTypeRegistry(),
		EdgeTypeRegistry: initEdgeTypeRegistry(),
	}
}

func (nodeService *nodeServiceImpl) NewStepInstance(pipelineID uint, runID uint, stepType string, nodeDescription model.NodeDescription) (*steps.Step, error) {
	stepTypeStructName := nodeService.StepTypeRegistry[stepType]

	if stepTypeStructName == nil {
		errMessage := nodeService.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "steps.service.step.new-instance.failed",
			TemplateData: map[string]interface{}{
				"Type": stepType,
			},
			PluralCount: 1,
		})

		log.Printf(errMessage)
		return nil, errors.New(errMessage)
	}

	stepPtr := reflect.New(stepTypeStructName)
	setupStep := stepPtr.Interface().(steps.Step)

	setupStep.SetData(nodeDescription)
	setupStep.SetPipelineID(pipelineID)
	setupStep.SetRunID(runID)

	return &setupStep, nil
}

func (nodeService *nodeServiceImpl) NewEdgeInstance(pipelineID uint, runID uint, edgeType string, nodeDescription model.NodeDescription) (*steps.Edge, error) {
	edgeTypeStructName := nodeService.EdgeTypeRegistry[edgeType]

	if edgeTypeStructName == nil {
		errMessage := nodeService.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "steps.service.edge.new-instance.failed",
			TemplateData: map[string]interface{}{
				"Type": edgeType,
			},
			PluralCount: 1,
		})

		log.Printf(errMessage)
		return nil, errors.New(errMessage)
	}

	edge := reflect.New(edgeTypeStructName)
	setupEdge := edge.Interface().(steps.Edge)
	setupEdge.SetData(nodeDescription)

	return &setupEdge, nil
}

func initStepTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepo{}, steps.ShellScript{}, steps.Dataset{}, steps.Trainer{}, steps.Tester{}, steps.Trained{}, steps.CustomPyTorchModel{}, steps.HumanFeedbackNN{}, steps.CustomHITL{}, steps.ScikitTestingDataset{}, steps.ScikitTrainingDataset{}}

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

	stepTypes := []interface{}{steps.Smoothstep{}}

	for _, v := range stepTypes {
		splitString := strings.SplitAfter(fmt.Sprintf("%T", v), ".")
		camelCased := strcase.ToLowerCamel(splitString[len(splitString)-1])
		edgeTypeRegistry[camelCased] = reflect.TypeOf(v)
	}

	return edgeTypeRegistry
}
