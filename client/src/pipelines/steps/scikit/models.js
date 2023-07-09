import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, customDelay, removeDuplicates, golangType } from '@/util';
import { getNode, createMessage } from '@formkit/core';

import { leastSquaresStepConfig } from './leastSquares';
import { ridgeRegressionStepConfig, ridgeRegressionCVStepConfig } from './ridgeRegression';
import { ridgeClassifierStepConfig, ridgeClassifierCVStepConfig } from './ridgeClassifier';
import { lassoStepConfig, lassoCVStepConfig, lassoLarsStepConfig, lassoLarsCVStepConfig, lassoLarsICStepConfig, multiTaskLassoStepConfig, multiTaskLassoCVStepConfig } from './lasso';
import { elasticNetStepConfig, elasticNetCVStepConfig, multiTaskElasticNetStepConfig, multiTaskElasticNetCVStepConfig } from './elastic';
import { larsStepConfig, larsCVStepConfig } from './lars';
import { ompStepConfig, ompCVStepConfig } from './omp';
import { bayesianRidgeStepConfig, bayesianARDStepConfig } from './bayesian';
import { logisticRegressionStepConfig, logisticRegressionCVStepConfig } from './logistic';
import { tweedieRegressorStepConfig, poissonRegressorStepConfig, gammaRegressorStepConfig } from './generalizedLinears';
import { sgdClassifierStepConfig, sgdRegressorStepConfig } from './sgd';
import { perceptronStepConfig } from './perceptron';
import { passiveAgressiveClassifierStepConfig, passiveAgressiveRegressorStepConfig } from './passiveAgressive';
import { huberRegressorStepConfig, ransacRegressorStepConfig, theilSenRegressorStepConfig } from './robustnessRegression';
import { quantileRegressionStepConfig } from './quantileRegression';

const scikitUnsupervisedModelsForm = (stepConfigFields) => function (data, onSubmit) {
    const activeStep = ref('');
    const steps = reactive({});
    const visitedSteps = ref([]); // track visited steps

    // watch our activeStep and store visited steps
    // to know when to show errors
    watch(activeStep, (newStep, oldStep) => {
        if (oldStep && !visitedSteps.value.includes(oldStep)) {
            visitedSteps.value.push(oldStep)
        }
        // trigger showing validation on fields
        // within all visited steps
        visitedSteps.value.forEach((step) => {
            const node = getNode(step)
            node.walk((n) => {
                n.store.set(
                    createMessage({
                        key: 'submitted',
                        value: true,
                        visible: false,
                    })
                )
            })
        })
    })

    const setStep = (delta) => {
        const stepNames = Object.keys(steps);
        const currentIndex = stepNames.indexOf(activeStep.value);
        activeStep.value = stepNames[currentIndex + delta];
    }

    // pushes the steps (group nodes - $formkit: 'group') into the steps object
    const stepPlugin = (node) => {
        if (node.props.type == "group") {
            // builds an object of the top-level groups
            steps[node.name] = steps[node.name] || {}

            node.on('created', () => {
                // use 'on created' to ensure context object is available
                steps[node.name].valid = toRef(node.context.state, 'valid')
            })

            // listen for changes in error count and store it
            node.on('count:errors', ({ payload: count }) => {
                steps[node.name].errorCount = count
            })

            // listen for changes in count of blocking validations messages
            node.on('count:blocking', ({ payload: count }) => {
                steps[node.name].blockingCount = count
            })

            // set the active tab to the 1st tab
            if (activeStep.value === '') {
                activeStep.value = node.name
            }

            // Stop plugin inheritance to descendant nodes
            return false
        }
    }

    const formkitData = reactive({
        steps,
        visitedSteps,
        activeStep,
        plugins: [
            stepPlugin
        ],
        setStep: target => () => {
            setStep(target)
        },
        setActiveStep: stepName => () => {
            activeStep.value = stepName
        },
        showStepErrors: stepName => {
            return (steps[stepName].errorCount > 0 || steps[stepName].blockingCount > 0) && (visitedSteps.value && visitedSteps.value.includes(stepName))
        },
        stepIsValid: stepName => {
            return steps[stepName].valid && steps[stepName].errorCount === 0
        },
        submitForm: async (formData, node) => {
            try {
                await customDelay(formData);
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
            }
        },
        stringify: (value) => JSON.stringify(value, null, 2),
        camel2title
    })

    const formSchema = [
        {
            $cmp: 'FormKit',
            props: {
                type: 'form',
                id: 'form',
                onSubmit: '$submitForm',
                plugins: '$plugins',
                actions: false,
                value: { ...data }
            },
            children: [
                {
                    $el: 'ul',
                    attrs: {
                        class: "steps"
                    },
                    children: [
                        {
                            $el: 'li',
                            for: ['step', 'stepName', '$steps'],
                            attrs: {
                                class: {
                                    'step': true,
                                    'has-errors': '$showStepErrors($stepName)'
                                },
                                style: {
                                    if: '$activeNodeType == ""',
                                    then: 'display: none;'
                                },
                                onClick: '$setActiveStep($stepName)',
                                'data-step-active': '$activeStep === $stepName',
                                'data-step-valid': '$stepIsValid($stepName)'
                            },
                            children: [
                                {
                                    $el: 'span',
                                    if: '$showStepErrors($stepName)',
                                    attrs: {
                                        class: 'step--errors'
                                    },
                                    children: '$step.errorCount + $step.blockingCount'
                                },
                                '$camel2title($stepName)'
                            ]
                        }
                    ]
                },
                {
                    $el: 'div',
                    attrs: {
                        class: 'form-body'
                    },
                    children: [
                        {
                            $el: 'section',
                            attrs: {
                                style: {
                                    if: '$activeStep !== "name"',
                                    then: 'display: none;'
                                }
                            },
                            children: [
                                {
                                    $formkit: 'group',
                                    id: 'name',
                                    name: 'name',
                                    children: [
                                        {
                                            $formkit: 'text',
                                            name: 'nodeName',
                                            label: 'Step Name',
                                            placeholder: 'Step Name',
                                            validation: 'required'
                                        },
                                    ]
                                }
                            ]
                        },
                        {
                            $el: 'section',
                            attrs: {
                                style: {
                                    if: '$activeStep !== "stepConfig"',
                                    then: 'display: none;'
                                }
                            },
                            children: [
                                {
                                    $formkit: 'group',
                                    id: 'stepConfig',
                                    name: 'stepConfig',
                                    children: stepConfigFields,
                                }
                            ]
                        },
                        {
                            $el: 'div',
                            attrs: {
                                class: 'step-nav'
                            },
                            children: [
                                {
                                    $formkit: 'button',
                                    disabled: '$activeStep === "name"',
                                    onClick: '$setStep(-1)',
                                    children: 'Back'
                                },
                                {
                                    $formkit: 'button',
                                    disabled: '$activeStep === "stepConfig"',
                                    onClick: '$setStep(1)',
                                    children: 'Next'
                                }
                            ]
                        },
                    ]
                },
                {
                    $el: 'div',
                    attrs: {
                        class: 'formkit-bottom-buttons'
                    },
                    children: [
                        {
                            $formkit: 'button',
                            label: 'Cancel',
                            id: 'cancel-create-step-button'
                        },
                        {
                            $formkit: 'submit',
                            label: 'Submit',
                            disabled: '$get(form).state.valid !== true',
                        },
                    ]
                },
            ]
        },
    ];

    return { formSchema, formkitData }
};

