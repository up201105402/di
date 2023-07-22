<script setup>
import { computed, ref } from "vue";
import { mdiEye, mdiTrashCan } from "@mdi/js";
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

</script>

<template>
    <table>
        <thead>
            <tr>
                <th>{{ $t('pages.pipelines.table.headers.id') }}</th>
                <th>{{ $t('pages.pipelines.table.headers.name') }}</th>
                <th>{{ $t('pages.pipelines.table.headers.modified') }}</th>
                <th>{{ $t('pages.pipelines.table.headers.created') }}</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <tr v-for="pipeline in itemsPaginated" :key="pipeline.id">
                <td :data-label="$t('pages.pipelines.table.headers.name')">
                    {{ pipeline.ID }}
                </td>
                <td :data-label="$t('pages.pipelines.table.headers.name')">
                    {{ pipeline.name }}
                </td>
                <td :data-label="$t('pages.pipelines.table.headers.modified')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(pipeline.UpdatedAt)">{{ formatDate(pipeline.UpdatedAt) }}</small>
                </td>
                <td :data-label="$t('pages.pipelines.table.headers.created')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(pipeline.CreatedAt)">{{ formatDate(pipeline.CreatedAt) }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="info" :icon="mdiEye" small :to="'/pipelines/edit/' + pipeline.ID" />
                        <BaseButton color="danger" :icon="mdiTrashCan" small :target-id="pipeline.ID" @clicked="deleteButtonClicked" />
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
</template>