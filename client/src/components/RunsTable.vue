<script setup>
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { mdiEye, mdiPlayOutline, mdiRunFast } from "@mdi/js";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import BaseLevel from "@/components/BaseLevel.vue";
import CardBoxModal from '@/components/CardBoxModal.vue';
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
    },
    pipelineID: {
        type: Number,
        required: true,
    }
});

const isRunModalActive = ref(false);
const runIDtoExecute = ref(null);

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

// FETCH RUNS
const {
    isLoading: isFetchingRuns,
    state: fetchRunsResponse,
    isReady: isFetchRunsFinished,
    error: fetchRunsError,
    execute: fetchRuns
} = useAsyncState(
    () => {
        return doRequest({
            url: `/api/run/${props.pipelineID}`,
            method: 'GET',
            headers: {
                Authorization: `${accessToken.value}`
            },
        })
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
        let header = t('global.errors.generic.header');
        let detail = value.error.message;

        if (value.status == 401) {
            header = t('global.errors.authorization.header');
            detail = t('global.errors.authorization.detail');
        }

        toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
    } else {
        fetchRuns();
    }
})

watch(fetchRunsResponse, (value) => {
    if (value.error) {
        let header = t('global.errors.generic.header');
        let detail = value.error.message;

        if (value.status == 401) {
            header = t('global.errors.authorization.header');
            detail = t('global.errors.authorization.detail');
        }

        toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
    } else {
        fetchedRows.value = fetchRunsResponse.value?.data ? fetchRunsResponse.value.data.runs : [];
    }
})

const onSubRowButtonClicked = (e, subRowID) => {
    isRunModalActive.value = true;
    runIDtoExecute.value = subRowID;
}

const onExecuteRunConfirmed = (targetID) => {
    executeSubrow(null, targetID);
}

const isRunButtonDisabled = (run) => {
    return run.RunStatusID == 2;
}

const perPage = ref(5);
const currentPage = ref(0);
const providedRows = computed(() => props.rows);
const activeRows = ref([]);
const paginatedRows = ref([]);
const fetchedRows = ref([]);

watch(providedRows, (newValue, oldValue) => {
    activeRows.value = newValue;
    paginatedRows.value = newValue ? newValue.slice(
        perPage.value * currentPage.value,
        perPage.value * (currentPage.value + 1)
    ) : [];
})

watch(fetchedRows, (newValue, oldValue) => {
    activeRows.value = newValue;
    paginatedRows.value = newValue ? newValue.slice(
        perPage.value * currentPage.value,
        perPage.value * (currentPage.value + 1)
    ) : [];
})

watch(currentPage, (newValue, oldValue) => {
    paginatedRows.value = activeRows.value ? activeRows.value.slice(
        perPage.value * newValue,
        perPage.value * (newValue + 1)
    ) : [];
})

// const paginatedRows = computed(() => {
//     return fetchedRows.value.length ? fetchedRows.value.slice(
//         perPage.value * currentPage.value,
//         perPage.value * (currentPage.value + 1)
//     ) : props.rows.slice(
//         perPage.value * currentPage.value,
//         perPage.value * (currentPage.value + 1)
//     );
// });

const numPages = computed(() => Math.ceil(fetchedRows.value.length ? fetchedRows.value.length / perPage.value : props.rows.length / perPage.value));

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
                <th>{{ $t('pages.runs.table.headers.id') }}</th>
                <th>{{ $t('pages.runs.table.headers.status') }}</th>
                <th>{{ $t('pages.runs.table.headers.created') }}</th>
                <th>{{ $t('pages.runs.table.headers.lastRun') }}</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in paginatedRows" :key="row.ID">
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.ID }}
                </td>
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.RunStatus.Name }}
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.CreatedAt)">{{
                        formatDate(row.CreatedAt) }}</small>
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.LastRun)">{{
                        formatDate(row.LastRun) }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="success" :icon="mdiEye" small :to="`/runresults/${row.ID}`" />
                        <BaseButton color="success" :icon="mdiPlayOutline" small :disabled="isRunButtonDisabled(row)"
                            @click.prevent="(e) => onSubRowButtonClicked(e, row.ID)" />
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
            <small>{{ $t('tables.page', { page: currentPageHuman, count: numPages }) }}</small>
        </BaseLevel>
    </div>
    <CardBoxModal v-model="isRunModalActive" @confirm="onExecuteRunConfirmed" :targetId="runIDtoExecute"
        :title="$t('pages.runs.table.dialog.execute.header', { id: runIDtoExecute })" button="success" has-cancel>
        <div>{{ $t('pages.runs.table.dialog.execute.body') }}</div>
</CardBoxModal></template>