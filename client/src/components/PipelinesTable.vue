<script setup>
import { computed, ref, watch } from "vue";
import { mdiEye, mdiTrashCan } from "@mdi/js";
import CardBoxModal from "@/components/CardBoxModal.vue";
import TableCheckboxCell from "@/components/TableCheckboxCell.vue";
import BaseLevel from "@/components/BaseLevel.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import UserAvatar from "@/components/UserAvatar.vue";
import { storeToRefs } from 'pinia';
import { useAuthStore } from "@/stores/auth";
import { useAsyncState } from '@vueuse/core'
import { doRequest } from "@/util";

const props = defineProps({
    items: Array,
    checkable: Boolean,
});

const { idToken } = storeToRefs(useAuthStore());

const { isLoading, state: deleteResponse, isReady: deleteFinished, execute: deletePipeline } = useAsyncState(
    (pipelineID) => {
        if (pipelineID) {
            return doRequest({
                url: '/api/pipeline',
                method: 'DELETE',
                headers: {
                    Authorization: `${idToken.value}`,
                },
                data: {
                    ID: pipelineID
                },
            });
        }

        return {};
    },
    {},
    {
        delay: 200,
        resetOnExecute: false,
        immediate: false,
    },
)

const emit = defineEmits(["pipelineDeleted"]);

watch(deleteFinished, () => {
    if (deleteFinished.value) {
        emit("pipelineDeleted");
    }
})

const isErrorModalActive = ref(false);

const isModalDangerActive = ref(false);

const pipelineIdToDelete = ref(0);

const perPage = ref(5);

const currentPage = ref(0);

const checkedRows = ref([]);

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

const showDeleteModal = (id) => {
    isModalDangerActive.value = true;
    pipelineIdToDelete.value = id;
}

</script>

<template>
    <CardBoxModal v-model="deleteResponse.error" title="Sample modal">
        <p>Lorem ipsum dolor sit amet <b>adipiscing elit</b></p>
        <p>This is sample modal</p>
    </CardBoxModal>

    <CardBoxModal v-model="isModalDangerActive" title="Confirm Delete" :target-id="pipelineIdToDelete"
        @confirm="deletePipeline(200, pipelineIdToDelete)" button="danger" has-cancel>
        <p>This will permanently delete this pipeline.</p>
    </CardBoxModal>

    <div v-if="checkedRows.length" class="p-3 bg-gray-100/50 dark:bg-slate-800">
        <span v-for="checkedRow in checkedRows" :key="checkedRow.id"
            class="inline-block px-2 py-1 rounded-sm mr-2 text-sm bg-gray-100 dark:bg-slate-700">
            {{ checkedRow.name }}
        </span>
    </div>

    <table>
        <thead>
            <tr>
                <th v-if="checkable" />
                <th />
                <th>Name</th>
                <th>Progress</th>
                <th>Created</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <tr v-for="pipeline in itemsPaginated" :key="pipeline.id">
                <TableCheckboxCell v-if="checkable" @checked="checked($event, pipeline)" />
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
                        <BaseButton color="info" :icon="mdiEye" small :to="'/pipeline/' + pipeline.ID" />
                        <BaseButton color="danger" :icon="mdiTrashCan" small :target-id="pipeline.ID"
                            @clicked="showDeleteModal" />
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