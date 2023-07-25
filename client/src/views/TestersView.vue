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
import TestersTable from "@/components/TestersTable.vue";
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
const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchTesters } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/tester',
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

const createTesterForm = reactive({
  name: "",
  entryPoint: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createTester } = useAsyncState(
  (name) => {
    if (name && name != "") {
      return doRequest({
        url: '/api/tester',
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          name: createTesterForm.name,
          entryPoint: createTesterForm.entryPoint
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
    fetchTesters();
  }
})

// DELETE DATASET
const isDeleteModalActive = ref(false);
const testerIdToDelete = ref(null);

const onDeleteTesterClicked = (id) => {
  isDeleteModalActive.value = true;
  testerIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deleteTester } = useAsyncState(
  (testerID) => {
    if (testerID) {
      return doRequest({
        url: '/api/tester',
        method: 'DELETE',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: testerID
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
    fetchTesters();
  }
})

// EDIT DATASET
const isUploadFileModalActive = ref(false);
const testerIdToEdit = ref(-1);
const editCount = ref(0);

const onEditTesterClicked = (id) => {
  isUploadFileModalActive.value = true;
  testerIdToEdit.value = id;
  // trick to force the file picker to update its file's name
  editCount.value++;
}

const getTesterFilename = (id) => {
  const tester = testers.value.find(tester => tester.ID == id)
  return tester?.path?.replace(/^.*[\\\/]/, '')
}

const { isLoading: isEditing, state: editResponse, isReady: editFinished, execute: editTester } = useAsyncState(
  (testerID) => {
    if (testerID) {
      return doRequest({
        url: `/api/tester/${testerID}/file`,
        method: 'POST',
        headers: {
          Authorization: `${accessToken.value}`,
        },
        data: {
          ID: testerID
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
    fetchTesters();
  }
})

const testers = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.testers : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isEditing.value || isDeleting.value);

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :hasButton="false" :icon="mdiDataMatrix" :title="$t('pages.testers.header')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <TestersTable :items="testers" @editButtonClicked="onEditTesterClicked" @deleteButtonClicked="onDeleteTesterClicked" />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createTester(200, createTester.name)" :title="$t('pages.testers.dialog.create.header')" button="success" has-cancel>
        <FormField :label="$t('pages.testers.dialog.create.name.label')" :help="$t('pages.testers.dialog.create.name.help')">
            <FormControl v-model="createTesterForm.name" name="name" autocomplete="name" placeholder="Name" :focus="isCreateModalActive" />
        </FormField>
    </CardBoxModal>

    <CardBoxModal v-model="isUploadFileModalActive" :title="$t('pages.testers.dialog.create.header')" button="success" has-cancel @confirm="fetchTesters()">
      <FormFilePicker :key="'picker_' + editCount" id="script-file-upload" :filename="getTesterFilename(testerIdToEdit)" :url="`/api/tester/${testerIdToEdit}/file`" :label="t('pages.pipelines.edit.dialog.stepConfig.scriptFile.button')" />
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" :title="$t('pages.testers.dialog.delete.header')"
      :target-id="testerIdToDelete" @confirm="deleteTester(200, testerIdToDelete)" button="danger" has-cancel>
      <p>{{ $t('pages.testers.dialog.delete.body') }}</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>