<script setup>
    import { ref, computed } from 'vue';
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
    import { mdiRunFast, mdiPlus } from '@mdi/js';
    import SectionTitleLineWithButton from '@/components/SectionTitleLineWithButton.vue';
    import BaseButton from "@/components/BaseButton.vue";
    import RunsTable from '@/components/RunsTable.vue';
    import CardBoxModal from "@/components/CardBoxModal.vue";
    import Loading from "vue-loading-overlay";

    const { accessToken } = storeToRefs(useAuthStore());
    const route = useRoute();
    const toast = useToast();

    const isCreateModalActive = ref(false);

    const pipelineID = route.params.id;

    // FETCH RUNS

    const {
        isLoading: isFetchingRuns,
        state: fetchRunsResponse,
        isReady: isFetchRunsFinished,
        error: fetchError,
        execute: fetchRuns
    } = useAsyncState(
        () => {
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

    // CREATE RUN

    const {
        isLoading: isCreatingRun,
        state: createRunResponse,
        isReady: isCreateRunFinished,
        error: createError,
        execute: createRun
    } = useAsyncState(
        (pipelineID) => {
            return doRequest({
                url: `/api/run/${pipelineID}`,
                method: 'POST',
                headers: {
                    Authorization: `${accessToken.value}`,
                },
                data: {
                    Execute: false
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

    watch(fetchRunsResponse, (value) => {
        if (value.error) {
            toast.add({ severity: 'error', summary: 'Error', detail: value.error.message, life: 3000 });
        }
    })

    watch(createRunResponse, (value) => {
        if (value.error) {
            toast.add({ severity: 'error', summary: 'Error', detail: value.error.message, life: 3000 });
        } else {
            fetchRuns(null);
        }
    })

    const runs = computed(() => fetchRunsResponse.value?.data ? fetchRunsResponse.value.data.runs : []);
    const isLoading = computed(() => isFetchingRuns.value || isCreatingRun.value );

    const onNewRunClicked = () => {
        isCreateModalActive.value = true;
    }

    const onCreateRunConfirmed = () => {
        createRun(null, pipelineID);
    }

</script>

<template>
    <LayoutAuthenticated>
        <SectionMain>
            <loading v-model:active="isLoading" :is-full-page="false" />
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiRunFast" :title="`Pipeline ${pipelineID} Runs`" main>
                <BaseButton :icon="mdiPlus" color="success" @click="onNewRunClicked" />
            </SectionTitleLineWithButton>
            <RunsTable :rows="runs" />
        </SectionMain>
        <CardBoxModal v-model="isCreateModalActive" @confirm="onCreateRunConfirmed" :title="`Create Run for Pipeline ${pipelineID}?`" button="success" has-cancel />
        <Toast />
    </LayoutAuthenticated>
</template>