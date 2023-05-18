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
                'regex:^(.+)\\/([^\\/]+)$',
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
                'regex:^(.+)\\/([^\\/]+)$',
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

export default {
    checkoutRepo,
    loadTrainingDataset,
    trainModel
}