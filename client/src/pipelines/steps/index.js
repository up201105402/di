import { markRaw } from 'vue';
import { camel2title } from '@/util';
import General from '@/pipelines/steps/components/nodes/General.vue';
import ScikitDatasets from '@/pipelines/steps/components/nodes/ScikitDatasets.vue';
import ScikitUnsupervisedModels from '@/pipelines/steps/components/nodes/ScikitUnsupervisedModels.vue';
import { checkoutRepoForm } from '@/pipelines/steps/checkoutRepo';
import { scriptForm } from '@/pipelines/steps/script';
import { scikitDatasetForm } from '@/pipelines/steps/scikit/datasets';
import { scikitUnsupervisedModels } from '@/pipelines/steps/scikit/models';
import ScriptEditor from '@/components/ScriptEditor.vue';
import FormFilePicker from '@/components/FormFilePicker.vue';
import PythonScriptNode from '@/pipelines/steps/components/nodes/PythonScript.vue';
import ShellScriptNode from '@/pipelines/steps/components/nodes/ShellScript.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseCancelAndSubmitButtons from '@/components/BaseCancelAndSubmitButtons.vue';

export const nodeTypes = {
  general: markRaw(General),
  scikitDatasets: markRaw(ScikitDatasets),
  scikitUnsupervisedModels: markRaw(ScikitUnsupervisedModels),
  checkoutRepo: markRaw(General),
  shellScript: markRaw(ShellScriptNode),
  pythonScript: markRaw(PythonScriptNode),
};

Object.keys(scikitUnsupervisedModels).forEach(key => {
  nodeTypes[key] = markRaw(ScikitUnsupervisedModels);
});

const scikitUnsupervisedSteps = Object.entries(scikitUnsupervisedModels)
  .map(item => {
    return {
      group: 'scikitUnsupervisedModels',
      type: item[0],
      label: camel2title(item[0]),
      form: item[1]
    }
  });

export const menubarSteps = [
  {
    type: 'general',
    label: 'General',
    items: [
      {
        type: 'checkoutRepo',
        label: 'Checkout Repository',
        form: checkoutRepoForm
      },
      {
        separator: true,
      },
      {
        type: 'shellScript',
        label: 'Shell Script',
        form: scriptForm
      },
      {
        type: 'pythonScript',
        label: 'Python Script',
        form: scriptForm
      }
    ]
  },
  {
    type: 'scikit',
    label: 'Scikit',
    items: [
      {
        type: 'scikitDatasets',
        label: 'Scikit Datasets',
        items: [
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
        items: scikitUnsupervisedSteps
      },
    ]
  },
];

export const library = markRaw({
  ScriptEditor: ScriptEditor,
  FormFilePicker: FormFilePicker,
  BaseButton: BaseButton,
  BaseCancelAndSubmitButtons: BaseCancelAndSubmitButtons,
})