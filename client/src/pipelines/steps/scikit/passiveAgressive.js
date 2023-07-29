
export const passiveAgressiveClassifierStepConfig = [
    {
        $formkit: 'number',
        name: "C",
        label: "C",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'checkbox',
        name: "early_stopping",
        label: "early_stopping",
    },
    {
        $formkit: 'number',
        name: "validation_fraction",
        label: "validation_fraction",
    },
    {
        $formkit: 'number',
        name: "n_iter_no_change",
        label: "n_iter_no_change",
    },
    {
        $formkit: 'number',
        name: "shuffle",
        label: "shuffle",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'select',
        name: "loss",
        label: "loss",
        options: [ 'hinge', 'squared_hinge' ]
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'checkbox',
        name: "average",
        label: "average",
    },
]

export const passiveAgressiveRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "C",
        label: "C",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'checkbox',
        name: "early_stopping",
        label: "early_stopping",
    },
    {
        $formkit: 'number',
        name: "validation_fraction",
        label: "validation_fraction",
    },
    {
        $formkit: 'number',
        name: "n_iter_no_change",
        label: "n_iter_no_change",
    },
    {
        $formkit: 'number',
        name: "shuffle",
        label: "shuffle",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'select',
        name: "loss",
        label: "loss",
        options: [ 'hinge', 'squared_hinge' ]
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'checkbox',
        name: "average",
        label: "average",
    },
]