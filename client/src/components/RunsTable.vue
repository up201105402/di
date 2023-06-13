<script setup>
    import { computed, ref } from "vue";
    import { storeToRefs } from "pinia";
    import { useAsyncState } from "@vueuse/core";
    import { doRequest } from "@/util";
    import { useAuthStore } from "@/stores/auth";
    import { mdiRefresh, mdiPlus, mdiChevronRight, mdiChevronDown } from "@mdi/js";
    import CardBoxModal from "@/components/CardBoxModal.vue";
    import TableCheckboxCell from "@/components/TableCheckboxCell.vue";
    import BaseLevel from "@/components/BaseLevel.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import UserAvatar from "@/components/UserAvatar.vue";
    import BaseIcon from "@/components/BaseIcon.vue";
    import RunsTableRow from "@/components/RunsTableRow.vue";

    const { accessToken, requireAuthRoute } = storeToRefs(useAuthStore());

    const props = defineProps({
        items: Array,
        checkable: Boolean,
    });

    // FETCH RUNS

    const {
        isLoading: isFetchingPipelineRuns,
        state: fetchPipelineRunsResponse,
        isReady: isFetchPipelineRunsFinished,
        execute: fetchPipelineRuns
    } = useAsyncState(
        (pipelineID) => {
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
            immediate: false,
        },
    );

    // ITEMS PROCESSING

    const perPage = ref(5);

    const currentPage = ref(0);

    const checkedRows = ref([]);

    const openedPipelines = ref([]);

    const runsByPipeline = ref({});

    const itemsPaginated = computed(() => {
        return props.items ? props.items.slice(
            perPage.value * currentPage.value,
            perPage.value * (currentPage.value + 1)
        ) : [];
    });

    const numPages = computed(() => Math.ceil(props.items?.length / perPage.value));

    const currentPageHuman = computed(() => currentPage.value + 1);

    const pagesList = computed(() => {
        const pagesList = [];

        for (let i = 0; i < numPages.value; i++) {
            pagesList.push(i);
        }

        return pagesList;
    });

    const remove = (arr, cb) => {
        const newArr = [];

        arr.forEach((item) => {
            if (!cb(item)) {
                newArr.push(item);
            }
        });

        return newArr;
    };

    const checked = (isChecked, pipeline) => {
        if (isChecked) {
            checkedRows.value.push(pipeline);
        } else {
            checkedRows.value = remove(
                checkedRows.value,
                (row) => row.id === pipeline.id
            );
        }
    };

    const expandOrCollapseRow = (pipelineID, isOpen) => {
        if (!isOpen) {
            const index = openedPipelines.value.findIndex(elem => elem == pipelineID);
            openedPipelines.value.splice(index, 1);
        } else {
            openedPipelines.value.push(pipelineID);
            // fetchPipelineRuns(500, pipelineID);
        }
    }

    const isPipelineOpen = (id) => {
        return openedPipelines.value.find(elem => elem == id) != null;
    }

    const isCreateModalActive = ref(false);
    const createPipelineID = ref(null);

</script>

<template>
    <div v-if="checkedRows.length" class="p-3 bg-gray-100/50 dark:bg-slate-800">
        <span v-for="checkedRow in checkedRows" :key="checkedRow.id"
            class="inline-block px-2 py-1 rounded-sm mr-2 text-sm bg-gray-100 dark:bg-slate-700">
            {{ checkedRow.name }}
        </span>
    </div>

    <table>
        <thead>
            <tr>
                <th />
                <th>Name</th>
                <th>Status</th>
                <th>Last Run</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <RunsTableRow v-for="pipeline in itemsPaginated" 
                                :parentRow="pipeline" 
                                subrowsFetchBaseURL="/api/run/" 
                                @expand-collapse-row="expandOrCollapseRow" />
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