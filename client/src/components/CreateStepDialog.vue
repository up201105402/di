<script setup>
import { ref } from 'vue';
import stepTypes from '@/pipelines/steps'
import CheckoutRepositoryForm from '@/pipelines/steps/components/CheckoutRepositoryForm.vue';
import LoadTrainingDatasetForm from '@/pipelines/steps/components/LoadTrainingDatasetForm.vue';
import TrainModelForm from '@/pipelines/steps/components/TrainModelForm.vue';
import { watch } from 'vue';

const data = ref(null);

const emit = defineEmits(["onSubmit"]);

const onCreate = (e) => {
  emit("onSubmit", e.data)
}

</script>

<template>
  <Vueform v-model="data" sync :endpoint="false" @submit="onCreate">
    <template #empty>
      <FormSteps>
        <FormStep name="step1" label="Select Type" :elements="['stepName', 'stepType']" />
        <FormStep name="step2" label="Configure" :elements="Object.values(stepTypes).map(stepType => stepType.name)"
          :labels="{ next: 'Add' }" />
      </FormSteps>

      <FormElements>
        <TextElement name="stepName" label="Name" :rules="['required']" />
        <SelectElement name="stepType" :search="true" :native="false" label="Type" input-type="search" autocomplete="off"
          :items="Object.values(stepTypes).map(stepType => ({
            'value': stepType.name,
            'label': stepType.label,
          }))" :rules="['required']" />
        <CheckoutRepositoryForm />
        <LoadTrainingDatasetForm />
        <TrainModelForm />
      </FormElements>

      <FormStepsControls />
    </template>
  </Vueform>
</template>