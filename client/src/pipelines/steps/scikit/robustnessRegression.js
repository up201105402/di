export const huberRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "epsilon",
        label: "epsilon",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
]

export const ransacRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "epsilon",
        label: "epsilon",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
]

export const theilSenRegressorStepConfig = [
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
        name: "max_subpopulation",
        label: "max_subpopulation",
    },
    {
        $formkit: 'number',
        name: "n_subsamples",
        label: "n_subsamples",
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
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
]