package service

import (
	"di/model"
	"di/steps"
	"fmt"
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

func (stepService *stepServiceImpl) NewStepInstance(stepType string) (*model.Step, error) {
	step := reflect.New(stepService.StepTypeRegistry[stepType]).Elem().Interface().(model.Step)
	return &step, nil
}

func initStepTypeRegistry() map[string]reflect.Type {
	var stepTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepoStep{}}

	for _, v := range stepTypes {
		stepTypeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	}

	return stepTypeRegistry
}
