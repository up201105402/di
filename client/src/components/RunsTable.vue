<script setup>
    import { computed, ref } from "vue";
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

    const { accessToken, requireAuthRoute } = useAuthStore();

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
                    Authorization: `${accessToken}`
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

    const {
        isLoading: isCreatingPipelineRun,
        state: createPipelineRunResponse,
        isReady: createPipelineRunFinished,
        execute: createPipelineRun
    } = useAsyncState(
        (pipelineID) => {
            return doRequest({
                url: `/api/run/${pipelineID}`,
                method: 'POST',
                headers: {
                    Authorization: `${accessToken}`,
                },
            });
        },
        {},
        {
            delay: 500,
            resetOnExecute: false,
            immediate: false,
        },
    )


    // EMITS

    const emit = defineEmits(["deleteButtonClicked"]);

    const deleteButtonClicked = (id) => {
        emit("deleteButtonClicked", id);
    }

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

    const expandRow = (e, pipelineID) => {
        if (isPipelineOpen(pipelineID)) {
            const index = openedPipelines.value.findIndex(elem => elem == pipelineID);
            openedPipelines.value.splice(index, 1);
        } else {
            openedPipelines.value.push(pipelineID);
            fetchPipelineRuns(500, pipelineID);
        }
    }

    const isPipelineOpen = (id) => {
        return openedPipelines.value.find(elem => elem == id) != null;
    }

    const isCreateModalActive = ref(false);
    const createPipelineID = ref(null);

    const onNewPipelineRunClicked = (e, pipelineID) => {
        isCreateModalActive.value = true;
        createPipelineID.value = pipelineID;
    }

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
            <tr v-for="pipeline in itemsPaginated" :key="pipeline.id">
                <td class="border-b-0 lg:w-6 before:hidden">
                    <BaseIcon :path="isPipelineOpen(pipeline.ID) ? mdiChevronDown : mdiChevronRight"
                        @click.prevent="(e) => expandRow(e, pipeline.ID)" />
                </td>
                <td class="border-b-0 lg:w-6 before:hidden">
                    <UserAvatar :username="pipeline.name" class="w-24 h-24 mx-auto lg:w-6 lg:h-6" />
                </td>
                <td data-label="Name">
                    {{ pipeline.name }}
                </td>
                <td data-label="Progress" class="lg:w-32">
                    <progress class="flex w-2/5 self-center lg:w-full" max="100" :value="pipeline.progress">
                        {{ pipeline.progress }}
                    </progress>
                </td>
                <td data-label="Created" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="pipeline.CreatedAt">{{ pipeline.CreatedAt
                        }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="success" :icon="mdiPlus" small
                            @click.prevent="(e) => onNewPipelineRunClicked(e, pipeline.ID)" />
                        <BaseButton color="success" :icon="mdiRefresh" small />
                    </BaseButtons>
                </td>
                <td v-for="run in runsByPipeline[pipeline.ID]" :key="run.ID" />
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

    <CardBoxModal v-model="isCreateModalActive" @confirm="createPipelineRun(500, createPipelineID)"
        :title="'Create Run for Pipeline ' + createPipelineID" button="success" has-cancel />
</template>