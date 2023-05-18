<script setup>
import { ref, reactive, computed, watch } from "vue";
import {
  mdiChartTimelineVariant,
  mdiPlus
} from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import RunsTable from "@/components/RunsTable.vue";
import CardBoxModal from "@/components/CardBoxModal.vue";
import BaseButton from "@/components/BaseButton.vue";
import { useAuthStore } from "@/stores/auth";
import { doRequest } from "@/util";
import { useAsyncState } from "@vueuse/core";
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/css/index.css";
import router from "@/router";

const { accessToken, requireAuthRoute } = useAuthStore();

// FETCH PIPELINES

const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchPipelines } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/pipeline',
      method: 'GET',
      headers: {
        Authorization: `${accessToken}`
      },
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
  },
)

// CREATE PIPELINE

const isCreateModalActive = ref(false);
const onNewPipelineClicked = (e) => isCreateModalActive.value = true;

const createPipelineForm = reactive({
  name: "",
});

const { isLoading: isCreating, state: createResponse, isReady: createFinished, execute: createPipeline } = useAsyncState(
  (name) => {
    if (name && name != "") {
      return doRequest({
        url: '/api/pipeline',
        method: 'POST',
        headers: {
          Authorization: `${accessToken}`,
        },
        data: {
          name: createPipelineForm.name
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

watch(createFinished, () => {
  if (createFinished.value) {
    fetchPipelines();
  }
})

// DELETE PIPELINE

const isDeleteModalActive = ref(false);
const pipelineIdToDelete = ref(null);

const onDeletePipelineClicked = (id) => {
  isDeleteModalActive.value = true;
  pipelineIdToDelete.value = id;
}

const { isLoading: isDeleting, state: deleteResponse, isReady: deleteFinished, execute: deletePipeline } = useAsyncState(
  (pipelineID) => {
    if (pipelineID) {
      return doRequest({
        url: '/api/pipeline',
        method: 'DELETE',
        headers: {
          Authorization: `${accessToken}`,
        },
        data: {
          ID: pipelineID
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

watch(deleteFinished, () => {
  if (deleteFinished.value) {
    fetchPipelines();
  }
})

watch(fetchResponse, () => {
  if (fetchResponse.value.status === 401) {
    router.push(requireAuthRoute);
  }
})

const pipelines = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.pipelines : []);
const isLoading = computed(() => isFetching.value || isCreating.value || isDeleting.value)

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="'Runs'" main>
        <BaseButton :icon="mdiPlus" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <RunsTable :items="pipelines" @deleteButtonClicked="onDeletePipelineClicked" checkable />
    </SectionMain>

    <CardBoxModal v-model="isDeleteModalActive" title="Confirm Delete" :target-id="pipelineIdToDelete"
      @confirm="deletePipeline(200, pipelineIdToDelete)" button="danger" has-cancel>
      <p>This will permanently delete this pipeline.</p>
    </CardBoxModal>

  </LayoutAuthenticated>
</template>