<script setup>
  import { ref, reactive, computed, watch } from "vue";
  import { storeToRefs } from "pinia";
  import {
    mdiChartTimelineVariant,
    mdiPlus
  } from "@mdi/js";
  import { useRouter } from 'vue-router';
  import SectionMain from "@/components/SectionMain.vue";
  import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
  import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
  import PipelinesTable from "@/components/PipelinesTable.vue";
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
  import router from "@/router";

  const { accessToken, requireAuthRoute } = storeToRefs(useAuthStore());
  const toast = useToast();

  // FETCH PIPELINES

  const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchPipelines } = useAsyncState(
    () => {
      return doRequest({
        url: '/api/pipeline',
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

  watch(fetchResponse, (newVal) => {
    if (newVal.error) {
      toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
    }
  })

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
            Authorization: `${accessToken.value}`,
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

  watch(createResponse, (newVal) => {
    if (newVal.error) {
      toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
    } else {
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
            Authorization: `${accessToken.value}`,
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

  watch(deleteResponse, (newVal) => {
    if (newVal.error) {
      toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
    } else {
      fetchPipelines();
    }
  })

  const pipelines = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.pipelines : []);
  const isLoading = computed(() => isFetching.value || isCreating.value || isDeleting.value);

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <PipelinesTable :items="pipelines" @deleteButtonClicked="onDeletePipelineClicked" checkable />
    </SectionMain>

    <CardBoxModal v-model="isCreateModalActive" @confirm="createPipeline(200, createPipelineForm.name)"
      title="Create Pipeline" button="success" has-cancel>
      <FormField label="Name" help="Please enter the pipeline name">
        <FormControl v-model="createPipelineForm.name" name="name" autocomplete="name" placeholder="Name"
          :focus="isCreateModalActive" />
      </FormField>
    </CardBoxModal>

    <CardBoxModal v-model="isDeleteModalActive" title="Confirm Delete" :target-id="pipelineIdToDelete"
      @confirm="deletePipeline(200, pipelineIdToDelete)" button="danger" has-cancel>
      <p>This will permanently delete this pipeline.</p>
    </CardBoxModal>

    <Toast />
  </LayoutAuthenticated>
</template>