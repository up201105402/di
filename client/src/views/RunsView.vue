<script setup>
import { computed, watch } from "vue";
import { storeToRefs } from "pinia";
import { mdiChartTimelineVariant } from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import PipelinesRunsTable from "@/components/PipelinesRunsTable.vue";
import { useAuthStore } from "@/stores/auth";
import { doRequest } from "@/util";
import { useAsyncState } from "@vueuse/core";
import Loading from "vue-loading-overlay";
import "vue-loading-overlay/dist/css/index.css";
import router from "@/router";

const { accessToken, requireAuthRoute } = storeToRefs(useAuthStore());

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

watch(fetchResponse, () => {
  if (fetchResponse.value.status === 401) {
    router.push(requireAuthRoute);
  }
})

const pipelines = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.pipelines : []);
const isLoading = computed(() => isFetching.value)

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />
      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="'Runs'" main />
      <PipelinesRunsTable :items="pipelines" checkable />
    </SectionMain>
  </LayoutAuthenticated>
</template>