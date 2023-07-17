import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, i18nFromStepName } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { i18n } from '@/i18n';

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
import { stepTabs, cancelAndSubmitButtons, getFormBody } from '@/pipelines/steps/formBasics';

const { t } = i18n.global;

const nameAndTypeGroupChildren = [
    {
        $formkit: 'group',
        id: 'nameAndType',
        name: 'nameAndType',
        children: [
            {
                $formkit: 'text',
                name: 'name',
                label: t('pages.pipelines.edit.dialog.nameAndType.name'),
                placeholder: t('pages.pipelines.edit.dialog.nameAndType.name'),
                validation: 'required'
            },
            {
                $formkit: 'checkbox',
                name: 'isFirstStep',
                label: t('pages.pipelines.edit.dialog.nameAndType.isFirstLabel'),
                if: '$showIsFirstStep == true',
            },
        ]
    }
];

const scikitUnsupervisedModelsForm = (stepConfigGroupChildren) => function (data, onSubmit) {
    const activeStep = ref('');
    const steps = reactive({});
    const visitedSteps = ref([]); // track visited steps
    const showIsFirstStep = data.nameAndType == null || data?.nameAndType?.isFirstStep == false;

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
        showIsFirstStep,
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
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
            }
        },
        stringify: (value) => JSON.stringify(value, null, 2),
        camel2title,
        i18nFromStepName
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
                stepTabs,
                getFormBody(nameAndTypeGroupChildren, stepConfigGroupChildren),
                cancelAndSubmitButtons,
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
