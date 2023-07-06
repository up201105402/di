export const leastSquaresStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fitIntercept_1",
        label: "Fit Intercept 1",
    },
    {
        $formkit: 'checkbox',
        name: "copyX_1",
        label: "Copy X 1",
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

export const leastSquaresConfigSection = {
    $formkit: 'group',
    id: 'stepConfig_1',
    name: 'stepConfig',
    children: leastSquaresStepConfig,
    if: '$isActiveNodeType("ScikitModel") && $isScikitModel("leastSquares")',
}