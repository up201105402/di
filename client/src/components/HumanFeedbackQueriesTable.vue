<script setup>
import { mdiEye, mdiMessageAlertOutline } from "@mdi/js";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import BaseLevel from "@/components/BaseLevel.vue";
import Tag from 'primevue/tag';
import { useToast } from 'primevue/usetoast';
import { ref, computed } from "vue";
import { formatDate, getStatusTagSeverity } from '@/util';
import { i18n } from '@/i18n';

const toast = useToast();
const { t } = i18n.global;

const props = defineProps({
    rows: {
        type: Array,
        required: true,
    },
    runID: {
        type: Number,
        required: true,
    }
});

const perPage = ref(5);
const currentPage = ref(0);
const fetchedRows = ref([]);

const paginatedRows = computed(() => {
    return props.rows ? props.rows.slice(
        perPage.value * currentPage.value,
        perPage.value * (currentPage.value + 1)
    ) : [];
});

const numPages = computed(() => Math.ceil(fetchedRows.value.length ? fetchedRows.value.length / perPage.value : props.rows.length / perPage.value));
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
                <th>{{ $t('pages.runs.results.humanFeedbackQueries.table.headers.id') }}</th>
                <th>{{ $t('pages.runs.results.humanFeedbackQueries.table.headers.status') }}</th>
                <th>{{ $t('pages.runs.results.humanFeedbackQueries.table.headers.created') }}</th>
                <th>{{ $t('pages.runs.results.humanFeedbackQueries.table.headers.updated') }}</th>
                <th>{{ $t('pages.runs.results.humanFeedbackQueries.table.headers.stepID') }}</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <tr v-for="row in paginatedRows" :key="row.ID">
                <td class="border-b-0 lg:w-6 before:hidden">
                    {{ row.ID }}
                </td>
                <td class="border-b-0 lg:w-6 before:hidden">
                    <Tag :severity="getStatusTagSeverity(row.QueryStatusID)" :value="row.QueryStatus.Name" />
                </td>
                <td :data-label="$t('pages.runs.results.humanFeedbackQueries.table.headers.created')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.CreatedAt)">{{ formatDate(row.CreatedAt) }}</small>
                </td>
                <td :data-label="$t('pages.runs.results.humanFeedbackQueries.table.headers.updated')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="formatDate(row.UpdatedAt)">{{ formatDate(row.UpdatedAt) }}</small>
                </td>
                <td :data-label="$t('pages.runs.results.humanFeedbackQueries.table.headers.stepID')" class="lg:w-1 whitespace-nowrap">
                    <small class="text-gray-500 dark:text-slate-400" :title="row.StepID">{{ formatDate(row.StepID) }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton v-if="row.QueryStatusID == 3" :to="`/feedback/${runID}/query/${row.ID}`" :icon="mdiEye" color="info" />
                        <BaseButton v-if="row.QueryStatusID < 3" :to="`/feedback/${runID}/query/${row.ID}`" :icon="mdiMessageAlertOutline" color="warning" />
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
            <small>{{ $t('tables.page', { page: currentPageHuman, count: numPages }) }}</small>
        </BaseLevel>
    </div>
</template>