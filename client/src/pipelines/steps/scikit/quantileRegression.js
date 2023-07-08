export const quantileRegressionStepConfig = [
    {
        $formkit: 'number',
        name: "quantile",
        label: "quantile",
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
        name: "solver",
        label: "solver",
        options: [ 'interior-point', 'highs-ds', 'highs-ipm', 'highs', 'revised simplex' ]
    },
    {
        $formkit: 'text',
        name: "solver_options",
        label: "solver_options",
        validation: 'dict'
    },
]