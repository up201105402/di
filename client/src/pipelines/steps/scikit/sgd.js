export const sgdClassifierStepConfig = [
    {
        $formkit: 'select',
        name: "loss",
        label: "loss",
        options: [ 'hinge', 'log_loss', 'modified_huber', 'squared_hinge', 'perceptron', 'squared_error', 'huber', 'epsilon_insensitive', 'squared_epsilon_insensitive' ]
    },
    {
        $formkit: 'select',
        name: "penalty",
        label: "penalty",
        options: [ 'l2', 'l1', 'elasticnet' ]
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'number',
        name: "l1_ratio",
        label: "l1_ratio",
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
        name: "shuffle",
        label: "shuffle",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "epsilon",
        label: "epsilon",
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
        name: "learning_rate",
        label: "learning_rate",
        options: [ 'optimal', 'constant', 'invscaling', 'adaptive' ]
    },
    {
        $formkit: 'number',
        name: "eta0",
        label: "eta0",
    },
    {
        $formkit: 'number',
        name: "power_t",
        label: "power_t",
    },
    {
        $formkit: 'checkbox',
        name: "early_stopping",
        label: "early_stopping",
    },
    {
        $formkit: 'number',
        name: "validation_fraction",
        label: "validation_fraction",
    },
    {
        $formkit: 'number',
        name: "n_iter_no_change",
        label: "n_iter_no_change",
    },
    {
        $formkit: 'text',
        name: "class_weight",
        label: "class_weight",
        validation: 'dict'
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'checkbox', // this should be a bool or int in reality
        name: "average",
        label: "average",
    },
]

export const sgdRegressorStepConfig = [
    {
        $formkit: 'select',
        name: "loss",
        label: "loss",
        options: [ 'squared_error', 'huber', 'epsilon_insensitive', 'squared_epsilon_insensitive' ]
    },
    {
        $formkit: 'select',
        name: "penalty",
        label: "penalty",
        options: [ 'l2', 'l1', 'elasticnet' ]
    },
    {
        $formkit: 'number',
        name: "alpha",
        label: "alpha",
    },
    {
        $formkit: 'number',
        name: "l1_ratio",
        label: "l1_ratio",
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
        name: "shuffle",
        label: "shuffle",
    },
    {
        $formkit: 'number',
        name: "verbose",
        label: "verbose",
    },
    {
        $formkit: 'number',
        name: "epsilon",
        label: "epsilon",
    },
    {
        $formkit: 'number',
        name: "random_state",
        label: "random_state",
    },
    {
        $formkit: 'select',
        name: "learning_rate",
        label: "learning_rate",
        options: [ 'invscaling', 'optimal', 'constant', 'adaptive' ]
    },
    {
        $formkit: 'number',
        name: "eta0",
        label: "eta0",
    },
    {
        $formkit: 'number',
        name: "power_t",
        label: "power_t",
    },
    {
        $formkit: 'checkbox',
        name: "early_stopping",
        label: "early_stopping",
    },
    {
        $formkit: 'number',
        name: "validation_fraction",
        label: "validation_fraction",
    },
    {
        $formkit: 'number',
        name: "n_iter_no_change",
        label: "n_iter_no_change",
    },
    {
        $formkit: 'checkbox',
        name: "warm_start",
        label: "warm_start",
    },
    {
        $formkit: 'checkbox', // this should be a bool or int in reality
        name: "average",
        label: "average",
    },
]