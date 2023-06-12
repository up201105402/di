<script setup>
    import { computed, ref } from "vue";
    import { useAsyncState } from "@vueuse/core";
    import { doRequest } from "@/util";
    import { useAuthStore } from "@/stores/auth";
    import { mdiRefresh, mdiPlus, mdiChevronRight, mdiChevronDown, mdiRunFast } from "@mdi/js";
    import CardBoxModal from "@/components/CardBoxModal.vue";
    import TableCheckboxCell from "@/components/TableCheckboxCell.vue";
    import BaseLevel from "@/components/BaseLevel.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import UserAvatar from "@/components/UserAvatar.vue";
    import BaseIcon from "@/components/BaseIcon.vue";
    import Loading from "vue-loading-overlay";

    const { accessToken, requireAuthRoute } = useAuthStore();

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
        execute: fetchSubrows
    } = useAsyncState(
        (rowID) => {
            return doRequest({
                url: `${props.subrowsFetchBaseURL}${rowID}`,
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

    // CREATE RUN

    const {
        isLoading: isCreatingSubrow,
        state: createSubrowResponse,
        isReady: createSubrowFinished,
        execute: createSubrow
    } = useAsyncState(
        (pipelineID) => {
            return doRequest({
                url: `/api/run/${pipelineID}`,
                method: 'POST',
                headers: {
                    Authorization: `${accessToken}`,
                },
                data: {
                    Execute: false
                }
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
    const subRows = computed(() => fetchSubrowsResponse.value?.data ? fetchSubrowsResponse.value.data.runs : []);

    const isCreateModalActive = ref(false);

    const isRowOpen = ref(props.isRowOpened);

    const emit = defineEmits(["expand-collapse-row", "create-subrow"]);

    const expandOrCollapseRow = (e) => {
        isRowOpen.value = !isRowOpen.value;

        if (isRowOpen.value) {
            fetchSubrows(null, props.parentRow.ID)
        }

        emit("expand-collapse-row", props.parentRow.ID, isRowOpen.value)
    }

    const onSubRowCreateButtonClicked = (e) => {
        isCreateModalActive.value = true;
        emit("create-subrow", props.parentRow.ID);
    }

    const onCreateSubrow = (e) => {
        createSubrow(null, props.parentRow.ID);
        fetchSubrows(null, props.parentRow.ID);
    }

</script>

<template>

    <loading v-model:active="isLoading" :is-full-page="false" />

    <tr :key="parentRow.ID">
        <td class="border-b-0 lg:w-6 before:hidden">
            <BaseIcon :path="isRowOpen ? mdiChevronDown : mdiChevronRight" @click.prevent="(e) => expandOrCollapseRow(e, parentRow.ID)" />
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
            <small class="text-gray-500 dark:text-slate-400" :title="parentRow.CreatedAt">{{ parentRow.CreatedAt
                }}</small>
        </td>
        <td class="before:hidden lg:w-1 whitespace-nowrap">
            <BaseButtons type="justify-start lg:justify-end" no-wrap>
                <BaseButton color="success" :icon="mdiPlus" small @click.prevent="onSubRowCreateButtonClicked" />
                <BaseButton color="success" :icon="mdiRefresh" small />
            </BaseButtons>
        </td>
    </tr>
    <tr v-if="isRowOpen" v-for="subRow in subRows" :key="subRow.ID">
        <td class="border-b-0 lg:w-6 before:hidden">
            <BaseIcon :path="mdiRunFast" />
        </td>
        <td class="border-b-0 lg:w-6 before:hidden">
            {{ subRow.Status.Name }}
        </td>
        <td data-label="Name">
            {{ subRow.ID }}
        </td>
        <td data-label="Progress" class="lg:w-32">
            <progress class="flex w-2/5 self-center lg:w-full" max="100" :value="parentRow.progress">
                {{ subRow.progress }}
            </progress>
        </td>
        <td data-label="Created" class="lg:w-1 whitespace-nowrap">
            <small class="text-gray-500 dark:text-slate-400" :title="parentRow.CreatedAt">{{ parentRow.CreatedAt}}</small>
        </td>
        <td class="before:hidden lg:w-1 whitespace-nowrap">
            <BaseButtons type="justify-start lg:justify-end" no-wrap>
                <BaseButton color="success" :icon="mdiPlus" small
                    @click.prevent="(e) => onSubRowCreateButtonClicked(e, parentRow.ID)" />
                <BaseButton color="success" :icon="mdiRefresh" small />
            </BaseButtons>
        </td>
    </tr>

    <CardBoxModal v-model="isCreateModalActive" @confirm="onCreateSubrow"
        :title="`Create Run for Pipeline ${props.parentRow.ID}?`" button="success" has-cancel />

</template>