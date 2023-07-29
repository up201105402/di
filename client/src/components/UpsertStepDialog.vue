<script setup>
  import { computed, onMounted, markRaw } from 'vue';
  import { library } from '@/pipelines/steps';
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
  });

  onMounted(() => {
    $("#cancel-create-step-button").off('click').on('click', function (e) {
      emit("onCancel");
    });
  });

  const emit = defineEmits(["onSubmit", "onCancel"]);

  const onSubmit = (formData) => {
    emit("onSubmit", { id: props.nodeId, group: props.group, type: props.type, data: formData });
  }

</script>

<template>
  <FormKitSchema v-if="formSchema" :schema="formSchema" :data="nodeData" :library="library" />
</template>

<style>
  .formkit-step-schema .ql-editor {
    counter-reset: line;
    padding-left: 0;
  }

  .formkit-step-schema .ql-editor p:before {
    counter-increment: line;
    content: counter(line);
    display: inline-block;
    border-right: 1px solid #ddd;
    padding: 0 .5em;
    margin-right: .5em;
    color: #888
  }

</style>