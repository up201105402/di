<script setup>
import { Panel, PanelPosition, VueFlow, isNode, useVueFlow } from '@vue-flow/core';
import { storeToRefs } from "pinia";
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import { MiniMap } from '@vue-flow/minimap';
import { ref, computed, watch } from "vue";
import { useAsyncState } from "@vueuse/core";
import { doRequest, parsePipelineDefinition, deepFilterMenuBarSteps, validateCron } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRouter, useRoute } from 'vue-router';
import { mdiChartTimelineVariant, mdiPlus, mdiCalendarEdit } from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import Menubar from 'primevue/menubar';
import CardBoxModal from '@/components/CardBoxModal.vue';
import UpsertStepDialog from '@/components/UpsertStepDialog.vue';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import Loading from "vue-loading-overlay";
import { nodeTypes, menubarSteps } from "@/pipelines/steps";
import deepEqual from 'deep-equal';
import $ from 'jquery';
import PipelineScheduleTable from '@/components/PipelineScheduleTable.vue';
import TabView from 'primevue/tabview';
import TabPanel from 'primevue/tabpanel';
import Dropdown from 'primevue/dropdown';
import Calendar from 'primevue/calendar';
import InputText from 'primevue/inputtext';
import { datasetForm } from '@/pipelines/steps/datasets'
import { trainerForm } from '@/pipelines/steps/trainers'
import { useField } from 'vee-validate';
import cronTime from "cron-time-generator";
import { i18n } from '@/i18n';

const { t } = i18n.global;

const { accessToken } = storeToRefs(useAuthStore());
const router = useRouter();
const route = useRoute();
const elements = ref([]);
const pipelineSchedules = ref([]);
const pipelineTitle = ref('');
const toast = useToast();

const { value: cronValue, errorMessage: cronError } = useField('input-cron', validateCron);

const emit = defineEmits(["onUpdate", "onStepEdited"]);

// FETCH PIPELINE
const { isLoading: isFetchingPipeline, state: fetchPipelineResponse, isReady: isFetchPipelineFinished, execute: fetchPipeline } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/pipeline/${route.params.id}`,
      method: 'GET',
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
  },
);

// FETCH PIPELINE SCHEDULE
const { isLoading: isFetchingPipelineSchedules, state: fetchPipelineSchedulesResponse, isReady: isPipelineScheduleFetchFinished, execute: fetchPipelineSchedules } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/pipeline/${route.params.id}/schedule`,
      method: 'GET',
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
  },
);

// UPDATE PIPELINE
const { isLoading: isUpdatingPipeline, state: updatePipelineResponse, isReady: isPipelineUpdateFinished, execute: updatePipeline } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/pipeline/${route.params.id}`,
      method: 'POST',
      data: {
        id: parseInt(route.params.id),
        definition: JSON.stringify(elements.value)
      },
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
)

// CREATE PIPELINE SCHEDULE
const { isLoading: isCreatingPipelineSchedule, state: createPipelineScheduleResponse, isReady: isPipelineScheduleCreateFinished, execute: createPipelineSchedule } = useAsyncState(
  (pipelineSchedule) => {
    return doRequest({
      url: `/api/pipeline/${route.params.id}/schedule`,
      method: 'POST',
      data: {
        ...pipelineSchedule
      },
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
);

// DELETE PIPELINE SCHEDULE
const { isLoading: isDeletingPipelineSchedule, state: deletePipelineScheduleResponse, isReady: isPipelineScheduleDeleteFinished, execute: deletePipelineSchedule } = useAsyncState(
  (pipelineSchedule) => {
    return doRequest({
      url: `/api/pipeline/${route.params.id}/schedule`,
      method: 'DELETE',
      data: {
        ...pipelineSchedule
      },
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
);

// FETCH DATASETS
const { isLoading: isFetchingDatasets, state: fetchDatasetsResponse, isReady: isFetchDatasetsFinished, execute: fetchDatasets } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/dataset`,
      method: 'GET',
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    immediate: true,
    delay: 500,
    resetOnExecute: false,
  },
);

