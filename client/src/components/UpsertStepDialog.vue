<script setup>
  import { ref, computed, reactive, markRaw } from 'vue';
  import stepTypes from '@/pipelines/steps'
  import { watch } from 'vue';
  import useSteps from '@/pipelines/steps'

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

  const onSubmit = (formData) => {
    emit("onSubmit", { id: props.nodeId, data: formData });
  }

  let { formSchema, formkitData, steps, visitedSteps, activeStep, nodeTypesOptions, activeNodeType, setStep, stepPlugin } = useSteps(props.nodeData, onSubmit);

</script>

<template>
  <!-- <Vueform v-model="data" sync :endpoint="false" @submit="onSubmit">
    <template #empty>
      <FormSteps>
        <FormStep name="step1" label="Select Type" :elements="['stepName', 'stepType']" />
        <FormStep name="step2" label="Configure" :elements="Object.values(stepTypes).map(stepType => stepType.name)" :labels="{ next: 'Submit' }" />
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
  </Vueform> -->
  <FormKitSchema :schema="formSchema" :data="formkitData" />
</template>

<style>
  @import "https://cdn.formk.it/web-assets/multistep-form.css";
</style>