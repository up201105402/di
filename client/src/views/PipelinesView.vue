<script setup>
import { ref, reactive } from "vue";
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

const onNewPipelineClicked = (e) => isNewPipelineModalActive.value = true;

const createNewPipeline = (e) => {
  doRequest({
    url: '/api/pipeline',
    method: 'POST',
    headers: {
      Authorization: `Bearer ${idToken}`,
    },
    data: {
      name: form.name
    },
  });
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

      <PipelinesTable checkable />

    </SectionMain>

    <CardBoxModal v-model="isNewPipelineModalActive" @confirm="createNewPipeline" title="Please confirm" button="success"
      has-cancel>
      <FormField label="Name" help="Please enter the pipeline name">
        <FormControl v-model="form.name" name="name" autocomplete="name" placeholder="Name" />
      </FormField>
    </CardBoxModal>
  </LayoutAuthenticated>
</template>