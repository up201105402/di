export const checkoutRepo = {
    id: 0,
    name: 'checkoutRepo',
    label: 'Checkout repository',
    conditions: [
        [
            'stepType',
            'in',
            [
                '0',
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
            'Step  Type',
            'in',
            [
                '1',
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
            'Step  Type',
            'in',
            [
                '2',
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

export default {
    checkoutRepo,
    loadTrainingDataset,
    trainModel
}