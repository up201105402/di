export const ridgeRegressionStepConfig = [
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'checkbox',
        name: "fitIntercept",
        label: "Fit Intercept",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'checkbox',
        name: "copyX",
        label: "Copy X",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'number',
        name: "maxIter",
        label: "max_iter",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'number',
        name: "toI",
        label: "toI",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "Solver",
        options: ['auto', 'svd', 'cholesky', 'lsqr', 'sparse_cg', 'sag', 'saga', 'lbfgs'],
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
    {
        $formkit: 'randomState',
        name: "randomState",
        label: "Random State",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegression")',
    },
]

export const ridgeRegressionCVStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fitIntercept",
        label: "Fit Intercept",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegressionCV")',
    },
    {
        $formkit: 'checkbox',
        name: "copyX",
        label: "Copy X",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegressionCV")',
    },
    {
        $formkit: 'number',
        name: "nJobs",
        label: "Computation number of Jobs",
        placeholder: '1',
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegressionCV")',
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("ridgeRegressionCV")',
    },
]