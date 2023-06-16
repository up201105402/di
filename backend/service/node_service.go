package service

import (
	"di/model"
	"di/steps"
	"di/util/errors"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type stepServiceImpl struct {
	StepTypeRegistry map[string]reflect.Type
	EdgeTypeRegistry map[string]reflect.Type
}

func NewStepService() NodeTypeService {
	return &stepServiceImpl{
		StepTypeRegistry: initStepTypeRegistry(),
		EdgeTypeRegistry: initEdgeTypeRegistry(),
	}
}

func (stepService *stepServiceImpl) NewStepInstance(pipelineID uint, stepType string, stepConfig model.StepDataConfig) (*steps.Step, error) {
	stepTypeStructName := stepService.StepTypeRegistry[stepType]

	if stepTypeStructName == nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, errors.NewNotFound("stepType", stepType)
	}

	step := reflect.New(stepTypeStructName).Elem()

	setStepPipelineID(step, pipelineID)
	setStepFields(step, stepConfig)

	setupStep := step.Interface().(steps.Step)
	marshalledStepDataConfig, err := json.Marshal(stepConfig)

	if err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	if err = json.Unmarshal(marshalledStepDataConfig, &setupStep); err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	return &setupStep, nil
}

func setStepPipelineID(stp reflect.Value, pipelineID uint) {
	t := reflect.ValueOf(stp)
	pipelineIDField := t.FieldByName("PipelineID")

	if pipelineIDField.IsValid() {
		if pipelineIDField.CanSet() {
			if pipelineIDField.Kind() == reflect.Uint {
				pipelineIDField.SetUint(uint64(pipelineID))
			}
		}
	}
}

func setStepFields(stp reflect.Value, stepConfig model.StepDataConfig) {
	t := reflect.TypeOf(stp)
	d := reflect.ValueOf(stepConfig)

	for i := 0; i < t.NumField(); i++ {
		for j := 0; j < d.NumField(); j++ {
			if t.Field(i).Name == d.Field(j).Type().Name() {
				v := reflect.ValueOf(stp).Elem()
				field := v.Field(i)

				if field.IsValid() {
					if field.CanSet() {
						field.Set(reflect.ValueOf(d.Field(j)))
					}
				}
			}
		}
	}
}

func (stepService *stepServiceImpl) NewEdgeInstance(pipelineID uint, stepType string, stepConfig model.StepDataConfig) (*steps.Edge, error) {
	stepTypeStructName := stepService.StepTypeRegistry[stepType]

	if stepTypeStructName == nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, errors.NewNotFound("stepType", stepType)
	}

	step := reflect.New(stepTypeStructName).Elem()

	pipelineIDField := step.FieldByName("PipelineID")

	if pipelineIDField.IsValid() {
		if pipelineIDField.CanSet() {
			if pipelineIDField.Kind() == reflect.Uint {
				pipelineIDField.SetUint(uint64(pipelineID))
			}
		}
	}

	setupEdge := step.Interface().(steps.Edge)
	marshalledStepDataConfig, err := json.Marshal(stepConfig)

	if err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	if err = json.Unmarshal(marshalledStepDataConfig, &setupEdge); err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	return &setupEdge, nil
}

func initStepTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepo{}}

	for _, v := range stepTypes {
		splitString := strings.SplitAfter(fmt.Sprintf("%T", v), ".")
		stepTypeRegistry[splitString[len(splitString)-1]] = reflect.TypeOf(v)
	}

	return stepTypeRegistry
}

func initEdgeTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.SmoothStep{}}

	for _, v := range stepTypes {
		stepTypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}

	return stepTypeRegistry
}
