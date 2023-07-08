<script setup>
import { computed, watch } from "vue";
import { storeToRefs } from "pinia";
import { mdiRunFast } from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import PipelinesRunsTable from "@/components/PipelinesRunsTable.vue";
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

const pipelines = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.pipelines : []);
const isLoading = computed(() => isFetching.value)

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />
      <SectionTitleLineWithButton :icon="mdiRunFast" :title="'Runs'" main />
      <PipelinesRunsTable :items="pipelines" checkable />
    </SectionMain>
    <Toast />
  </LayoutAuthenticated>
</template>