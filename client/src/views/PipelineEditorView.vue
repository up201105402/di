<script setup>
  import { computed, reactive, ref } from "vue";
  import { useMainStore } from "@/stores/main";
  import {
    mdiChartTimelineVariant,
    mdiPlus
  } from "@mdi/js";
  import SectionMain from "@/components/SectionMain.vue";
  import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
  import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
  import BaseButton from "@/components/BaseButton.vue";
  import FlowChart from "@/components/FlowChart.vue";
  import CardBoxModal from '@/components/CardBoxModal.vue';
  import CreateStepDialog from '@/components/CreateStepDialog.vue';
  import { initialElements } from '@/flowChart.js'

  const elements = ref(initialElements);

  const isCreateStepActive = ref(false);

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
  }

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onCreateStepClick" />
      </SectionTitleLineWithButton>
      <FlowChart v-if="$route.params.id" v-model="elements" />
      <CardBoxModal v-model="isCreateStepActive" :has-submit="false" title="Create Step">
        <CreateStepDialog @onSubmit="onStepCreate" />
      </CardBoxModal>
    </SectionMain>
  </LayoutAuthenticated>
</template>