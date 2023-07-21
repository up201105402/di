<script setup>
import { computed, ref } from "vue";    
import BaseLevel from "@/components/BaseLevel.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import { formatDate } from '@/util';

const props = defineProps({
    items: Array,
    checkedRows: {
        type: Array,
        default: []
    },
    queryID: Number,
    checkable: Boolean,
});

// EMITS
const emit = defineEmits(["checked"]);

const count = ref(0);

// ITEMS PROCESSING
// ITEMS PROCESSING

// ITEMS PROCESSING

const perPage = ref(5);
const currentPage = ref(0);

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

const checked = (item) => {
    emit('checked', props.queryID, item);
};

const isRowChecked = (itemID) => {
    return props.checkedRows.find(row => row.ID == itemID) != null;
}

</script>

<template>
    <div v-if="checkedRows.length" class="p-3 bg-gray-100/50 dark:bg-slate-800">
        <span v-for="checkedRow in checkedRows" :key="checkedRow.ID"
            class="inline-block px-2 py-1 rounded-sm mr-2 text-sm bg-gray-100 dark:bg-slate-700">
            {{ checkedRow.ID }}
        </span>
    </div>

    <table>
        <thead>
            <tr>
                <th v-if="checkable" />
                <th>{{ $t('pages.runs.feedback.table.headers.id') }}</th>
                <th>{{ $t('pages.runs.feedback.table.headers.x1') }}</th>
                <th>{{ $t('pages.runs.feedback.table.headers.y1') }}</th>
                <th>{{ $t('pages.runs.feedback.table.headers.x2') }}</th>
                <th>{{ $t('pages.runs.feedback.table.headers.y2') }}</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="item in itemsPaginated" :key="item.ID">
                <td class="lg:w-1 whitespace-nowrap">
                    <label class="checkbox">
                        <input :checked="isRowChecked(item.ID)" type="checkbox" @click="(e) => checked(item)" />
                        <span class="check" />
                    </label>
                </td>
                <!-- <TableCheckboxCell v-if="checkable" @checked="checked($event, item)" :isChecked="isRowChecked(item.ID)" /> -->
                <td :data-label="$t('pages.runs.feedback.table.headers.id')">
                    {{ item.ID }}
                </td>
                <td :data-label="$t('pages.runs.feedback.table.headers.x1')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(item.X1)">{{ item.X1 }}</small>
                </td>
                <td :data-label="$t('pages.runs.feedback.table.headers.y1')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(item.Y1)">{{ item.Y1 }}</small>
                </td>
                <td :data-label="$t('pages.runs.feedback.table.headers.x2')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(item.X2)">{{ item.X2 }}</small>
                </td>
                <td :data-label="$t('pages.runs.feedback.table.headers.y2')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(item.Y2)">{{ item.Y2 }}</small>
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
            <small>{{ $t('tables.page', { page: currentPageHuman, count: numPages }) }}</small>
        </BaseLevel>
    </div>
</template>