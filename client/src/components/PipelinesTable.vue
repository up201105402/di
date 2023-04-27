<script setup>
    import { computed, ref, watch } from "vue";
    import { usePipelinesStore } from "@/stores/pipelines";
    import { mdiEye, mdiTrashCan } from "@mdi/js";
    import CardBoxModal from "@/components/CardBoxModal.vue";
    import TableCheckboxCell from "@/components/TableCheckboxCell.vue";
    import BaseLevel from "@/components/BaseLevel.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import UserAvatar from "@/components/UserAvatar.vue";

    const props = defineProps({
        items: Array,
        checkable: Boolean,
    });

    const pipelinesStore = usePipelinesStore();

    const pipelines = ref(props.items.value)

    watch(props.items, () => pipelines = props.items);

    const isModalActive = ref(false);

    const isModalDangerActive = ref(false);

    const idToDelete = ref(0);

    const perPage = ref(5);

    const currentPage = ref(0);

    const checkedRows = ref([]);

    const itemsPaginated = computed(() =>
        props.items.slice(
            perPage.value * currentPage.value,
            perPage.value * (currentPage.value + 1)
        )
    );

    const numPages = computed(() => Math.ceil(props.items.length / perPage.value));

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

    const deletePipeline = (e) => {
        console.error(e);
    }
</script>

<template>
    <CardBoxModal v-model="isModalActive" title="Sample modal">
        <p>Lorem ipsum dolor sit amet <b>adipiscing elit</b></p>
        <p>This is sample modal</p>
    </CardBoxModal>

    <CardBoxModal v-model="isModalDangerActive" title="Confirm Delete" :id="idToDelete" @confirm="deletePipeline" button="danger" has-cancel>
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
                    <small class="text-gray-500 dark:text-slate-400" :title="pipeline.CreatedAt">{{ pipeline.CreatedAt }}</small>
                </td>
                <td class="before:hidden lg:w-1 whitespace-nowrap">
                    <BaseButtons type="justify-start lg:justify-end" no-wrap>
                        <BaseButton color="info" :icon="mdiEye" small :to="'/pipeline/' + pipeline.id"/>
                        <BaseButton color="danger" :icon="mdiTrashCan" small @click="() => { isModalDangerActive = true; idToDelete.value = pipeline.id }" />
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