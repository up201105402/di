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
import TrainersTable from "@/components/TrainersTable.vue";
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
const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchTrainers } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/trainer',
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

const createTrainerForm = reactive({
  name: "",
  entryPoint: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createTrainer } = useAsyncState(
  (name) => {
    if (name && name != "") {
      return doRequest({
        url: '/api/trainer',
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          name: createTrainerForm.name,
          entryPoint: createTrainerForm.entryPoint
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
    fetchTrainers();
  }
})

// DELETE DATASET
const isDeleteModalActive = ref(false);
const trainerIdToDelete = ref(null);

const onDeleteTrainerClicked = (id) => {
  isDeleteModalActive.value = true;
  trainerIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deleteTrainer } = useAsyncState(
  (trainerID) => {
    if (trainerID) {
      return doRequest({
        url: '/api/trainer',
        method: 'DELETE',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: trainerID
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
    fetchTrainers();
  }
})

// EDIT DATASET
const isUploadFileModalActive = ref(false);
const trainerIdToEdit = ref(-1);
const editCount = ref(0);

const onEditTrainerClicked = (id) => {
  isUploadFileModalActive.value = true;
  trainerIdToEdit.value = id;
  // trick to force the file picker to update its file's name
  editCount.value++;
}

const getTrainerFilename = (id) => {
  const trainer = trainers.value.find(trainer => trainer.ID == id)
  return trainer?.path?.replace(/^.*[\\\/]/, '')
}

const { isLoading: isEditing, state: editResponse, isReady: editFinished, execute: editTrainer } = useAsyncState(
  (trainerID) => {
    if (trainerID) {
      return doRequest({
        url: `/api/trainer/${trainerID}/file`,
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: trainerID
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
    fetchTrainers();
  }
})

const trainers = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.trainers : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isEditing.value || isDeleting.value);

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :hasButton="false" :icon="mdiDataMatrix" :title="$t('pages.trainers.header')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <TrainersTable :items="trainers" @editButtonClicked="onEditTrainerClicked" @deleteButtonClicked="onDeleteTrainerClicked" />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createTrainer(200, createTrainer.name)" :title="$t('pages.trainers.dialog.create.header')" button="success" has-cancel>
        <FormField :label="$t('pages.trainers.dialog.create.name.label')" :help="$t('pages.trainers.dialog.create.name.help')">
            <FormControl v-model="createTrainerForm.name" name="name" autocomplete="name" placeholder="Name" :focus="isCreateModalActive" />
        </FormField>
    </CardBoxModal>

    <CardBoxModal v-model="isUploadFileModalActive" :title="$t('pages.trainers.dialog.create.header')" button="success" has-cancel @confirm="fetchTrainers()">
      <FormFilePicker :key="'picker_' + editCount" id="script-file-upload" :filename="getTrainerFilename(trainerIdToEdit)" :url="`/api/trainer/${trainerIdToEdit}/file`" :label="t('pages.pipelines.edit.dialog.stepConfig.scriptFile.button')" />
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" :title="$t('pages.trainers.dialog.delete.header')"
      :target-id="trainerIdToDelete" @confirm="deleteTrainer(200, trainerIdToDelete)" button="danger" has-cancel>
      <p>{{ $t('pages.trainers.dialog.delete.body') }}</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>