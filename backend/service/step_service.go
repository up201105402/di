package service

import (
	"di/model"
	"di/steps"
	"di/util/errors"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type stepServiceImpl struct {
	StepTypeRegistry map[string]reflect.Type
}

func NewStepService() StepTypeService {
	return &stepServiceImpl{
		StepTypeRegistry: initStepTypeRegistry(),
	}
}

func (stepService *stepServiceImpl) NewStepInstance(stepType string, stepConfig model.StepDataConfig) (*model.Step, error) {
	stepTypeStructName := stepService.StepTypeRegistry[stepType]

	if stepTypeStructName == nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, errors.NewNotFound("stepType", stepType)
	}

	step := reflect.New(stepTypeStructName).Elem().Interface().(model.Step)

	marshalledStepDataConfig, err := json.Marshal(stepConfig)

	if err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	if err = json.Unmarshal(marshalledStepDataConfig, &step); err != nil {
		log.Printf("Unable to create new step of type %v\n", stepType)
		return nil, err
	}

	return &step, nil
}

func initStepTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepoStep{}, steps.SmoothStepStep{}}

	for _, v := range stepTypes {
		stepTypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}

	return stepTypeRegistry
}
