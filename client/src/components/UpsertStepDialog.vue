<script setup>
  import { ref, computed } from 'vue';
  import stepTypes from '@/pipelines/steps'
  import CheckoutRepositoryForm from '@/pipelines/steps/components/CheckoutRepositoryForm.vue';
  import LoadTrainingDatasetForm from '@/pipelines/steps/components/LoadTrainingDatasetForm.vue';
  import TrainModelForm from '@/pipelines/steps/components/TrainModelForm.vue';
  import { watch } from 'vue';

  const props = defineProps({
    nodeId: {
      type: String,
      required: false,
      default: null
    },
    nodeData: {
      type: Object,
      required: false,
      default: null,
    },
  });

  const data = computed({
    get: () => props.nodeData,
    set: (value) => { },
  });

  const emit = defineEmits(["onSubmit"]);

  const onSubmit = (e) => {
    emit("onSubmit", { id: props.nodeId, data: e.data })
  }

</script>

<template>
  <Vueform v-model="data" sync :endpoint="false" @submit="onSubmit">
    <template #empty>
      <FormSteps>
        <FormStep name="step1" label="Select Type" :elements="['stepName', 'stepType']" />
        <FormStep name="step2" label="Configure" :elements="Object.values(stepTypes).map(stepType => stepType.name)"
          :labels="{ next: 'Submit' }" />
      </FormSteps>

      <FormElements>
        <TextElement name="stepName" label="Name" :rules="['required']" />
        <SelectElement name="stepType" :search="true" :native="false" label="Type" input-type="search"
          autocomplete="off" :items="Object.values(stepTypes).map(stepType => ({
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