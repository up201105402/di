export const logisticRegressionStepConfig = [
    {
        $formkit: 'select',
        name: "penalty",
        label: "penalty",
        options: [ 'l2', 'l1', 'elasticnet' ]  // None was deprecated
    },
    {
        $formkit: 'checkbox',
        name: "dual",
        label: "dual",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'number',
        name: "C",
        label: "C",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "intercept_scaling",
        label: "intercept_scaling",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "Solver",
        options: ['lbfgs', 'liblinear', 'newton-cg', 'newton-cholesky', 'sag', 'saga'],
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'select',
        name: "multi_class",
        label: "multi_class",
        options: ['auto', 'ovr', 'multinomial'],
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "l1_ratio",
        label: "l1_ratio",
    },
]

export const logisticRegressionCVStepConfig = [
    {
        $formkit: 'number', //only ints
        name: "Cs",
        label: "Cs",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'checkbox',
        name: "dual",
        label: "dual",
    },
    {
        $formkit: 'select',
        name: "penalty",
        label: "penalty",
        options: [ 'l2', 'l1', 'elasticnet' ]
    },
    {
        $formkit: 'text',
        name: "scoring",
        label: "scoring",
    },
    {
        $formkit: 'select',
        name: "solver",
        label: "Solver",
        options: ['lbfgs', 'liblinear', 'newton-cg', 'newton-cholesky', 'sag', 'saga'],
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'checkbox',
        name: "refit",
        label: "refit",
    },
    {
        $formkit: 'number',
        name: "intercept_scaling",
        label: "intercept_scaling",
    },
    {
        $formkit: 'select',
        name: "multi_class",
        label: "multi_class",
        options: ['auto', 'ovr', 'multinomial'],
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'number',
        name: "l1_ratios",
        label: "l1_ratios",
        validation: 'floats'
    },
]