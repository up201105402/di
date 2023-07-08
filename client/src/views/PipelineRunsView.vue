<script setup>
import { computed } from 'vue';
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { useRoute } from 'vue-router';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import { watch } from "vue";
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue';
import SectionMain from '@/components/SectionMain.vue';
import { mdiRunFast } from '@mdi/js';
import SectionTitleLineWithButton from '@/components/SectionTitleLineWithButton.vue';
import RunsTable from '@/components/RunsTable.vue';
import Loading from "vue-loading-overlay";


const { accessToken } = storeToRefs(useAuthStore());
const route = useRoute();
const toast = useToast();

const props = defineProps({
    rows: {
        type: Array,
        required: true,
    }
});

const pipelineID = route.params.id;

// FETCH RUNS

const {
    isLoading: isFetchingRuns,
    state: fetchRunsResponse,
    isReady: isFetchRunsFinished,
    error: fetchError,
    execute: fetchRuns
} = useAsyncState(
    (rowID) => {
        return doRequest({
            url: `/api/run/${pipelineID}`,
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
    },
);

watch(fetchRunsResponse, (value) => {
    if (value.error) {
        toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
    }
})

const runs = computed(() => fetchRunsResponse.value?.data ? fetchRunsResponse.value.data.runs : []);
const isLoading = computed(() => isFetchingRuns.value);

</script>

<template>
    <LayoutAuthenticated>
        <SectionMain>
            <loading v-model:active="isLoading" :is-full-page="false" />
            <SectionTitleLineWithButton :icon="mdiRunFast" :title="`Pipeline ${pipelineID} Runs`" main />
            <RunsTable :rows="runs" />
        </SectionMain>
        <Toast />
    </LayoutAuthenticated>
</template>