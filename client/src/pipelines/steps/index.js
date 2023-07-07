import { markRaw } from 'vue';
import CheckoutRepoNode from '@/pipelines/steps/components/nodes/CheckoutRepoNode.vue';
import TrainModelNode from '@/pipelines/steps/components/nodes/TrainModelNode.vue';
import { checkoutRepoForm } from '@/pipelines/steps/checkoutRepo';
import { scikitDatasetForm } from '@/pipelines/steps/scikit/datasets';
import { scikitLeastSquaresForm, ridgeRegressionForm, ridgeRegressionCVForm } from '@/pipelines/steps/scikit/models';

const nodeTypesOptions = [
  { id: 0, value: "CheckoutRepo", label: "Checkout Repository" },
  { id: 1, value: "ScikitTrainingDataset", label: "Scikit - Load Training Dataset" },
  { id: 2, value: "ScikitTestingDataset", label: "Scikit - Load Testing Dataset" },
  { id: 3, value: "ScikitModel", label: "Scikit - Model" },
];

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
      } 
    ]
  },
];