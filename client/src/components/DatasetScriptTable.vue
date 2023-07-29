<script setup>
import { computed, ref } from "vue";
import { mdiPencilOutline, mdiTrashCan } from "@mdi/js";
import BaseLevel from "@/components/BaseLevel.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import { formatDate } from '@/util';

const props = defineProps({
    items: Array,
});

// EMITS

const emit = defineEmits(["deleteButtonClicked"]);

const deleteButtonClicked = (id) => {
    emit("deleteButtonClicked", id);
}

// ITEMS PROCESSING

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

</script>

<template>
    <table>
        <thead>
            <tr>
                <th>{{ $t('pages.datasets.table.headers.id') }}</th>
                <th>{{ $t('pages.datasets.table.headers.name') }}</th>
                <th>{{ $t('pages.datasets.table.headers.modified') }}</th>
                <th>{{ $t('pages.datasets.table.headers.created') }}</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <tr v-for="datasetScript in itemsPaginated" :key="datasetScript.id">
                <td :data-label="$t('pages.datasets.table.headers.name')">
                    {{ datasetScript.ID }}
                </td>
                <td :data-label="$t('pages.datasets.table.headers.name')">
                    {{ datasetScript.name }}
                </td>
                <td :data-label="$t('pages.datasets.table.headers.modified')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(datasetScript.UpdatedAt)">{{ formatDate(datasetScript.UpdatedAt) }}</small>
                </td>
                <td :data-label="$t('pages.datasets.table.headers.created')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(datasetScript.CreatedAt)">{{ formatDate(datasetScript.CreatedAt) }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="danger" :icon="mdiTrashCan" small :target-id="datasetScript.ID" @clicked="deleteButtonClicked" />
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
            <small>{{ $t('tables.page', {page: currentPageHuman, count: numPages}) }}</small>
        </BaseLevel>
    </div>
</template>s