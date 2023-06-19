<script setup>
    import { computed, ref, watch } from "vue";
    import { storeToRefs } from "pinia";
    import { useAsyncState } from "@vueuse/core";
    import { doRequest } from "@/util";
    import { useAuthStore } from "@/stores/auth";
    import { mdiEye, mdiPlus, mdiChevronRight, mdiChevronDown } from "@mdi/js";
    import $ from 'jquery';
    import RunsTable from "@/components/RunsTable.vue";
    import CardBoxModal from "@/components/CardBoxModal.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import UserAvatar from "@/components/UserAvatar.vue";
    import BaseIcon from "@/components/BaseIcon.vue";
    import { useToast } from 'primevue/usetoast';
    import Loading from "vue-loading-overlay";

    import ErrorModal from '@/components/ErrorModal.vue';

    const { accessToken } = storeToRefs(useAuthStore());
    const toast = useToast();

    const props = defineProps({
        parentRow: {
            type: Object,
            required: true,
        },
        isRowOpened: {
            type: Boolean,
            required: false,
            default: false,
        },
        subrowsFetchBaseURL: {
            type: String,
            required: false
        },
    });

    // FETCH RUNS

    const {
        isLoading: isFetchingSubrows,
        state: fetchSubrowsResponse,
        isReady: isFetchSubrowsFinished,
        error: fetchError,
        execute: fetchSubrows
    } = useAsyncState(
        (rowID) => {
            return doRequest({
                url: `${props.subrowsFetchBaseURL}${rowID}`,
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

    // CREATE RUN

    const {
        isLoading: isCreatingSubrow,
        state: createSubrowResponse,
        isReady: createSubrowFinished,
        error: createError,
        execute: createSubrow
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

    const isLoading = computed(() => isFetchingSubrows.value || isCreatingSubrow.value);
    const isRequestError = ref(false);
    const requestError = ref("");
    const subRows = computed(() => fetchSubrowsResponse.value?.data ? fetchSubrowsResponse.value.data.runs : []);

    watch(fetchSubrowsResponse, (newVal) => {
        if (newVal.error) {
            isRequestError.value = true;
            requestError.value = newVal.error.message;
            toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
        } else {
            slideDownSubRow();
        }
    })

    watch(createSubrowResponse, (newVal) => {
        if (newVal.error) {
            isRequestError.value = true;
            requestError.value = newVal.error.message;
            toast.add({ severity: 'error', summary: 'Error', detail: newVal.error.message, life: 3000 });
        }
    })

    const isCreateModalActive = ref(false);

    const isRowOpen = ref(props.isRowOpened);

    const emit = defineEmits(["expand-collapse-row", "create-subrow"]);

    const expandOrCollapseRow = (e) => {
        isRowOpen.value = !isRowOpen.value;

        if (isRowOpen.value) {
            const fetch = fetchSubrows(null, props.parentRow.ID)
        } else {
            slideUpSubRow();
        }

        emit("expand-collapse-row", props.parentRow.ID, isRowOpen.value)
    }

    const onRowCreateButtonClicked = (e) => {
        isCreateModalActive.value = true;
        emit("create-subrow", props.parentRow.ID);
    }

    const onSubRowButtonClicked = (e, subRowID) => {
        executeSubrow(null, subRowID);
    }

    const onCreateSubrow = (e) => {
        createSubrow(null, props.parentRow.ID);
        fetchSubrows(null, props.parentRow.ID);
    }

    const acknowledgeError = (e) => {
        isRequestError.value = null
    }

    const slideDownSubRow = (subRowID) => {
        const subRow = $("#subrow-" + props.parentRow.ID);
        subRow.removeClass("hidden-subrow");
        subRow.addClass("show-subrow");
    }

    const slideUpSubRow = (subRowID) => {
        const subRow = $("#subrow-" + props.parentRow.ID);
        subRow.removeClass("show-subrow");
        subRow.addClass("hidden-subrow");
    }

</script>

<template>

    <loading v-model:active="isLoading" :is-full-page="false" />

    <tr :key="parentRow.ID">
        <td class="border-b-0 lg:w-6 before:hidden">
            <BaseIcon :path="isRowOpen ? mdiChevronDown : mdiChevronRight"
                @click.prevent="(e) => expandOrCollapseRow(e, parentRow.ID)" />
        </td>
        <td class="border-b-0 lg:w-6 before:hidden">
            <UserAvatar :username="parentRow.name" class="w-24 h-24 mx-auto lg:w-6 lg:h-6" />
        </td>
        <td data-label="Name">
            {{ parentRow.name }}
        </td>
        <td data-label="Progress" class="lg:w-32">
            <progress class="flex w-2/5 self-center lg:w-full" max="100" :value="parentRow.progress">
                {{ parentRow.progress }}
            </progress>
        </td>
        <td data-label="Created" class="lg:w-1 whitespace-nowrap">
            <small class="text-gray-500 dark:text-slate-400" :title="parentRow.CreatedAt">{{ parentRow.CreatedAt }}</small>
        </td>
        <td class="before:hidden lg:w-1 whitespace-nowrap">
            <BaseButtons type="justify-start lg:justify-end" no-wrap>
                <BaseButton color="info" :icon="mdiEye" small :to="'/run/' + parentRow.ID" />
                <BaseButton color="success" :icon="mdiPlus" small @click.prevent="onRowCreateButtonClicked" />
            </BaseButtons>
        </td>
    </tr>
    <tr>
        <td class="border-b-0 lg:w-6 before:hidden" colspan="100">
            <div :id="'subrow-' + parentRow.ID" class="hidden-subrow">
                <RunsTable :rows="subRows" />
            </div>
        </td>
    </tr>

    <CardBoxModal v-model="isCreateModalActive" @confirm="onCreateSubrow" :title="`Create Run for Pipeline ${props.parentRow.ID}?`" button="success" has-cancel />

</template>

<style>
    .hidden-subrow {
        max-height: 0;
        transition: max-height 0.15s ease-out;
        overflow: hidden;
    }

    .show-subrow {
        max-height: 500px;
        transition: max-height 0.25s ease-in;
    }
</style>