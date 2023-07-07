export const ridgeRegressionStepConfig = [
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'checkbox',
        name: "fitIntercept",
        label: "Fit Intercept",
    },
    {
        $formkit: 'checkbox',
        name: "copyX",
        label: "Copy X",
    },
    {
        $formkit: 'number',
        name: "maxIter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "toI",
        label: "toI",
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
        name: "randomState",
        label: "Random State",
    },
]

export const ridgeRegressionCVStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fitIntercept",
        label: "Fit Intercept",
    },
    {
        $formkit: 'checkbox',
        name: "copyX",
        label: "Copy X",
    },
    {
        $formkit: 'number',
        name: "nJobs",
        label: "Computation number of Jobs",
        placeholder: '1',
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
    },
]
