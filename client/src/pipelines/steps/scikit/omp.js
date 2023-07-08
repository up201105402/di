export const ompStepConfig = [
    {
        $formkit: 'number',
        name: "n_nonzero_coefs",
        label: "n_nonzero_coefs",
    },
    {
        $formkit: 'number',
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'checkbox',
        name: "precompute",
        label: "precompute",
    },
]

export const ompCVStepConfig = [
    {
        $formkit: 'checkbox',
        name: "copy",
        label: "copy",
    },
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
]