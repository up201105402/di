<script setup>
    import { Panel, PanelPosition, VueFlow, isNode, useVueFlow } from '@vue-flow/core';
    import { storeToRefs } from "pinia";
    import { ref, computed, watch, onMounted } from "vue";
    import { useAsyncState } from "@vueuse/core";
    import { doRequest, deepFilterMenuBarSteps } from "@/util";
    import { useAuthStore } from "@/stores/auth";
    import { useRoute } from 'vue-router';
    import {
        mdiChartTimelineVariant,
        mdiFileDocumentOutline
    } from "@mdi/js";
    import SectionMain from "@/components/SectionMain.vue";
    import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
    import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
    import CardBoxModal from '@/components/CardBoxModal.vue';
    import UpsertStepDialog from '@/components/UpsertStepDialog.vue';
    import FeedbackRectsTable from '@/components/FeedbackRectsTable.vue';
    import FeedbackCanvas from '@/components/FeedbackCanvas.vue';
    import Carousel from 'primevue/carousel';
    import Tag from 'primevue/tag';
    import Toast from 'primevue/toast';
    import { useToast } from 'primevue/usetoast';
    import Loading from "vue-loading-overlay";
    import { camel2title } from '@/util';
    import Editor from 'primevue/editor';
    import { i18n } from '@/i18n';
    import $ from 'jquery';

    const { accessToken, requireAuthRoute } = storeToRefs(useAuthStore());
    const route = useRoute();
    const queries = ref([]);
    const toast = useToast();
    const { t } = i18n.global;

    const pageTitle = ref(t('pages.runs.feedback.header', { runID: route.params.id }));

    // FETCH RUN

    const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchRun } = useAsyncState(
        () => {
            return doRequest({
                url: `/api/feedback/${route.params.id}`,
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
        },
    );

    watch(fetchResponse, (value) => {
        if (value.error) {
            let header = t('global.errors.generic.header');
            let detail = value.error.message;

            if (value.status == 401) {
                header = t('global.errors.authorization.header');
                detail = t('global.errors.authorization.detail');
            }

            toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
        }
    })

    watch(isFetchFinished, () => {
        queries.value = fetchResponse.value?.data ? fetchResponse.value?.data.queries : []
    })

    const onRectChecked = (rectID, queryID) => {
        queries.value.forEach(query => {
            if (query.HumanFeedbackQuery.QueryID == queryID) {
                query.HumanFeedbackRects.forEach(rect => {
                    if (rect.ID == rectID) {
                        rect.Selected = !rect.Selected
                    }
                })
            }
        });
    }

    // const queries = ref([
    //     {
    //         id: '1000',
    //         code: 'f230fh0g3',
    //         name: 'Bamboo Watch',
    //         description: 'Product Description',
    //         image: 'bamboo-watch.jpg',
    //         price: 65,
    //         category: 'Accessories',
    //         quantity: 24,
    //         inventoryStatus: 'INSTOCK',
    //         rating: 5
    //     },
    //     {
    //         id: '1001',
    //         code: 'nvklal433',
    //         name: 'Black Watch',
    //         description: 'Product Description',
    //         image: 'black-watch.jpg',
    //         price: 72,
    //         category: 'Accessories',
    //         quantity: 61,
    //         inventoryStatus: 'INSTOCK',
    //         rating: 4
    //     },
    // ]);

    const responsiveOptions = ref([
        {
            breakpoint: '767px',
            numVisible: 1,
            numScroll: 1
        }
    ]);

    const getSeverity = (status) => {
        switch (status) {
            case 'INSTOCK':
                return 'success';

            case 'LOWSTOCK':
                return 'warning';

            case 'OUTOFSTOCK':
                return 'danger';

            default:
                return null;
        }
    };

    const isLoading = computed(() => isFetching.value);

    const onImageClicked = (queryID, ptX, ptY) => {
        const queryIndex = queries.value.findIndex(query => query.HumanFeedbackQuery.ID == queryID);
        const rectIndex = queries.value[queryIndex].HumanFeedbackRects.findIndex(rect => {
            return (ptX >= rect.X1) && (ptX <= rect.X2) && (ptY >= rect.Y1) && (ptY <= rect.Y2);
        });

        const rect = queries.value[queryIndex].HumanFeedbackRects[rectIndex]

        if (rect) {
            console.error(rect)
            queries.value[queryIndex].HumanFeedbackRects[rectIndex].Selected = !rect.Selected
        }
    }

</script>

<template>
    <LayoutAuthenticated>
        <SectionMain>
            <loading v-model:active="isLoading" :is-full-page="false" />
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiChartTimelineVariant" :title="pageTitle" main />
            <Carousel :value="queries" :numVisible="1" :numScroll="1" :responsiveOptions="responsiveOptions">
                <template #item="slotProps">
                    <div class="border-1 surface-border border-round m-2 text-center py-5 px-3">
                        <FeedbackCanvas :id="slotProps.data.HumanFeedbackQuery.ID" :imageURL="slotProps.data.ImageURL" @mouseClick="onImageClicked" />
                        <div>
                            <h4 class="mb-1">Query {{ slotProps.data.HumanFeedbackQuery.QueryID }}</h4>
                            <h4 class="mb-1">{{ slotProps.data.RunStepStatus.Name }} ({{ slotProps.data.RunStepStatus.ID }})</h4>
                            <Tag :value="slotProps.data.inventoryStatus" :severity="getSeverity(slotProps.data.inventoryStatus)" />
                        </div>
                        <FeedbackRectsTable :items="slotProps.data.HumanFeedbackRects" :queryID="slotProps.data.HumanFeedbackQuery.QueryID" @checked="onRectChecked" checkable />
                    </div>
                </template>
            </Carousel>
        </SectionMain>
        <Toast />
    </LayoutAuthenticated>
</template>

<style>
    #run-results-log .p-editor-toolbar {
        display: none
    }

    #run-results-log div.p-editor-container div.ql-editor p {
        max-width: 100%;
    }
</style>