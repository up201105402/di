<script setup>
import { ref, reactive, computed, watch } from "vue";
import { storeToRefs } from "pinia";
import {
  mdiDataMatrix,
  mdiPlus
} from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import DatasetsTable from "@/components/DatasetsTable.vue";
import CardBoxModal from "@/components/CardBoxModal.vue";
import BaseButton from "@/components/BaseButton.vue";
import FormControl from "@/components/FormControl.vue";
import FormField from "@/components/FormField.vue";
import { useAuthStore } from "@/stores/auth";
import { doRequest } from "@/util";
import { useAsyncState } from "@vueuse/core";
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/css/index.css";
import { i18n } from '@/i18n';
import FormFilePicker from "@/components/FormFilePicker.vue";

const { t } = i18n.global;
const { accessToken } = storeToRefs(useAuthStore());
const toast = useToast();

// FETCH DATASETS
const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchDatasets } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/dataset',
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
)

watch(fetchResponse, (value) => {
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

// CREATE DATASET
const isCreateModalActive = ref(false);
const onNewPipelineClicked = (e) => isCreateModalActive.value = true;

const createDatasetForm = reactive({
  name: "",
  entryPoint: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createDataset } = useAsyncState(
  (name) => {
    if (name && name != "") {
      return doRequest({
        url: '/api/dataset',
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          name: createDatasetForm.name,
          entryPoint: createDatasetForm.entryPoint
        },
      });
    }

    return {};
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
)

watch(createResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    fetchDatasets();
  }
})

// DELETE DATASET
const isDeleteModalActive = ref(false);
const datasetIdToDelete = ref(null);

const onDeleteDatasetClicked = (id) => {
  isDeleteModalActive.value = true;
  datasetIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deleteDataset } = useAsyncState(
  (datasetID) => {
    if (datasetID) {
      return doRequest({
        url: '/api/dataset',
        method: 'DELETE',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: datasetID
        },
      });
    }

    return {};
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
)

watch(deleteResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    fetchDatasets();
  }
})

// EDIT DATASET
const isUploadFileModalActive = ref(false);
const datasetIdToEdit = ref(-1);
const editCount = ref(0);

const onEditDatasetClicked = (id) => {
  isUploadFileModalActive.value = true;
  datasetIdToEdit.value = id;
  // trick to force the file picker to update its file's name
  editCount.value++;
}

const getDatasetFilename = (id) => {
  const dataset = datasets.value.find(dataset => dataset.ID == id)
  return dataset?.path?.replace(/^.*[\\\/]/, '')
}

const { isLoading: isEditing, state: editResponse, isReady: editFinished, execute: editPipeline } = useAsyncState(
  (datasetID) => {
    if (datasetID) {
      return doRequest({
        url: `/api/dataset/${datasetID}/file`,
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: datasetID
        },
      });
    }

    return {};
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
)

watch(deleteResponse, (value) => {
  if (value.error) {
    let header = t('global.errors.generic.header');
    let detail = value.error;

    if (value.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    fetchDatasets();
  }
})

const datasets = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.datasets : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isEditing.value || isDeleting.value);

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :hasButton="false" :icon="mdiDataMatrix" :title="$t('pages.datasets.header')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <DatasetsTable :items="datasets" @editButtonClicked="onEditDatasetClicked" @deleteButtonClicked="onDeleteDatasetClicked" />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createDataset(200, createDatasetForm.name)" :title="$t('pages.datasets.dialog.create.header')" button="success" has-cancel>
        <FormField :label="$t('pages.datasets.dialog.create.name.label')" :help="$t('pages.datasets.dialog.create.name.help')">
            <FormControl v-model="createDatasetForm.name" name="name" autocomplete="name" placeholder="Name" :focus="isCreateModalActive" />
        </FormField>
    </CardBoxModal>

    <CardBoxModal v-model="isUploadFileModalActive" :title="$t('pages.datasets.dialog.create.header')" button="success" has-cancel @confirm="fetchDatasets()">
      <FormFilePicker :key="'picker_' + editCount" id="script-file-upload" :filename="getDatasetFilename(datasetIdToEdit)" :url="`/api/dataset/${datasetIdToEdit}/file`" :label="t('pages.pipelines.edit.dialog.stepConfig.scriptFile.button')" />
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" :title="$t('pages.datasets.dialog.delete.header')"
      :target-id="datasetIdToDelete" @confirm="deleteDataset(200, datasetIdToDelete)" button="danger" has-cancel>
      <p>{{ $t('pages.datasets.dialog.delete.body') }}</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>