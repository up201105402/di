import { markRaw } from 'vue';
import CheckoutRepoNode from '@/pipelines/steps/components/nodes/CheckoutRepoNode.vue';
import TrainModelNode from '@/pipelines/steps/components/nodes/TrainModelNode.vue';
import { checkoutRepoForm } from '@/pipelines/steps/checkoutRepo';
import { scikitDatasetForm } from '@/pipelines/steps/scikit/datasets';
import { scikitLeastSquaresForm, ridgeRegressionForm, ridgeRegressionCVForm } from '@/pipelines/steps/scikit/models';

export const nodeTypes = {
  CheckoutRepo: markRaw(CheckoutRepoNode),
  TrainModel: markRaw(TrainModelNode),
  ScikitTrainingDataset: markRaw(TrainModelNode),
  ScikitTestingDataset: markRaw(TrainModelNode),
  ScikitModel: markRaw(TrainModelNode),
};

export const steps = [
  {
    name: 'general',
    label: 'General',
    steps: [
      {
        name: 'checkoutRepo',
        label: 'Checkout Repository',
        form: checkoutRepoForm
      },
    ]
  },
  {
    name: 'scikit',
    label: 'Scikit Datasets',
    steps: [
      {
        name: 'scikitTrainingDataset',
        label: 'Load Training Dataset',
        form: scikitDatasetForm
      },
      {
        name: 'scikitTestingDataset',
        label: 'Load Testing Dataset',
        form: scikitDatasetForm
      },
    ]
  },
  {
    name: 'scikit',
    label: 'Scikit Models',
    steps: [
      {
        name: 'leastSquares',
        label: 'Least Squares',
        form: scikitLeastSquaresForm
      },
      {
        name: 'ridgeRegression',
        label: 'Ridge Regression',
        form: ridgeRegressionForm
      },
      {
        name: 'ridgeRegressionCV',
        label: 'Ridge Regression CV',
        form: ridgeRegressionCVForm
      },
      // {
      //   name: 'ridgeClassifier',
      //   label: 'Ridge Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'ridgeClassifierCV',
      //   label: 'Ridge Classifier CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lasso',
      //   label: 'Lasso',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lassoCV',
      //   label: 'Lasso CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lassoLars',
      //   label: 'Lasso Lars',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lassoLarsCV',
      //   label: 'Lasso Lars CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lassoLarsIc',
      //   label: 'Lasso Lars IC',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'multiTaskLasso',
      //   label: 'Mult-Task Lasso',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'multiTaskLassoCV',
      //   label: 'Mult-Task Lasso CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'elasticNet',
      //   label: 'Elastic Net',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'elasticNetCV',
      //   label: 'Elastic Net CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'multiTaskElasticNet',
      //   label: 'Multi-Task Elastic Net',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'multiTaskElasticNetCV',
      //   label: 'Multi-Task Elastic Net CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'lars',
      //   label: 'LARS',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'larsCV',
      //   label: 'LARS CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'omp',
      //   label: 'Orthogonal Matching Puirsuit (OMP)',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'ompCV',
      //   label: 'Orthogonal Matching Puirsuit (OMP) CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'bayesianRidgeRegression',
      //   label: 'Bayesian Ridge Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'bayesianARDRegression',
      //   label: 'Bayesian ARD Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'logisticRegression',
      //   label: 'Logistic Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'logisticRegressionCV',
      //   label: 'Logistic Regression CV',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'tweedieRegressor',
      //   label: 'Tweedie Regressor',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'poissonRegressor',
      //   label: 'Poisson Regressor',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'gammaRegressor',
      //   label: 'Gamma Regressor',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'sgdRegressor',
      //   label: 'SGD Regressor',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'perceptron',
      //   label: 'Perceptron',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'passiveAgressiveClassifier',
      //   label: 'Passive Aggressive Classifier',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'passiveAgressiveRegressor',
      //   label: 'Passive Aggressive Regressor',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'huberRegression',
      //   label: 'Huber Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'ransacRegression',
      //   label: 'RANSAC Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'theilSenRegression',
      //   label: 'Theil-Sen Regression',
      //   form: ridgeRegressionCVForm
      // },
      // {
      //   name: 'quantileRegression',
      //   label: 'Quantile Regression',
      //   form: ridgeRegressionCVForm
      // },
    ]
  },
];