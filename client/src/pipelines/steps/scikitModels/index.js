import { leastSquaresStepConfig } from './leastSquares';
import { ridgeRegressionStepConfig } from './ridgeRegression';

export const learningTypes = [
    { id: 0, value: "unsupervised", label: "Unsupervised" },
    { id: 1, value: "supervised", label: "Supervised" },
];

export const scikitUnsupervisedModelOptions = [
    { id: 0, value: "leastSquares", label: "Ordinary Least Squares" },
    { id: 1, value: "ridgeRegression", label: "Ridge Regression" },
    { id: 2, value: "ridgeRegressionCV", label: "Ridge Regression CV" },
    { id: 3, value: "lasso", label: "Lasso" },
    { id: 4, value: "multiTaskLasso", label: "Multi-task Lasso" },
    { id: 5, value: "elasticNet", label: "Elastic-Net" },
    { id: 6, value: "multiTaskElasticNet", label: "Multi-task Elastic-Net" },
    { id: 7, value: "leastAngleRegression", label: "Least Angle Regression" },
    { id: 8, value: "larsLasso", label: "LARS Lasso" },
    { id: 9, value: "omp", label: "Orthogonal Matching Pursuit" },
    { id: 10, value: "bayesianRegression", label: "Bayesian Regression" },
    { id: 11, value: "logisticRegression", label: "Logistic regression" },
    // { id: 11, value: "generalizedLinearModels", label: "Generalized Linear Models" },
    { id: 12, value: "sgd", label: "Stochastic Gradient Descent" },
    { id: 13, value: "perceptron", label: "Perceptron" },
    { id: 14, value: "passiveAgressiveAlgorithms", label: "Passive Aggressive Algorithms" },
    { id: 15, value: "robustnessRegression", label: "Robustness regression" },
    { id: 16, value: "quantileRegression", label: "Quantile Regression" },
    { id: 17, value: "polynomialRegression", label: "Polynomial regression" },
];

export const scikitSupervisedModelOptions = [];

export const scikitLearningTypeSelect = [
    {
        $formkit: 'select',
        name: 'learningType',
        label: 'Learning Type',
        placeholder: "",
        options: learningTypes,
        validation: 'required',
        if: '$isActiveNodeType("ScikitModel")',
        onChange: "$setLearningType",
    },
]

export const scikitUnsupervisedModelSelect = [
    {
        $formkit: 'select',
        name: 'scikitModel',
        label: 'Scikit Model',
        placeholder: "",
        options: scikitUnsupervisedModelOptions,
        validation: 'required',
        if: '$isActiveNodeType("ScikitModel") && $isUnsupervisedLearning()',
        onChange: "$setScikitModel"
    },
]

export const scikitSupervisedModelSelect = [
    {
        $formkit: 'select',
        name: 'scikitModel',
        label: 'Scikit Model',
        placeholder: "",
        options: scikitSupervisedModelOptions,
        validation: 'required',
        if: '$isActiveNodeType("ScikitModel") && $isSupervisedLearning()',
    },
]

export const scikitModelsStepConfig = [
    ...leastSquaresStepConfig,
    ...ridgeRegressionStepConfig
]