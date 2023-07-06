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
        $formkit: 'randomState',
        name: "randomState",
        label: "Random State",
    },
]

export const ridgeRegressionConfigSection = {
    $formkit: 'group',
    id: 'stepConfig',
    name: 'stepConfig',
    children: ridgeRegressionStepConfig,
    if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
}

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

export const ridgeRegressionCVConfigSection = {
    $formkit: 'group',
    id: 'stepConfig',
    name: 'stepConfig',
    children: ridgeRegressionCVStepConfig,
    if: '$isScikitModel("ridgeRegressionCV")',
}