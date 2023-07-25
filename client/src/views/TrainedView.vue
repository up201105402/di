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
import TrainedTable from "@/components/TrainedTable.vue";
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
const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchTrained } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/trained',
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

const createTrainedForm = reactive({
  name: "",
  entryPoint: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createTrained } = useAsyncState(
  (name) => {
    if (name && name != "") {
      return doRequest({
        url: '/api/trained',
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          name: createTrainedForm.name,
          entryPoint: createTrainedForm.entryPoint
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
    fetchTrained();
  }
})

// DELETE DATASET
const isDeleteModalActive = ref(false);
const trainedIdToDelete = ref(null);

const onDeleteTrainedClicked = (id) => {
  isDeleteModalActive.value = true;
  trainedIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deleteTrained } = useAsyncState(
  (trainedID) => {
    if (trainedID) {
      return doRequest({
        url: '/api/trained',
        method: 'DELETE',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: trainedID
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
    fetchTrained();
  }
})

// EDIT DATASET
const isUploadFileModalActive = ref(false);
const trainedIdToEdit = ref(-1);
const editCount = ref(0);

const onEditTrainedClicked = (id) => {
  isUploadFileModalActive.value = true;
  trainedIdToEdit.value = id;
  // trick to force the file picker to update its file's name
  editCount.value++;
}

const getTrainedFilename = (id) => {
  const trainedModel = trained.value.find(model => model.ID == id)
  return trainedModel?.path?.replace(/^.*[\\\/]/, '')
}

const { isLoading: isEditing, state: editResponse, isReady: editFinished, execute: editTrained } = useAsyncState(
  (trainedID) => {
    if (trainedID) {
      return doRequest({
        url: `/api/trained/${trainedID}/file`,
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: trainedID
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
    fetchTrained();
  }
})

const trained = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.trained : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isEditing.value || isDeleting.value);

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :hasButton="false" :icon="mdiDataMatrix" :title="$t('pages.trained.header')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <TrainedTable :items="trained" @editButtonClicked="onEditTrainedClicked" @deleteButtonClicked="onDeleteTrainedClicked" />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createTrained(200, createTrained.name)" :title="$t('pages.trained.dialog.create.header')" button="success" has-cancel>
        <FormField :label="$t('pages.trained.dialog.create.name.label')" :help="$t('pages.trained.dialog.create.name.help')">
            <FormControl v-model="createTrainedForm.name" name="name" autocomplete="name" placeholder="Name" :focus="isCreateModalActive" />
        </FormField>
    </CardBoxModal>

    <CardBoxModal v-model="isUploadFileModalActive" :title="$t('pages.trained.dialog.create.header')" button="success" has-cancel @confirm="fetchTrained()">
      <FormFilePicker :key="'picker_' + editCount" id="script-file-upload" :filename="getTrainedFilename(trainedIdToEdit)" :url="`/api/trained/${trainedIdToEdit}/file`" :label="t('pages.pipelines.edit.dialog.stepConfig.scriptFile.button')" />
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" :title="$t('pages.trained.dialog.delete.header')"
      :target-id="trainedIdToDelete" @confirm="deleteTrained(200, trainedIdToDelete)" button="danger" has-cancel>
      <p>{{ $t('pages.trained.dialog.delete.body') }}</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>