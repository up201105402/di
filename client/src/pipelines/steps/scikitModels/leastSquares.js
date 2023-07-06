export const leastSquaresStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fitIntercept",
        label: "Fit Intercept",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("leastSquares")',
    },
    {
        $formkit: 'checkbox',
        name: "copyX",
        label: "Copy X",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("leastSquares")',
    },
    {
        $formkit: 'number',
        name: "nJobs",
        label: "Computation number of Jobs",
        placeholder: '1',
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("leastSquares")',
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
        if: '$isActiveNodeType("ScikitModel") && $isScikitModel("leastSquares")',
    },
]