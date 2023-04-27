<script setup>
  import { ref, reactive, computed } from "vue";
  import { Promised, usePromise } from 'vue-promised';
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

  const isNewPipelineModalActive = ref(false);
  const { idToken } = storeToRefs(useAuthStore());

  const fetchPipelines = () => doRequest({
    url: '/api/pipeline',
    method: 'GET',
    headers: {
      Authorization: `${idToken.value}`,
      "Content-Type": 'application/json',
    },
  });

  // const pipelines = computed(() => fetchPipelines())
  var pipelines = fetchPipelines();

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

    pipelines = fetchPipelines();
  }

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

      <Promised :promise="pipelines">
        <!-- Use the "pending" slot to display a loading message -->
        <template v-slot:pending="previousData">
          <PipelinesTable :items="data.data.pipelines" checkable />
        </template>
        <!-- The default scoped slot will be used as the result -->
        <template v-slot="data">
          <PipelinesTable :items="data.data.pipelines" checkable />
        </template>
        <!-- The "rejected" scoped slot will be used if there is an error -->
        <template v-slot:rejected="error">
          <p>Error: {{ error.message }}</p>
        </template>
      </Promised>

    </SectionMain>

    <CardBoxModal v-model="isNewPipelineModalActive" @confirm="createNewPipeline" title="Create Pipeline"
      button="success" has-cancel>
      <FormField label="Name" help="Please enter the pipeline name">
        <FormControl v-model="form.name" name="name" autocomplete="name" placeholder="Name" />
      </FormField>
    </CardBoxModal>
  </LayoutAuthenticated>
</template>