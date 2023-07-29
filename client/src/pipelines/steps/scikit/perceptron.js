
export const perceptronStepConfig = [
    {
        $formkit: 'select',
        name: "penalty",
        label: "penalty",
        options: [ 'l2', 'l1', 'elasticnet' ]
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'number',
        name: "l1_ratio",
        label: "l1_ratio",
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
        name: "shuffle",
        label: "shuffle",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "eta0",
        label: "eta0",
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
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
]