<script setup>
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { mdiPlayOutline, mdiRunFast } from "@mdi/js";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import BaseLevel from "@/components/BaseLevel.vue";
import { useToast } from 'primevue/usetoast';
import { ref, computed, watch } from "vue";
import { formatDate } from '@/util';
import Loading from "vue-loading-overlay";

const { accessToken } = storeToRefs(useAuthStore());
const toast = useToast();

const props = defineProps({
    rows: {
        type: Array,
        required: true,
    }
});

const disabledRuns = ref([])

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
        toast.add({ severity: 'error', summary: 'Error', detail: value.error.message, life: 3000 });
    }
})

const onSubRowButtonClicked = (e, subRowID) => {
    disabledRuns.value.push(subRowID);
    executeSubrow(null, subRowID);
}

const isRunButtonDisabled = (run) => {
    return run.Status.ID != 1 || disabledRuns.value.includes(run.ID);
}

const perPage = ref(5);

const currentPage = ref(0);

const paginatedRows = computed(() => {
    return props.rows ? props.rows.slice(
        perPage.value * currentPage.value,
        perPage.value * (currentPage.value + 1)
    ) : [];
});

const numPages = computed(() => Math.ceil(props.rows.length / perPage.value));

const currentPageHuman = computed(() => currentPage.value + 1);

const pagesList = computed(() => {
    const pagesList = [];

    for (let i = 0; i < numPages.value; i++) {
        pagesList.push(i);
    }

    return pagesList;
});

const isLoading = computed(() => isExecutingSubrow.value);

</script>

<template>
    <loading v-model:active="isLoading" :is-full-page="false" />
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>Status</th>
                <th>Created</th>
                <th>Last Run</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in paginatedRows" :key="row.ID">
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.ID }}
                </td>
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.Status.Name }}
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.CreatedAt)">{{ formatDate(row.CreatedAt) }}</small>
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.LastRun)">{{ formatDate(row.LastRun) }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="success" :icon="mdiPlayOutline" small :disabled="isRunButtonDisabled(row)" @click.prevent="(e) => onSubRowButtonClicked(e, row.ID)" />
                    </BaseButtons>
                </td>
            </tr>
        </tbody>
    </table>
    <div class="p-3 lg:px-6 border-t border-gray-100 dark:border-slate-800">
        <BaseLevel>
            <BaseButtons>
                <BaseButton v-for="page in pagesList" :key="page" :active="page === currentPage" :label="page + 1"
                    :color="page === currentPage ? 'lightDark' : 'whiteDark'" small @click="currentPage = page" />
            </BaseButtons>
            <small>Page {{ currentPageHuman }} of {{ numPages }}</small>
        </BaseLevel>
    </div>
</template>