// FETCH MODELS
const { isLoading: isFetchingTrainers, state: fetchTrainersResponse, isReady: isFetchModelsFinished, execute: fetchTrainers } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/trainer`,
      method: 'GET',
      headers: {
        Authorization: `${accessToken.value}`
      },
    })
  },
  {},
  {
    immediate: true,
    delay: 500,
    resetOnExecute: false,
  },
);

const selectedStep = ref();

const onMenubarClick = (e) => {
  editStepNodeId.value = "-1";
  stepData.value = {
    pipelineID: route.params.id
  };
  selectedStep.value = e.item

  if (e.item) {
    isStepDialogActive.value = true;
    if (elements.value == null || elements.value.length == 0) {
      stepData.value = {
        nameAndType: {
          isFirstStep: true
        },
      }
    }
    let step = e.item.form(stepData.value, onStepEdited);
    formSchema.value = step.formSchema;
    dialogTitle.value = t('pages.pipelines.edit.dialog.create.header', { name: e.item.label });
    stepData.value = step.formkitData;
    stepData.value.group = e.item.group;
    stepData.value.type = e.item.type;
    count++;
  }
}

const addItemCommand = (item) => {
  if (item.items) {
    item.items.forEach(addItemCommand)
  } else {
    item.command = onMenubarClick
  }
}

watch(fetchDatasetsResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    const datasetGroupsIndex = menubarItems.value?.findIndex(group => group.type == 'datasets');
    if (datasetGroupsIndex > -1) {
      value.data?.datasets?.forEach(dataset => menubarItems.value[datasetGroupsIndex].items.push({
        group: 'datasets',
        type: 'dataset',
        label: dataset.name,
        form: datasetForm(dataset),
        command: onMenubarClick,
      }))
    }
  }
})

watch(fetchTrainersResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    const trainerGroupsIndex = menubarItems.value?.findIndex(group => group.type == 'trainers');
    if (trainerGroupsIndex > -1) {
      value.data?.trainers?.forEach(trainer => menubarItems.value[trainerGroupsIndex].items.push({
        group: 'trainers',
        type: 'trainer',
        label: trainer.name,
        form: trainerForm(trainer),
        command: onMenubarClick,
      }))
    }
  }
})

menubarSteps.forEach(addItemCommand)
const menubarItems = ref(menubarSteps)

watch(fetchPipelineResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  }
})

watch(fetchPipelineSchedulesResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  }
})

watch(isFetchPipelineFinished, () => {
  elements.value = fetchPipelineResponse.value?.data ? parsePipelineDefinition(fetchPipelineResponse.value.data.pipeline, toast) : [];
  pipelineTitle.value = fetchPipelineResponse.value?.data ? fetchPipelineResponse.value.data.pipeline.name : 'Untitled';
})

watch(isPipelineScheduleFetchFinished, () => {
  pipelineSchedules.value = fetchPipelineSchedulesResponse.value?.data ? fetchPipelineSchedulesResponse.value.data.schedules : [];
})

watch(updatePipelineResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  }

  if (value.status === 200) {
    hasChanges.value = false;
    router.push("/pipelines");
  }
})

watch(createPipelineScheduleResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  }

  if (value.status === 200) {
    fetchPipelineSchedules(null);
  }
})

watch(deletePipelineScheduleResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  }

  if (value.status === 200) {
    fetchPipelineSchedules(null);
  }
})

const isLoading = computed(() => isFetchingPipeline.value || isFetchingPipelineSchedules.value || isUpdatingPipeline.value || isCreatingPipelineSchedule.value || isFetchingDatasets.value || isFetchingTrainers.value);

const hasChanges = ref(false);
const isStepDialogActive = ref(false);
const isScheduleDialogActive = ref(false);
const isDeleteScheduleDialogActive = ref(false);
const pipelineScheduleIdToDelete = ref();
const formSchema = ref({});
const dialogTitle = ref('');
const stepData = ref({});
const editStepNodeId = ref("-1");
let count = 0;

const activeTabIndex = ref(0);
const weekDays = [
  { value: 'monday', label: 'Monday' },
  { value: 'tuesday', label: 'Tuesday' },
  { value: 'wednesday', label: 'Wednesday' },
  { value: 'thursday', label: 'Thursday' },
  { value: 'friday', label: 'Friday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' },
];
const selectedWeekDay = ref(weekDays[0].value);
const scheduledTime = ref(new Date(new Date(Date.now() + 60 * 60 * 1000).setSeconds(0))); // one hour in the future

onBeforeRouteLeave((to, from) => {
  if (hasChanges.value) {
    const answer = window.confirm(
      'Do you really want to leave? you have unsaved changes!'
    )
    // cancel the navigation and stay on the same page
    if (!answer) return false
  }
})

onBeforeRouteUpdate((to, from) => {
  if (hasChanges.value) {
    const answer = window.confirm(
      'Do you really want to leave? you have unsaved changes!'
    )
    // cancel the navigation and stay on the same page
    if (!answer) return false
  }
})

const onCreateScheduleClick = (e) => {
  isScheduleDialogActive.value = true;
}

const onDeletePipelineScheduleClick = (id) => {
  pipelineScheduleIdToDelete.value = id;
  isDeleteScheduleDialogActive.value = true;
}

const onCreateSchedule = (e) => {
  var pipelineSchedule = {}

  if (activeTabIndex.value == 0) {
    pipelineSchedule['uniqueOccurrence'] = scheduledTime.value.toISOString()
  } else if (activeTabIndex.value == 1) {
    pipelineSchedule['cronExpression'] = cronTime.onSpecificDaysAt([selectedWeekDay.value], scheduledTime.value.getHours(), scheduledTime.value.getMinutes())
  } else if (activeTabIndex.value == 2) {
    pipelineSchedule['cronExpression'] = cronValue.value;
  }

  createPipelineSchedule(null, pipelineSchedule);
}

const onDeleteSchedule = (id) => {
  deletePipelineSchedule(null, { ID: id })
}

const onNodeDoubleClick = (e) => {
  isStepDialogActive.value = true;
  editStepNodeId.value = e.node.id;
  stepData.value = { ...e.node.data }
  const formkitObject = deepFilterMenuBarSteps(menubarSteps, 'type', e.node.data.type).form({ ...e.node.data }, onStepEdited);
  formSchema.value = formkitObject.formSchema;
  const label = t('pages.pipelines.steps.' + e.node.data.type);
  dialogTitle.value = t('pages.pipelines.edit.dialog.edit.header', { name: label });
  stepData.value = formkitObject.formkitData;
  count++;
}

$(document).off("onNodeDelete");
$(document).on("onNodeDelete", function (e, details) {
  onNodeDelete(details.id);
})

const onNodeDelete = (id) => {
  const index = elements.value.findIndex(element => element.id === id);

  if (index > -1) {
    elements.value.splice(index, 1);
    hasChanges.value = true;
    count++;
  }
}

const getNextId = () => {
  const ids = elements.value.filter(element => isNode(element)).map(element => parseInt(element.id));
  return ids.length == 0 ? '0' : (Math.max.apply(Math, ids) + 1) + '';
}

const onStepCreate = (formData) => {
  const newId = getNextId();

  if (elements.value.length == 0) {
    formData.nameAndType.isFirstStep = true;
  } else if (formData.nameAndType.isFirstStep == null) {
    formData.nameAndType.isFirstStep = false;
  }

  elements.value.push({
    id: newId,
    label: formData.nameAndType.name,
    type: selectedStep.value.type,
    position: { x: 0, y: 0 },
    class: 'light',
    data: {
      ...formData,
      id: newId,
      group: selectedStep.value.group,
      model: selectedStep.value.type,
      type: selectedStep.value.type,
    },
  });
  isStepDialogActive.value = false;
  hasChanges.value = true;
  count++;
}

const onCancel = () => {
  isStepDialogActive.value = false;
  count++;
}

const onStepEdited = (formData) => {
  if (formData.nameAndType.isFirstStep) {
    elements.value.forEach(element => {
      if (element.data.nameAndType) {
        element.data.nameAndType.isFirstStep = false;
      }
    })
  }

  const index = elements.value.findIndex(element => element.id === formData.id);

  if (index === -1) {
    onStepCreate(formData);
  } else {
    const oldStepData = elements.value[index].data;
    const oldStepType = elements.value[index].type;
    const newStepData = {
      ...formData,
    };

    if (!deepEqual(oldStepData, newStepData)) {
      elements.value = elements.value.toSpliced(index, 1, {
        id: formData.id,
        type: oldStepType,
        label: formData.nameAndType.name,
        position: { x: 0, y: 0 },
        class: 'light',
        data: {
          ...formData,
        },
      });

      hasChanges.value = true
    }

    isStepDialogActive.value = false;
    count++;
  }
}

const onPipelineSave = () => {
  updatePipeline();
}

const onPipelineCancel = () => {
  hasChanges.value = false;
  router.push('/pipelines');
}

/**
   * useVueFlow provides all event handlers and store properties
   * You can pass the composable an object that has the same properties as the VueFlow component props
   */
const { onPaneReady, onNodesChange, onEdgesChange, onConnect, addEdges, isEdge, setTransform, toObject } = useVueFlow();

/**
 * This is a Vue Flow event-hook which can be listened to from anywhere you call the composable, instead of only on the main component
 *
 * onPaneReady is called when viewpane & nodes have visible dimensions
 */
onPaneReady(({ fitView }) => {
  fitView();
})

const onEdgeUpdate = (edge) => emit("onUpdate", elements.value);

/**
 * onConnect is called when a new connection is created.
 * You can add additional properties to your new edge (like a type or label) or block the creation altogether
 */
onConnect((edge) => {
  edge.updatable = true;
  edge.type = 'smoothstep';
  addEdges([edge]);
  hasChanges.value = true;
  emit("onUpdate", elements.value);
})

onNodesChange(events => {
  const removedNodeOrEdge = events.find(event => event.type == 'remove');
  if (removedNodeOrEdge != null) {
    hasChanges.value = true;
  }
})

onEdgesChange(events => {
  const removedNodeOrEdge = events.find(event => event.type == 'remove');
  if (removedNodeOrEdge != null) {
    hasChanges.value = true;
  }
})

const dark = ref(false)

/**
 * To update node properties you can simply use your elements v-model and mutate the elements directly
 * Changes should always be reflected on the graph reactively, without the need to overwrite the elements
 */
function updatePos() {
  return elements.value.forEach((el) => {
    if (isNode(el)) {
      el.position = {
        x: Math.random() * 400,
        y: Math.random() * 400,
      }
    }
  })
}

/**
 * toObject transforms your current graph data to an easily persist-able object
 */
function logToObject() {
  return console.log(toObject())
}

/**
 * Resets the current viewpane transformation (zoom & pan)
 */
function resetTransform() {
  return setTransform({ x: 0, y: 0, zoom: 1 })
}

function toggleClass() {
  return (dark.value = !dark.value)
}

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />
      <SectionTitleLineWithButton :hasButton="false" :icon="mdiChartTimelineVariant" :title="pipelineTitle" main>
        <div>
          <BaseButtons style="float:right">
            <BaseButton :disabled="!hasChanges" :label="$t('buttons.save')" color="success" @click="onPipelineSave" />
            <BaseButton :label="$t('buttons.cancel')" color="danger" @click="onPipelineCancel" />
          </BaseButtons>
        </div>
      </SectionTitleLineWithButton>
      <Menubar :model="menubarItems" style="margin-right: 10px; " />
      <VueFlow v-model="elements" :class="{ dark }" class="basicflow" :node-types="nodeTypes"
        @nodeDoubleClick="onNodeDoubleClick" @edge-update="onEdgeUpdate" :default-viewport="{ zoom: 1.5 }" :min-zoom="0.2"
        :max-zoom="4">
        <Background :pattern-color="dark ? '#FFFFFB' : '#aaa'" gap="8" />
        <MiniMap />
        <Controls />

        <Panel :position="PanelPosition.TopRight" class="controls">
          <button style="background-color: #113285; color: white" title="Reset Transform" @click="resetTransform">
            <svg width="16" height="16" viewBox="0 0 32 32">
              <path fill="#FFFFFB"
                d="M18 28A12 12 0 1 0 6 16v6.2l-3.6-3.6L1 20l6 6l6-6l-1.4-1.4L8 22.2V16a10 10 0 1 1 10 10Z" />
            </svg>
          </button>

          <button style="background-color: #6f3381" title="Shuffle Node Positions" @click="updatePos">
            <svg width="16" height="16" viewBox="0 0 24 24">
              <path fill="#FFFFFB"
                d="M14 20v-2h2.6l-3.2-3.2l1.425-1.425L18 16.55V14h2v6Zm-8.6 0L4 18.6L16.6 6H14V4h6v6h-2V7.4Zm3.775-9.425L4 5.4L5.4 4l5.175 5.175Z" />
            </svg>
          </button>

          <button :style="{ backgroundColor: dark ? '#FFFFFB' : '#292524', color: dark ? '#292524' : '#FFFFFB' }"
            @click="toggleClass">
            <template v-if="dark">
              <svg width="16" height="16" viewBox="0 0 24 24">
                <path fill="#292524"
                  d="M12 17q-2.075 0-3.537-1.463Q7 14.075 7 12t1.463-3.538Q9.925 7 12 7t3.538 1.462Q17 9.925 17 12q0 2.075-1.462 3.537Q14.075 17 12 17ZM2 13q-.425 0-.712-.288Q1 12.425 1 12t.288-.713Q1.575 11 2 11h2q.425 0 .713.287Q5 11.575 5 12t-.287.712Q4.425 13 4 13Zm18 0q-.425 0-.712-.288Q19 12.425 19 12t.288-.713Q19.575 11 20 11h2q.425 0 .712.287q.288.288.288.713t-.288.712Q22.425 13 22 13Zm-8-8q-.425 0-.712-.288Q11 4.425 11 4V2q0-.425.288-.713Q11.575 1 12 1t.713.287Q13 1.575 13 2v2q0 .425-.287.712Q12.425 5 12 5Zm0 18q-.425 0-.712-.288Q11 22.425 11 22v-2q0-.425.288-.712Q11.575 19 12 19t.713.288Q13 19.575 13 20v2q0 .425-.287.712Q12.425 23 12 23ZM5.65 7.05L4.575 6q-.3-.275-.288-.7q.013-.425.288-.725q.3-.3.725-.3t.7.3L7.05 5.65q.275.3.275.7q0 .4-.275.7q-.275.3-.687.287q-.413-.012-.713-.287ZM18 19.425l-1.05-1.075q-.275-.3-.275-.712q0-.413.275-.688q.275-.3.688-.287q.412.012.712.287L19.425 18q.3.275.288.7q-.013.425-.288.725q-.3.3-.725.3t-.7-.3ZM16.95 7.05q-.3-.275-.287-.688q.012-.412.287-.712L18 4.575q.275-.3.7-.288q.425.013.725.288q.3.3.3.725t-.3.7L18.35 7.05q-.3.275-.7.275q-.4 0-.7-.275ZM4.575 19.425q-.3-.3-.3-.725t.3-.7l1.075-1.05q.3-.275.713-.275q.412 0 .687.275q.3.275.288.688q-.013.412-.288.712L6 19.425q-.275.3-.7.287q-.425-.012-.725-.287Z" />
              </svg>
            </template>

            <template v-else>
              <svg width="16" height="16" viewBox="0 0 24 24">
                <path fill="#FFFFFB"
                  d="M12 21q-3.75 0-6.375-2.625T3 12q0-3.75 2.625-6.375T12 3q.35 0 .688.025q.337.025.662.075q-1.025.725-1.637 1.887Q11.1 6.15 11.1 7.5q0 2.25 1.575 3.825Q14.25 12.9 16.5 12.9q1.375 0 2.525-.613q1.15-.612 1.875-1.637q.05.325.075.662Q21 11.65 21 12q0 3.75-2.625 6.375T12 21Z" />
              </svg>
            </template>
          </button>

          <button title="Log `toObject`" @click="logToObject">
            <svg width="16" height="16" viewBox="0 0 24 24">
              <path fill="#292524"
                d="M20 19V7H4v12h16m0-16a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h16m-7 14v-2h5v2h-5m-3.42-4L5.57 9H8.4l3.3 3.3c.39.39.39 1.03 0 1.42L8.42 17H5.59l3.99-4Z" />
            </svg>
          </button>
        </Panel>
      </VueFlow>
      <CardBoxModal v-model="isStepDialogActive" :has-submit="false" :has-cancel="false" :title="dialogTitle"
        @cancel="count++">
        <UpsertStepDialog :key="'createStepDialog_' + count" :formSchema="formSchema" :nodeId="editStepNodeId"
          :nodeData="stepData" @onSubmit="onStepEdited" @onCancel="onCancel" />
      </CardBoxModal>
      <SectionTitleLineWithButton :hasButton="false" :icon="mdiCalendarEdit"
        :title="$t('pages.pipelines.edit.scheduling.header')">
        <BaseButton :label="$t('pages.pipelines.edit.scheduling.add')" :icon="mdiPlus" color="success"
          @click="onCreateScheduleClick" />
      </SectionTitleLineWithButton>
      <PipelineScheduleTable :items="pipelineSchedules" @delete-button-clicked="onDeletePipelineScheduleClick" />
      <CardBoxModal v-model="isScheduleDialogActive" has-cancel title="Create Schedule" @confirm="onCreateSchedule">
        <TabView :activeIndex="activeTabIndex">
          <TabPanel header="Once">
            <Calendar v-model="scheduledTime" showTime hourFormat="12" />
          </TabPanel>
          <TabPanel header="Repeat">
            <Dropdown style="max-width: 30%; margin-right: 10px;" v-model="selectedWeekDay" :options="weekDays"
              optionValue="value" optionLabel="label" placeholder="Select a day" class="w-full md:w-14rem" />
            <Calendar v-model="scheduledTime" timeOnly hourFormat="12" />
          </TabPanel>
          <TabPanel header="Cron Expression">
            <InputText id="input-cron" style="margin-right: 10px;" v-model="cronValue" type="text"
              placeholder="Cron expression" :class="{ 'p-invalid': cronError }" aria-describedby="text-error" />
            <span class="p-error" id="cron-error">{{ cronError || '&nbsp;' }}</span>
          </TabPanel>
        </TabView>
      </CardBoxModal>
      <CardBoxModal v-model="isDeleteScheduleDialogActive" :target-id="pipelineScheduleIdToDelete" has-cancel submitLabel="Confirm" title="Delete Schedule" @confirm="onDeleteSchedule" />
    </SectionMain>
    <Toast />
  </LayoutAuthenticated>
</template>@/pipelines/steps/trainer