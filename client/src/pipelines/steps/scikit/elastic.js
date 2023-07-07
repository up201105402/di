export const multiTaskElasticNetStepConfig = [
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
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
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
        $formkit: 'number',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]

export const multiTaskElasticNetCVStepConfig = [
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
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
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
        $formkit: 'number',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]