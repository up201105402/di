<script setup>
  import { reactive, ref, computed, watch } from "vue";
  import { useAsyncState } from "@vueuse/core";
  import { doRequest } from "@/util";
  import { useAuthStore } from "@/stores/auth";
  import { onBeforeRouteLeave, onBeforeRouteUpdate, useRouter, useRoute } from 'vue-router';
  import {
    mdiChartTimelineVariant,
    mdiPlus
  } from "@mdi/js";
  import SectionMain from "@/components/SectionMain.vue";
  import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
  import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
  import BaseButtons from "@/components/BaseButtons.vue";
  import BaseButton from "@/components/BaseButton.vue";
  import FlowChart from "@/components/FlowChart.vue";
  import CardBoxModal from '@/components/CardBoxModal.vue';
  import CreateStepDialog from '@/components/CreateStepDialog.vue';
  import Loading from "vue-loading-overlay";
  import { initialElements } from '@/flowChart.js';

  const { accessToken, requireAuthRoute } = useAuthStore();
  const router = useRouter();
  const route = useRoute();

  // FETCH PIPELINE

  const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchPipelines } = useAsyncState(
    () => {
      return doRequest({
        url: `/api/pipeline/${route.params.id}`,
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
  );

  // UPDATE PIPELINE

  const { isLoading: isUpdating, state: updateResponse, isReady: isUpdateFinished, execute: updatePipeline } = useAsyncState(
    () => {
      return doRequest({
        url: `/api/pipeline/${route.params.id}`,
        method: 'POST',
        data: {
          id: parseInt(route.params.id),
          definition: JSON.stringify(elements.value)
        },
        headers: {
          Authorization: `${accessToken}`
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

  watch(fetchResponse, () => {
    if (fetchResponse.value.status === 401) {
      router.push(requireAuthRoute);
    } else if (fetchResponse.value.status === 500) {
      router.push("/pipelines");
    }
  })

  watch(isFetchFinished, () => {
    elements.value = fetchResponse.value?.data ? JSON.parse(fetchResponse.value.data.pipeline.definition) : [];
  })

  watch(updateResponse, () => {
    if (updateResponse.value.status === 200) {
      router.push("/pipelines");
    }
  })

  const isLoading = computed(() => isFetching.value || isUpdating.value);
  const definition = computed(() => fetchResponse.value?.data ? fetchResponse.value.data.pipeline.definition : []);
  const elements = ref([]);

  const hasChanges = ref(false);
  const isCreateStepActive = ref(false);
  let count = 0;

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

  const onCreateStepClick = (e) => isCreateStepActive.value = !isCreateStepActive.value;

  const newStep = reactive({
    name: "",
    type: "",
  });

  const getNextId = () => {
    return Math.max(...elements.value.map(element => parseInt(element.id))) + 1;
  }

  const onStepCreate = (data) => {
    elements.value.push({ id: getNextId(), type: 'input', label: data.stepName, position: { x: 0, y: 0 }, class: 'light' })
    isCreateStepActive.value = false;
    hasChanges.value = true;
    count++;
  }

  const onPipelineSave = () => {
    updatePipeline();
  }

  const onPipelineCancel = () => {
    hasChanges.value = false;
    router.push('/pipelines');
  }

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <loading v-model:active="isLoading" :is-full-page="false" />

      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onCreateStepClick" />
      </SectionTitleLineWithButton>
      <FlowChart v-if="$route.params.id" v-model="elements" />
      <CardBoxModal :key="'createDialog_' + count" v-model="isCreateStepActive" :has-submit="false" :has-cancel="true"
        title="Create Step" @cancel="count++">
        <CreateStepDialog @onSubmit="onStepCreate" />
      </CardBoxModal>
      <BaseButtons style="float:right">
        <BaseButton :disabled="!hasChanges" :label="'Save'" color="success" @click="onPipelineSave" />
        <BaseButton :disabled="!hasChanges" :label="'Cancel'" color="danger" @click="onPipelineCancel" />
      </BaseButtons>
    </SectionMain>
  </LayoutAuthenticated>
</template>