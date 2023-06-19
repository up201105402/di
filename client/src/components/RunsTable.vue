<script setup>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { mdiPlayOutline, mdiRunFast } from "@mdi/js";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import BaseIcon from "@/components/BaseIcon.vue";
import Loading from "vue-loading-overlay";
import { useToast } from 'primevue/usetoast';
import { watch } from "vue";

const { accessToken } = storeToRefs(useAuthStore());
const toast = useToast();

const props = defineProps({
    rows: {
        type: Array,
        required: true,
    }
});

// EXECUTE RUN

const {
    isLoading: isExecutingSubrow,
    state: executeSubrowResponse,
    isReady: executeSubrowFinished,
    error: executeError,
    execute: executeSubrow
} = useAsyncState(
    (subRowID) => {
        return doRequest({
            url: `/api/run/execute/${subRowID}`,
            method: 'POST',
            headers: {
                Authorization: `${accessToken.value}`,
            },
        });
    },
    {},
    {
        delay: 500,
        resetOnExecute: false,
        immediate: false,
    },
);

watch(executeSubrowResponse, (value) => {
    if (value.error) {
        toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
    }
})

const isRequestError = ref(false);
const requestError = ref("");

const onSubRowButtonClicked = (e, subRowID) => {
    executeSubrow(null, subRowID);
}

const acknowledgeError = (e) => {
    isRequestError.value = null
}

</script>

<template>
    <table>
        <thead>
            <tr>
                <th />
                <th>Name</th>
                <th>Status</th>
                <th />
                <th>Last Run</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in rows" :key="row.ID">
                <td class="border-b-0 text-center lg:w-6 before:hidden">
                    {{ row.ID }}
                </td>
                <td class="border-b-0 text-center lg:w-6 before:hidden">
                    <BaseIcon :path="mdiRunFast" />
                </td>
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.Status.Name }}
                </td>
                <td data-label="Progress" class="lg:w-32">
                    <progress class="flex w-2/5 self-center lg:w-full" max="100" :value="row.progress">
                        {{ row.progress }}
                    </progress>
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="row.CreatedAt">{{ row.CreatedAt }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="success" :icon="mdiPlayOutline" small @click.prevent="(e) => onSubRowButtonClicked(e, row.ID)" />
                    </BaseButtons>
                </td>
            </tr>
        </tbody>
    </table>
</template>