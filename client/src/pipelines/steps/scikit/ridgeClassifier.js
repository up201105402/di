export const ridgeClassifierStepConfig = [
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
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
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

export const ridgeClassifierCVStepConfig = [
    {
        $formkit: 'text',
        name: "alphas",
        label: "alphas",
        validation: "required|floats"
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'text',
        name: "scoring",
        label: "scoring",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'checkbox',
        name: "store_cv_values",
        label: "store_cv_values",
    },
]
