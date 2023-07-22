<script setup>
import { computed, ref, watch } from "vue";
import { storeToRefs } from "pinia";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { mdiReload, mdiEye, mdiPlus, mdiChevronRight, mdiChevronDown } from "@mdi/js";
import $ from 'jquery';
import RunsTable from "@/components/RunsTable.vue";
import CardBoxModal from "@/components/CardBoxModal.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import BaseIcon from "@/components/BaseIcon.vue";
import { useToast } from 'primevue/usetoast';
import Loading from "vue-loading-overlay";
import { formatDate } from '@/util';
import { i18n } from '@/i18n';

const { accessToken } = storeToRefs(useAuthStore());
const toast = useToast();
const { t } = i18n.global;

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

const isLoading = computed(() => isFetchingSubrows.value || isCreatingSubrow.value);
const isRequestError = ref(false);
const requestError = ref("");
const subRows = computed(() => fetchSubrowsResponse.value?.data ? fetchSubrowsResponse.value.data.runs : []);

watch(fetchSubrowsResponse, (value) => {
    if (value.error) {
        isRequestError.value = true;
        let header = t('global.errors.generic.header');
        let detail = value.error.message;

        if (value.status == 401) {
            header = t('global.errors.authorization.header');
            detail = t('global.errors.authorization.detail');
        }

        toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
    } else {
        slideDownSubRow();
    }
})

watch(createSubrowResponse, (value) => {
    if (value.error) {
        isRequestError.value = true;
        let header = t('global.errors.generic.header');
        let detail = value.error.message;

        if (value.status == 401) {
            header = t('global.errors.authorization.header');
            detail = t('global.errors.authorization.detail');
        }

        toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
    } else {
        fetchSubrows(null, props.parentRow.ID);
    }
})

const isCreateModalActive = ref(false);

const isRowOpen = ref(props.isRowOpened);

const emit = defineEmits(["expand-collapse-row", "create-subrow"]);

const expandOrCollapseRow = () => {
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

const onCreateSubrow = () => {
    createSubrow(null, props.parentRow.ID);
}

const slideDownSubRow = () => {
    const subRow = $("#subrow-" + props.parentRow.ID);
    subRow.removeClass("hidden-subrow");
    subRow.addClass("show-subrow");
}

const slideUpSubRow = () => {
    const subRow = $("#subrow-" + props.parentRow.ID);
    subRow.removeClass("show-subrow");
    subRow.addClass("hidden-subrow");
}

const onReloadClicked = () => {
    if (isRowOpen.value) {
        fetchSubrows(null, props.parentRow.ID)
        emit("expand-collapse-row", props.parentRow.ID, isRowOpen.value);
    }
}

</script>

<template>
    <loading v-model:active="isLoading" :is-full-page="false" />

    <tr :key="parentRow.ID">
        <td class="border-b-0 lg:w-6 before:hidden">
            <BaseIcon :path="isRowOpen ? mdiChevronDown : mdiChevronRight" @click.prevent="(e) => expandOrCollapseRow(e, parentRow.ID)" />
        </td>
        <td :data-label="$t('pages.runs.pipelineRuns.table.headers.name')">
            {{ parentRow.name }}
        </td>
        <td :data-label="$t('pages.runs.pipelineRuns.table.headers.lastRun')" class="lg:w-1 whitespace-nowrap">
            <small class="text-gray-500 dark:text-slate-400" :title="formatDate(parentRow.LastRun)">{{
                formatDate(parentRow.LastRun) }}</small>
        </td>
        <td class="before:hidden lg:w-1 whitespace-nowrap">
            <BaseButtons type="justify-start lg:justify-end" no-wrap>
                <BaseButton color="success" :icon="mdiReload" @click="onReloadClicked" />
                <BaseButton color="info" :icon="mdiEye" small :to="'/pipelines/runs/' + parentRow.ID" />
                <BaseButton color="success" :icon="mdiPlus" small @click.prevent="onRowCreateButtonClicked" />
            </BaseButtons>
        </td>
    </tr>
    <tr>
        <td class="border-b-0 lg:w-6 before:hidden" colspan="100">
            <div :id="'subrow-' + parentRow.ID" class="hidden-subrow">
                <RunsTable :pipelineID="parentRow.ID" :rows="subRows" />
            </div>
        </td>
    </tr>

    <CardBoxModal v-model="isCreateModalActive" @confirm="onCreateSubrow"
        :title="$t('pages.runs.pipelineRuns.dialog.header', { id: props.parentRow.ID })" button="success" has-cancel />
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