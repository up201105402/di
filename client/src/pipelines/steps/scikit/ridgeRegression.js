export const ridgeRegressionStepConfig = [
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
        validation: "required",
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
        $formkit: 'select',
        name: "solver",
        label: "Solver",
        options: ['auto', 'svd', 'cholesky', 'lsqr', 'sparse_cg', 'sag', 'saga', 'lbfgs'],
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
]

export const ridgeRegressionCVStepConfig = [
    {
        $formkit: 'text',
        name: "alphas",
        label: "alphas",
        validation: "required|floats",
        help: "Example: 0.1 0.2 0.3"
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
        $formkit: 'select',
        name: "solver",
        label: "Solver",
        options: ['auto', 'svd', 'cholesky', 'lsqr', 'sparse_cg', 'sag', 'saga', 'lbfgs'],
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
]
