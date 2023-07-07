<script setup>
  import { computed, onMounted } from 'vue';
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
    formSchema: {
      type: Object,
      required: false
    },
    formkitData: {
      type: Object,
      required: false
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

  let stepCategory = props.stepCategory != '' ? props.stepCategory : 'default';
  // let { formSchema, formkitData } = useSteps[stepCategory](props.nodeData, onSubmit);

</script>

<template>
  <FormKitSchema v-if="formSchema" :schema="formSchema" :data="formkitData" />
</template>

<style>
  @import "@/css/formkit/multistep-form.css";
</style>