export const lassoStepConfig = [
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
        name: "precompute",
        label: "precompute",
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
        $formkit: 'number',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "Positive",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]

export const lassoCVStepConfig = [
    {
        $formkit: 'number',
        name: "eps",
        label: "eps",
        validation: "required",
    },
    {
        $formkit: 'number',
        name: "n_alphas",
        label: "n_alphas",
    },
    {
        $formkit: 'text',
        name: "alphas",
        label: "alphas",
        validation: "floats",
        help: "Example: 0.1 0.2 0.3"
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
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "positive",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]

export const lassoLarsStepConfig = [
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
        name: "max_iter",
        label: "max_iter",
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
        $formkit: 'checkbox',
        name: "positive",
        label: "positive",
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

export const lassoLarsCVStepConfig = [
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
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "positive",
    },
]

export const lassoLarsICStepConfig = [
    {
        $formkit: 'select',
        name: "criterion",
        label: "selection",
        options: ['aic', 'bic']
    },
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
        name: "max_iter",
        label: "max_iter",
    },
    {
        $formkit: 'number',
        name: "eps",
        label: "eps",
    },
    {
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    },
    {
        $formkit: 'checkbox',
        name: "positive",
        label: "positive",
    },
    {
        $formkit: 'number',
        name: "noise_variance",
        label: "noise_variance",
    },
]

export const multiTaskLassoStepConfig = [
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
        $formkit: 'number',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]

export const multiTaskLassoCVStepConfig = [
    {
        $formkit: 'number',
        name: "eps",
        label: "eps",
    },
    {
        $formkit: 'number',
        name: "n_alphas",
        label: "n_alphas",
    },
    {
        $formkit: 'text',
        name: "alphas",
        label: "alphas",
        validation: "floats",
        help: "Example: 0.1 0.2 0.3"
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
        name: "tol",
        label: "tol",
    },
    {
        $formkit: 'checkbox',
        name: "copy_X",
        label: "copy_X",
    },
    {
        $formkit: 'number',
        name: "cv",
        label: "cv",
    },
    {
        $formkit: 'checkbox',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "n_jobs",
        label: "n_jobs",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "selection",
        label: "selection",
        options: ['cyclic', 'random']
    },
]