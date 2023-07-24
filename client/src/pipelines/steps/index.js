import { markRaw } from 'vue';
import General from '@/pipelines/steps/components/nodes/General.vue';
import ScikitDatasets from '@/pipelines/steps/components/nodes/ScikitDatasets.vue';
import ScikitUnsupervisedModels from '@/pipelines/steps/components/nodes/ScikitUnsupervisedModels.vue';
import { checkoutRepoForm } from '@/pipelines/steps/checkoutRepo';
import { scriptForm } from '@/pipelines/steps/script';
import { hitlForm } from '@/pipelines/steps/hitl';
import { customHitlForm } from '@/pipelines/steps/customHitl';
import { scikitDatasetForm } from '@/pipelines/steps/scikit/datasets';
import { scikitUnsupervisedModels } from '@/pipelines/steps/scikit/models';
import ScriptEditor from '@/components/ScriptEditor.vue';
import FormFilePicker from '@/components/FormFilePicker.vue';
import PythonScriptNode from '@/pipelines/steps/components/nodes/PythonScript.vue';
import ShellScriptNode from '@/pipelines/steps/components/nodes/ShellScript.vue';
import HITLNode from '@/pipelines/steps/components/nodes/HITL.vue';
import BaseButton from '@/components/BaseButton.vue';
import BaseCancelAndSubmitButtons from '@/components/BaseCancelAndSubmitButtons.vue';
import { i18n } from '@/i18n';

const { t } = i18n.global;

export const nodeTypes = {
  general: markRaw(General),
  scikitDatasets: markRaw(ScikitDatasets),
  scikitUnsupervisedModels: markRaw(ScikitUnsupervisedModels),
  checkoutRepo: markRaw(General),
  shellScript: markRaw(ShellScriptNode),
  pythonScript: markRaw(PythonScriptNode),
  humanFeedbackNN: markRaw(HITLNode),
  customPyTorchModel: markRaw(HITLNode),
};

Object.keys(scikitUnsupervisedModels).forEach(key => {
  nodeTypes[key] = markRaw(ScikitUnsupervisedModels);
});

const scikitUnsupervisedSteps = Object.entries(scikitUnsupervisedModels)
  .map(item => {
    return {
      group: 'scikitUnsupervisedModels',
      type: item[0],
      label: t('pages.pipelines.steps.' + item[0]),
      form: item[1]
    }
  });

export const menubarSteps = [
  {
    type: 'general',
    label: t('pages.pipelines.edit.menubar.general'),
    items: [
      {
        type: 'checkoutRepo',
        label: t('pages.pipelines.steps.checkoutRepo'),
        form: checkoutRepoForm
      },
      {
        separator: true,
      },
      {
        type: 'shellScript',
        label: t('pages.pipelines.steps.shellScript'),
        form: scriptForm
      },
      {
        type: 'pythonScript',
        label: t('pages.pipelines.steps.pythonScript'),
        form: scriptForm
      }
    ]
  },
  /*
  {
    type: 'scikit',
    label: t('pages.pipelines.edit.menubar.scikit'),
    items: [
      {
        type: 'scikitDatasets',
        label: 'Scikit Datasets',
        items: [
          {
            group: 'scikitDatasets',
            type: 'scikitTrainingDataset',
            label: t('pages.pipelines.steps.scikitTrainingDataset'),
            form: scikitDatasetForm
          },
          {
            group: 'scikitDatasets',
            type: 'scikitTestingDataset',
            label: t('pages.pipelines.steps.scikitTestingDataset'),
            form: scikitDatasetForm
          },
        ]
      },
      {
        type: 'scikitUnsupervisedModels',
        label: t('pages.pipelines.steps.scikitUnsupervisedModels'),
        items: scikitUnsupervisedSteps
      },
    ]
  },
  */
  {
    type: 'hitl',
    label: t('pages.pipelines.edit.menubar.hitl'),
    items: [
      {
        group: 'hitl',
        type: 'customPyTorchModel',
        label: t('pages.pipelines.steps.customPyTorchModel'),
        form: scriptForm,
      },
      {
        group: 'hitl',
        type: 'humanFeedbackNN',
        label: t('pages.pipelines.steps.humanFeedbackNN'),
        form: hitlForm,
      },
      {
        separator: true,
      },
      {
        group: 'hitl',
        type: 'customHitl',
        label: t('pages.pipelines.steps.customHitl'),
        form: customHitlForm,
      },
    ]
  }
];

export const library = markRaw({
  ScriptEditor: ScriptEditor,
  FormFilePicker: FormFilePicker,
  BaseButton: BaseButton,
  BaseCancelAndSubmitButtons: BaseCancelAndSubmitButtons,
})