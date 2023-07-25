<script setup>
import { ref, reactive, computed, watch } from "vue";
import { storeToRefs } from "pinia";
import {
  mdiDataMatrix,
  mdiPlus
} from "@mdi/js";
import { useRoute } from 'vue-router';
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import DatasetScriptTable from "@/components/DatasetScriptTable.vue";
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
import FileUpload from "primevue/fileupload";

const { t } = i18n.global;
const { accessToken } = storeToRefs(useAuthStore());
const route = useRoute();
const toast = useToast();

// FETCH DATASET
const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchDatasets } = useAsyncState(
  () => {
    return doRequest({
      url: `/api/dataset/${route.params.id}`,
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

// CREATE DATASET SCRIPT
const isCreateModalActive = ref(false);
const onNewPipelineClicked = (e) => isCreateModalActive.value = true;

const createDatasetForm = reactive({
  name: "",
  entryPoint: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createPipeline } = useAsyncState(
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
const datasetScriptIdToDelete = ref(null);

const onDeleteDatasetClicked = (id) => {
  isDeleteModalActive.value = true;
  datasetScriptIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deletePipeline } = useAsyncState(
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

const datasetScripts = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.datasetScripts : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isDeleting.value);

const customUploader = async (event) => {
    const file = event.files[0];
};


</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :hasButton="false" :icon="mdiDataMatrix" :title="$t('pages.datasets.header')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <DatasetScriptTable :items="datasetScripts" @deleteButtonClicked="onDeleteDatasetClicked" />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createPipeline(200, createDatasetForm.name)" :title="$t('pages.datasets.dialog.create.header')" button="success" has-cancel>
      <FileUpload name="demo[]" url="./upload.php" customUpload @uploader="customUploader" :multiple="true" accept="text/plain" :maxFileSize="1000000">
        <template #empty>
            <p>Drag and drop files to here to upload.</p>
        </template>
    </FileUpload>
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" :title="$t('pages.datasets.dialog.delete.header')"
      :target-id="datasetScriptIdToDelete" @confirm="deletePipeline(200, datasetScriptIdToDelete)" button="danger" has-cancel>
      <p>{{ $t('pages.datasets.dialog.delete.body') }}</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>