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
  customHitl: markRaw(HITLNode),
  dataset: markRaw(HITLNode),
  trainer: markRaw(HITLNode),
  trained: markRaw(HITLNode),
  tester: markRaw(HITLNode),
};

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
  {
    type: 'hitl',
    label: t('pages.pipelines.edit.menubar.hitl'),
    items: [
      {
        group: 'hitl',
        type: 'humanFeedbackNN',
        label: t('pages.pipelines.steps.humanFeedbackNN'),
        form: hitlForm,
      },
      {
        separator: true,
      },
    ]
  },
  {
    type: 'datasets',
    label: t('pages.pipelines.edit.menubar.datasets'),
    items: [

    ]
  },
  {
    type: 'trainers',
    label: t('pages.pipelines.edit.menubar.trainers'),
    items: [
      
    ]
  },
  {
    type: 'trainedModels',
    label: t('pages.pipelines.edit.menubar.trainedModels'),
    items: [
      
    ]
  },
  {
    type: 'testers',
    label: t('pages.pipelines.edit.menubar.testers'),
    items: [
      
    ]
  }
];

export const library = markRaw({
  ScriptEditor: ScriptEditor,
  FormFilePicker: FormFilePicker,
  BaseButton: BaseButton,
  BaseCancelAndSubmitButtons: BaseCancelAndSubmitButtons,
})