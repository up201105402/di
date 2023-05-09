<script setup>
import { reactive, ref } from "vue";
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRouter } from 'vue-router';
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
import { initialElements } from '@/flowChart.js'

const elements = ref(initialElements);

const hasChanges = ref(false);
const isCreateStepActive = ref(false);
const router = useRouter();
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
  console.error("Save");
}

const onPipelineCancel = () => {
  hasChanges.value = false;
  router.push('/pipelines');
}

</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
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