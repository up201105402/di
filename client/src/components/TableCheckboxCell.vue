<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  type: {
    type: String,
    default: "td",
  },
  isChecked: {
    type: Boolean,
    default: false,
  },
  isDisabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["checked"]);

const checked = ref(props.isChecked);

watch(checked, (newVal) => {
  emit("checked", newVal);
});
</script>

<template>
  <component :is="type" class="lg:w-1">
    <label class="checkbox">
      <input type="checkbox" :checked="isChecked" isDisabled />
      <span class="check" />
    </label>
  </component>
</template>
