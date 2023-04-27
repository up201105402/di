<script setup>
import { ref, reactive, computed } from "vue";
import {
  mdiChartTimelineVariant,
  mdiPlus
} from "@mdi/js";
import { storeToRefs } from 'pinia';
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
import { useAsyncState } from '@vueuse/core'

const isNewPipelineModalActive = ref(false);
const { idToken } = storeToRefs(useAuthStore());

const { isLoading, state: fetchResponse, isReady, execute: fetchPipelines } = useAsyncState(
  () => {
    return doRequest({
      url: '/api/pipeline',
      method: 'GET',
      headers: {
        Authorization: `${idToken.value}`
      },
    })
  },
  {},
  {
    delay: 200,
    resetOnExecute: false,
  },
)

const onNewPipelineClicked = (e) => isNewPipelineModalActive.value = true;

const createNewPipeline = (e) => {
  doRequest({
    url: '/api/pipeline',
    method: 'POST',
    headers: {
      Authorization: `${idToken.value}`,
    },
    data: {
      name: form.name
    },
  });

  fetchPipelines(200);
}

const pipelines = computed(() => fetchResponse?.data?.pipelines)

const form = reactive({
  name: "",
});

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
        <BaseButton :icon="mdiPlus" @click="onNewPipelineClicked" />
      </SectionTitleLineWithButton>

      <PipelinesTable :items="fetchResponse?.data ? fetchResponse.data.pipelines : []" @pipelineDeleted="fetchPipelines(200)" checkable />

    </SectionMain>

    <CardBoxModal v-model="isNewPipelineModalActive" @confirm="createNewPipeline" title="Create Pipeline" button="success"
      has-cancel>
      <FormField label="Name" help="Please enter the pipeline name">
        <FormControl v-model="form.name" name="name" autocomplete="name" placeholder="Name" :focus="isNewPipelineModalActive" />
      </FormField>
    </CardBoxModal>
  </LayoutAuthenticated>
</template>