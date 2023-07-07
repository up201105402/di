import { reactive, toRef, ref, watch, markRaw } from 'vue';
import CheckoutRepoNode from '@/pipelines/steps/components/nodes/CheckoutRepoNode.vue';
import TrainModelNode from '@/pipelines/steps/components/nodes/TrainModelNode.vue';
import { camel2title, customDelay } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { checkoutRepoConfigSection } from '@/pipelines/steps/checkoutRepo';
import { datasetConfigSection, scikitDatasetSelect, scikitDatasets } from '@/pipelines/steps/datasets';
import { learningTypes, scikitLearningTypeSelect, scikitUnsupervisedModelSelect, scikitUnsupervisedModelOptions, scikitSupervisedModelSelect, scikitModelConfigSections } from '@/pipelines/steps/scikitModels';

const nodeTypesOptions = [
  { id: 0, value: "CheckoutRepo", label: "Checkout Repository" },
  { id: 1, value: "ScikitTrainingDataset", label: "Scikit - Load Training Dataset" },
  { id: 2, value: "ScikitTestingDataset", label: "Scikit - Load Testing Dataset" },
  { id: 3, value: "ScikitModel", label: "Scikit - Model" },
];

export const nodeTypes = {
  CheckoutRepo: markRaw(CheckoutRepoNode),
  TrainModel: markRaw(TrainModelNode),
  ScikitTrainingDataset: markRaw(TrainModelNode),
  ScikitTestingDataset: markRaw(TrainModelNode),
  ScikitModel: markRaw(TrainModelNode),
};

export default {
  default: function (data, onSubmit) {

    const formkitData = reactive({});

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
        }
      },
    ];

    return { formSchema, formkitData, }
  },
  checkoutRepo: function (data, onSubmit) {
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
      const stepNames = Object.keys(steps)
      const currentIndex = stepNames.indexOf(activeStep.value)
      activeStep.value = stepNames[currentIndex + delta]
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
                  checkoutRepoConfigSection,
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

    return { formSchema, formkitData, activeStep, nodeTypesOptions, visitedSteps, steps, stepPlugin, setStep }
  },
  scikitTrainingDataset: function (data, onSubmit) {
    const activeStep = ref('');

    const activeNodeType = ref(nodeTypesOptions[0].value);
    const activeLearningType = ref(learningTypes[0].value);
    const activeScikitDataset = ref(scikitDatasets[0].value);
    const activeScikitModel = ref(scikitUnsupervisedModelOptions[0].value);

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
      activeScikitDataset,
      activeScikitModel,
      plugins: [
        stepPlugin
      ],
      isActiveNodeType: (nodeType) => {
        return nodeType === activeNodeType.value;
      },
      isScikitDataset: (dataset) => {
        return activeScikitDataset.value === dataset;
      },
      isScikitModel: (model) => {
        return activeScikitModel.value === model;
      },
      isSupervisedLearning: () => {
        return activeLearningType.value === "supervised";
      },
      isUnsupervisedLearning: () => {
        return activeLearningType.value === "unsupervised";
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
      setSciKitDataset: changeEvent => {
        activeScikitDataset.value = changeEvent.target.value;
      },
      setScikitModel: changeEvent => {
        activeScikitModel.value = changeEvent.target.value;
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
                  checkoutRepoConfigSection,
                  datasetConfigSection,
                  ...scikitModelConfigSections
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

    return { formSchema, formkitData, activeStep, nodeTypesOptions, activeNodeType, visitedSteps, steps, stepPlugin, setStep }
  },
}