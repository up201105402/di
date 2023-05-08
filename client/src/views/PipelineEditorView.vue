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

  const mainStore = useMainStore();

  const clientBarItems = computed(() => mainStore.clients.slice(0, 4));

  const isCreateStepActive = ref(false);

  const onCreateStepClick = (e) => isCreateStepActive.value = !isCreateStepActive.value;

  const stepTypes = [
    { id: 1, label: "Checkout GitHub" },
    { id: 2, label: "Training Dataset" },
  ];

  const newStep = reactive({
    name: "",
    type: "",
  });

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
        <BaseButton :icon="mdiPlus" color="success" @click="onCreateStepClick" />
      </SectionTitleLineWithButton>
      <FlowChart v-if="$route.params.id" />
      <CardBoxModal v-model="isCreateStepActive" :has-submit="false" title="Create Step">
        <CreateStepDialog />
      </CardBoxModal>
    </SectionMain>
  </LayoutAuthenticated>
</template>