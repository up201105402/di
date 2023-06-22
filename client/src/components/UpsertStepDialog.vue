<script setup>
  import { ref, computed, reactive, markRaw, onMounted } from 'vue';
  import stepTypes from '@/pipelines/steps'
  import { watch } from 'vue';
  import useSteps from '@/pipelines/steps'
  import $ from 'jquery';

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

  onMounted(() => {
    $("#cancel-create-step-button").off('click').on('click', function (e) {
      emit("onCancel");
    });
  });

  const emit = defineEmits(["onSubmit", "onCancel"]);

  const onSubmit = (formData) => {
    emit("onSubmit", { id: props.nodeId, data: formData });
  }

  let { formSchema, formkitData, steps, visitedSteps, activeStep, nodeTypesOptions, activeNodeType, setStep, stepPlugin } = useSteps(props.nodeData, onSubmit);

</script>

<template>
  <FormKitSchema :schema="formSchema" :data="formkitData" />
</template>

<style>
  @import "@/css/formkit/multistep-form.css";
</style>