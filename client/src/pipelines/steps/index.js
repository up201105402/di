import { markRaw } from 'vue';
import TrainModelNode from './components/nodes/TrainModelNode.vue';
import CheckoutRepoNode from '@/pipelines/steps/components/nodes/CheckoutRepoNode.vue';
import { camel2title, customDelay } from '@/util';

const stepConfigFields = [
  {
    $formkit: 'text',
    label: 'URL',
    name: 'repoURL',
    validation: 'required|url',
    if: '$isActiveNodeType("CheckoutRepo")',
  },
  {
    $formkit: 'text',
    name: 'trainingSetDirectory',
    label: 'Directory',
    validation: 'required|isDirectoryPath',
    if: '$isActiveNodeType("TrainModel")',
  },
  {
    $formkit: 'number',
    name: "fraction",
    label: "Fraction of the training data to use",
    min: 0,
    max: 100,
    validation: 'required',
    if: '$isActiveNodeType("TrainModel")',
  },
  {
    $formkit: 'text',
    name: "modelDirectory",
    label: "Model Directory",
    validation: 'required|isDirectoryPath',
    if: '$isActiveNodeType("TrainModel")',
  },
  {
    $formkit: 'number',
    name: "epochs",
    label: "Number of Epochs",
    validation: 'required',
    min: 1,
    if: '$isActiveNodeType("TrainModel")',
  },
  {
    $formkit: 'number',
    name: "epochs",
    label: "Number of Epochs",
    validation: 'required',
    min: 1,
    if: '$isActiveNodeType("TrainModel")',
  },
]

const nodeTypesOptions = [
  { id: 0, value: "CheckoutRepo", label: "Checkout Repository" },
  { id: 1, value: "TrainModel", label: "Train Model" },
  { id: 2, value: "Scikit", label: "Scikit" },
];

export const nodeTypes = {
  CheckoutRepo: markRaw(CheckoutRepoNode),
  TrainModel: markRaw(TrainModelNode),
  Scikit: markRaw(TrainModelNode),
};

const learningTypes = [
  { id: 0, value: "unsupervised", label: "Unsupervised" },
  { id: 1, value: "supervised", label: "Supervised" },
];

const scikitUnsupervisedModelOptions = [
  { id: 0, value: "ordinaryLeastSquares", label: "Ordinary Least Squares" },
  { id: 1, value: "ridgeRegressionAndClassification", label: "Ridge regression and classification" },
  { id: 2, value: "lasso", label: "Lasso" },
  { id: 3, value: "multiTaskLasso", label: "Multi-task Lasso" },
  { id: 4, value: "elasticNet", label: "Elastic-Net" },
  { id: 5, value: "multiTaskElasticNet", label: "Multi-task Elastic-Net" },
  { id: 6, value: "leastAngleRegression", label: "Least Angle Regression" },
  { id: 7, value: "larsLasso", label: "LARS Lasso" },
  { id: 8, value: "omp", label: "Orthogonal Matching Pursuit" },
  { id: 9, value: "bayesianRegression", label: "Bayesian Regression" },
  { id: 10, value: "logisticRegression", label: "Logistic regression" },
  // { id: 11, value: "generalizedLinearModels", label: "Generalized Linear Models" },
  { id: 12, value: "sgd", label: "Stochastic Gradient Descent" },
  { id: 13, value: "perceptron", label: "Perceptron" },
  { id: 14, value: "passiveAgressiveAlgorithms", label: "Passive Aggressive Algorithms" },
  { id: 15, value: "robustnessRegression", label: "Robustness regression" },
  { id: 16, value: "quantileRegression", label: "Quantile Regression" },
  { id: 17, value: "polynomialRegression", label: "Polynomial regression" },
];

const scikitSupervisedModelOptions = [];

import { reactive, toRef, ref, watch } from 'vue';
import { getNode, createMessage } from '@formkit/core';

export default function useSteps(data, onSubmit) {
  const activeStep = ref('');

  const activeNodeType = ref(nodeTypesOptions[0].value);
  const activeLearningType = ref(learningTypes[0].value);

  const steps = reactive({});
  const visitedSteps = ref([]); // track visited steps

  // NEW: watch our activeStep and store visited steps
  // to know when to show errors
  watch(activeStep, (newStep, oldStep) => {
    if (oldStep && !visitedSteps.value.includes(oldStep)) {
      visitedSteps.value.push(oldStep)
    }
    // NEW: trigger showing validation on fields
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
    if (activeNodeType.value !== "") {
      const stepNames = Object.keys(steps)
      const currentIndex = stepNames.indexOf(activeStep.value)
      activeStep.value = stepNames[currentIndex + delta]
    }
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
    activeNodeType,
    activeLearningType,
    plugins: [
      stepPlugin
    ],
    isActiveNodeType: (nodeType) => {
      return nodeType === activeNodeType.value;
    },
    isSupervisedLearning: () => {
      return activeLearningType.value == "supervised";
    },
    isUnsupervisedLearning: () => {
      return activeLearningType.value == "unsupervised";
    },
    setStep: target => () => {
      setStep(target)
    },
    setActiveNodeType: changeEvent => {
      activeNodeType.value = changeEvent.target.value;
    },
    setLearningType: changeEvent => {
      activeLearningType.value = changeEvent.target.value;
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
                  if: '$activeStep !== "nameAndType"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'nameAndType',
                  name: 'nameAndType',
                  children: [
                    {
                      $formkit: 'text',
                      name: 'nodeName',
                      label: 'Step Name',
                      placeholder: 'Step Name',
                      validation: 'required'
                    },
                    {
                      $formkit: 'select',
                      name: 'nodeType',
                      label: 'Node Type',
                      placeholder: "",
                      options: nodeTypesOptions,
                      validation: 'required',
                      onChange: "$setActiveNodeType",
                    },
                    {
                      $formkit: 'select',
                      name: 'learningType',
                      label: 'Learning Type',
                      placeholder: "",
                      options: learningTypes,
                      validation: 'required',
                      if: '$isActiveNodeType("Scikit")',
                      onChange: "$setLearningType",
                    },
                    {
                      $formkit: 'select',
                      name: 'scikitModel',
                      label: 'Scikit Model',
                      placeholder: "",
                      options: scikitUnsupervisedModelOptions,
                      validation: 'required',
                      if: '$isActiveNodeType("Scikit") && $isUnsupervisedLearning()',
                    },
                    {
                      $formkit: 'select',
                      name: 'scikitModel',
                      label: 'Scikit Model',
                      placeholder: "",
                      options: scikitSupervisedModelOptions,
                      validation: 'required',
                      if: '$isActiveNodeType("Scikit") && $isSupervisedLearning()',
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
                  disabled: '$activeStep === "nameAndType"',
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

  // NEW: include visitedSteps in our return
  return { formSchema, formkitData, activeStep, nodeTypesOptions, activeNodeType, stepConfigFields, visitedSteps, steps, stepPlugin, setStep }
}