export const scikitUnsupervisedModels = {
    leastSquares: scikitUnsupervisedModelsForm(leastSquaresStepConfig),
    ridgeRegression: scikitUnsupervisedModelsForm(ridgeRegressionStepConfig),
    ridgeRegressionCV: scikitUnsupervisedModelsForm(ridgeRegressionCVStepConfig),
    ridgeClassifier: scikitUnsupervisedModelsForm(ridgeClassifierStepConfig),
    ridgeClassifierCV: scikitUnsupervisedModelsForm(ridgeClassifierCVStepConfig),
    lasso: scikitUnsupervisedModelsForm(lassoStepConfig),
    lassoCV: scikitUnsupervisedModelsForm(lassoCVStepConfig),
    lassoLars: scikitUnsupervisedModelsForm(lassoLarsStepConfig),
    lassoLarsCV: scikitUnsupervisedModelsForm(lassoLarsCVStepConfig),
    lassoLarsIC: scikitUnsupervisedModelsForm(lassoLarsICStepConfig),
    multiTaskLasso: scikitUnsupervisedModelsForm(multiTaskLassoStepConfig),
    multiTaskLassoCV: scikitUnsupervisedModelsForm(multiTaskLassoCVStepConfig),
    elasticNet: scikitUnsupervisedModelsForm(elasticNetStepConfig),
    elasticNetCV: scikitUnsupervisedModelsForm(elasticNetCVStepConfig),
    multiTaskElasticNet: scikitUnsupervisedModelsForm(multiTaskElasticNetStepConfig),
    multiTaskElasticNetCV: scikitUnsupervisedModelsForm(multiTaskElasticNetCVStepConfig),
    lars: scikitUnsupervisedModelsForm(larsStepConfig),
    larsCV: scikitUnsupervisedModelsForm(larsCVStepConfig),
    omp: scikitUnsupervisedModelsForm(ompStepConfig),
    ompCV: scikitUnsupervisedModelsForm(ompCVStepConfig),
    bayesianRidge: scikitUnsupervisedModelsForm(bayesianRidgeStepConfig),
    bayesianARD: scikitUnsupervisedModelsForm(bayesianARDStepConfig),
    logisticRegression: scikitUnsupervisedModelsForm(logisticRegressionStepConfig),
    logisticRegressionCV: scikitUnsupervisedModelsForm(logisticRegressionCVStepConfig),
    tweedieRegressor: scikitUnsupervisedModelsForm(tweedieRegressorStepConfig),
    poissonRegressor: scikitUnsupervisedModelsForm(poissonRegressorStepConfig),
    gammaRegressor: scikitUnsupervisedModelsForm(gammaRegressorStepConfig),
    sgdClassifier: scikitUnsupervisedModelsForm(sgdClassifierStepConfig),
    sgdRegressor: scikitUnsupervisedModelsForm(sgdRegressorStepConfig),
    perceptron: scikitUnsupervisedModelsForm(perceptronStepConfig),
    passiveAgressiveClassifier: scikitUnsupervisedModelsForm(passiveAgressiveClassifierStepConfig),
    passiveAgressiveRegressor: scikitUnsupervisedModelsForm(passiveAgressiveRegressorStepConfig),
    huberRegression: scikitUnsupervisedModelsForm(huberRegressorStepConfig),
    ransacRegression: scikitUnsupervisedModelsForm(ransacRegressorStepConfig),
    theilSenRegression: scikitUnsupervisedModelsForm(theilSenRegressorStepConfig),
    quantileRegression: scikitUnsupervisedModelsForm(quantileRegressionStepConfig),
}
