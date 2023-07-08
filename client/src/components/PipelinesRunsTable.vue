<script setup>
    import { computed, ref } from "vue";
    import BaseLevel from "@/components/BaseLevel.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import PipelineRunsTableRow from "@/components/PipelineRunsTableRow.vue";

    const props = defineProps({
        items: Array,
        checkable: Boolean,
    });
    
    // ITEMS PROCESSING

    const perPage = ref(5);

    const currentPage = ref(0);

    const checkedRows = ref([]);

    const openedPipelines = ref([]);

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
                <th>Last Run</th>
                <th />
            </tr>
        </thead>
        <tbody>
            <PipelineRunsTableRow v-for="pipeline in itemsPaginated" :parentRow="pipeline" subrowsFetchBaseURL="/api/run/" />
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