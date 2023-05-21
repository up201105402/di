import { markRaw } from 'vue';
import CustomNode from '@/pipelines/steps/components/nodes/CustomNode.vue';
import CheckoutRepoNode from '@/pipelines/steps/components/nodes/CheckoutRepoNode.vue';

export const checkoutRepo = {
    id: 0,
    name: 'checkoutRepo',
    label: 'Checkout repository',
    conditions: [
        [
            'stepType',
            'in',
            [
                'checkoutRepo',
            ],
        ],
    ],
    fields: {
        repoURL: {
            name: "repoURL",
            label: "URL",
            rules: [
                'required',
                'regex:((git|ssh|http(s)?)|(git@[\\w\\.]+))(:(//)?)([\\w\\.@\\:/\\-~]+)(\\.git)(/)?/',
            ]
        }
    }
};

export const loadTrainingDataset = {
    id: 1,
    name: 'loadTrainingDataset',
    label: 'Load Training Dataset',
    conditions: [
        [
            'stepType',
            'in',
            [
                'loadTrainingDataset',
            ],
        ],
    ],
    fields: {
        trainingDatasetDirectory: {
            name: "trainingSetDirectory",
            label: "Directory",
            rules: [
                'required',
                'regex:^(.+)\\/([^\\/]+)$/',
            ]
        },
        fraction: {
            name: "fraction",
            label: "Fraction of the training data to use",
        }
    }
};

export const trainModel = {
    id: 2,
    name: 'trainModel',
    label: 'Train Model',
    conditions: [
        [
            'stepType',
            'in',
            [
                'trainModel',
            ],
        ],
    ],
    fields: {
        modelDirectory: {
            name: "modelDirectory",
            label: "Model Directory",
            rules: [
                'required',
                'regex:^(.+)\\/([^\\/]+)$/',
            ]
        },
        epochs: {
            name: "epochs",
            label: "Number of Epochs",
            rules: [
                'required',
                'min:1',
                'numeric',
            ]
        }
    }
};

export const nodeTypes = {
    checkoutRepo: markRaw(CheckoutRepoNode),
    custom: markRaw(CustomNode)
};

// export default {
//     checkoutRepo,
//     loadTrainingDataset,
//     trainModel
// }

import { reactive, toRef, ref, watch } from 'vue'
import { getNode, createMessage } from '@formkit/core'

export default function useSteps () {
  const activeStep = ref('')
  const steps = reactive({})
  const visitedSteps = ref([]) // track visited steps

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
    const stepNames = Object.keys(steps)
    const currentIndex = stepNames.indexOf(activeStep.value)
    activeStep.value = stepNames[currentIndex + delta]
  }

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

  // NEW: include visitedSteps in our return
  return { activeStep, visitedSteps, steps, stepPlugin, setStep }
}