import { markRaw } from 'vue';
import { camel2title } from '@/util';
import General from '@/pipelines/steps/components/nodes/General.vue';
import ScikitDatasets from '@/pipelines/steps/components/nodes/ScikitDatasets.vue';
import ScikitUnsupervisedModels from '@/pipelines/steps/components/nodes/ScikitUnsupervisedModels.vue';
import { checkoutRepoForm } from '@/pipelines/steps/checkoutRepo';
import { scikitDatasetForm } from '@/pipelines/steps/scikit/datasets';
import { scikitUnsupervisedModels } from '@/pipelines/steps/scikit/models';

export const nodeTypes = {
  general: markRaw(General),
  scikitDatasets: markRaw(ScikitDatasets),
  scikitUnsupervisedModels: markRaw(ScikitUnsupervisedModels),
};

const scikitUnsupervisedSteps = Object.entries(scikitUnsupervisedModels)
  .map(item => {
    return {
      group: 'scikitUnsupervisedModels',
      type: item[0],
      label: camel2title(item[0]),
      form: item[1]
    }
  });

export const steps = [
  {
    type: 'general',
    label: 'General',
    steps: [
      {
        group: 'general',
        type: 'checkoutRepo',
        label: 'Checkout Repository',
        form: checkoutRepoForm
      },
    ]
  },
  {
    type: 'scikitDatasets',
    label: 'Scikit Datasets',
    steps: [
      {
        group: 'scikitDatasets',
        type: 'scikitTrainingDataset',
        label: 'Load Training Dataset',
        form: scikitDatasetForm
      },
      {
        group: 'scikitDatasets',
        type: 'scikitTestingDataset',
        label: 'Load Testing Dataset',
        form: scikitDatasetForm
      },
    ]
  },
  {
    type: 'scikitUnsupervisedModels',
    label: 'Scikit Unsupervised Models',
    steps: scikitUnsupervisedSteps
  },
];