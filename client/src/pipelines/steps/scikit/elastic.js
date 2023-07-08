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
        $formkit: 'number',
        name: "eps",
        label: "eps",
    },
    {
        $formkit: 'number',
        name: "n_alphas",
        label: "n_alphas",
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
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
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
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]