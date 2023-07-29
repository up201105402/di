export const larsStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'checkbox',
        name: "precompute",
        label: "precompute",
    },
    {
        $formkit: 'number',
        name: "n_nonzero_coefs",
        label: "n_nonzero_coefs",
    },
    {
        $formkit: 'number',
        name: "eps",
        label: "eps",
        validation: "required",
    },
    {
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    },
    {
        $formkit: 'checkbox',
        name: "fit_path",
        label: "fit_path",
    },
    {
        $formkit: 'number',
        name: "jitter",
        label: "jitter",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
]

export const larsCVStepConfig = [
    {
        $formkit: 'checkbox',
        name: "fit_intercept",
        label: "fit_intercept",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'checkbox',
        name: "precompute",
        label: "precompute",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'number',
        name: "max_n_alphas",
        label: "max_n_alphas",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "eps",
        label: "eps",
        validation: "required",
    },
    {
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    }
]