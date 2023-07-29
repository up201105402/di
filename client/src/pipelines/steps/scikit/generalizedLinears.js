export const tweedieRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "power",
        label: "power",
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'select',
        name: "link",
        label: "link",
        options: [ 'auto', 'identity', 'log' ]
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "solver",
        options: [ 'lbfgs', 'newton-cholesky' ]
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
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
]

export const poissonRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "solver",
        options: [ 'lbfgs', 'newton-cholesky' ]
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
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
]

export const gammaRegressorStepConfig = [
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "solver",
        options: [ 'lbfgs', 'newton-cholesky' ]
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
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
]