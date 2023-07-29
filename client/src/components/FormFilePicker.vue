<script setup>
import { mdiUpload } from "@mdi/js";
import { computed, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { useToast } from 'primevue/usetoast';
import BaseButton from "@/components/BaseButton.vue";
import ProgressBar from 'primevue/progressbar';
import Toast from 'primevue/toast';
import { i18n } from "@/i18n";

const props = defineProps({
  modelValue: {
    type: [Object, File, Array],
    default: null,
  },
  filename: {
    type: String,
    default: null,
  },
  label: {
    type: String,
    default: null,
  },
  icon: {
    type: String,
    default: mdiUpload,
  },
  accept: {
    type: String,
    default: null,
  },
  color: {
    type: String,
    default: "info",
  },
  isRoundIcon: Boolean,
  url: {
    type: String,
    required: true
  }
});

const emit = defineEmits(["update:modelValue", "fileUpdated"]);

const { t } = i18n.global;

const { accessToken } = storeToRefs(useAuthStore());
const toast = useToast();

const root = ref(null);
const file = ref(props.modelValue);
const uploadPercent = ref(0)
const filename = ref(props.filename);

// UPLOAD FILE
const { isLoading: isUpdatingFile, state: uploadFileResponse, isReady: isFileUploadFinished, execute: uploadFile } = useAsyncState(
  (formData) => {
    return doRequest({
      url: props.url,
      method: 'POST',
      headers: {
        Authorization: `${accessToken.value}`,
        'Content-Type': 'multipart/form-data'
      },
      data: formData,
      onUploadProgress: progressEvent
    })
  },
  {},
  {
    delay: 500,
    resetOnExecute: false,
    immediate: false,
  },
)

watch(uploadFileResponse, (response) => {
  if (response.error) {
    let header = t('global.errors.generic.header');
    let detail = response.error.message;

    if (response.status == 401) {
      header = t('global.errors.authorization.header');
      detail = t('global.errors.authorization.detail');
    }

    toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
  } else {
    toast.add({ severity: 'info', summary: 'Success', detail: 'Upload Completed', life: 3000 });
    filename.value = response.data.filename;
    emit('fileUpdated', filename.value);
  }
})

const showFilename = computed(() => !props.isRoundIcon && filename.value);

const modelValueProp = computed(() => props.modelValue);

watch(modelValueProp, (value) => {
  file.value = value;

  if (!value) {
    root.value.input.value = null;
  }
});

const upload = (event) => {
  const value = event.target.files || event.dataTransfer.files;

  file.value = value[0];

  emit("update:modelValue", file.value);

  let formData = new FormData()
  formData.append('file', file.value)

  uploadPercent.value = 0;

  uploadFile(null, formData);
};

const progressEvent = progressEvent => {
  uploadPercent.value = Math.round((progressEvent.loaded * 100) / progressEvent.total);
}

const isLoading = computed(() => isUpdatingFile.value);

</script>

<template>
  <div>
    <div class="flex items-stretch justify-start relative mb-3">
      <label class="inline-flex">
        <BaseButton
          as="a"
          :class="{ 'w-12 h-12': isRoundIcon, 'rounded-r-none': showFilename }"
          :icon-size="isRoundIcon ? 24 : undefined"
          :label="isRoundIcon ? null : label"
          :icon="icon"
          :color="color"
          :rounded-full="isRoundIcon"
        />
        <input
          ref="root"
          type="file"
          class="absolute top-0 left-0 w-full h-full opacity-0 outline-none cursor-pointer -z-1"
          :accept="accept"
          @input="upload"
        />
      </label>
      <div
        v-if="showFilename"
        class="px-4 py-2 bg-gray-100 dark:bg-slate-800 border-gray-200 dark:border-slate-700 border rounded-r w-full"
      >
        <span class="text-ellipsis line-clamp-1">
          {{ filename }}
        </span>
      </div>
    </div>
    <ProgressBar :value="uploadPercent" />
    <Toast />
  </div>
</template>
