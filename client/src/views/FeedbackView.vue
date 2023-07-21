<script setup>
    import { Panel, PanelPosition, VueFlow, isNode, useVueFlow } from '@vue-flow/core';
    import { storeToRefs } from "pinia";
    import { ref, computed, watch, onMounted } from "vue";
    import { useAsyncState } from "@vueuse/core";
    import { doRequest, deepFilterMenuBarSteps } from "@/util";
    import { useAuthStore } from "@/stores/auth";
    import { useRouter, useRoute } from 'vue-router';
    import {
        mdiChartTimelineVariant,
        mdiFileDocumentOutline
    } from "@mdi/js";
    import SectionMain from "@/components/SectionMain.vue";
    import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
    import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import FeedbackQuery from '@/components/FeedbackQuery.vue';
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
    const router = useRouter();
    const queries = ref([]);
    const toast = useToast();
    const { t } = i18n.global;

    const pageTitle = ref('');

    // FETCH FEEDBACK
    const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchFeedback } = useAsyncState(
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

    // UPDATE FEEDBACK
    const { isLoading: isUpdating, state: updateResponse, isReady: isUpdateFinished, execute: updateFeedback } = useAsyncState(
        (data) => {
            return doRequest({
                url: `/api/feedback/${route.params.id}`,
                method: 'POST',
                data: data,
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
    )

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
        pageTitle.value = t('pages.runs.feedback.header', { runID: route.params.id, stepName: queries.value[0].RunStepStatus.Name })
    })

    watch(updateResponse, (value) => {
        if (value.error) {
            let header = t('global.errors.generic.header');
            let detail = value.error.message;

            if (value.status == 401) {
                header = t('global.errors.authorization.header');
                detail = t('global.errors.authorization.detail');
            }

            toast.add({ severity: 'error', summary: header, detail: detail, life: 3000 });
        }

        if (value.status === 200) {
            router.push(`/runresults/${route.params.id}`);
        }
    })

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

    const onFeedbackCancel = () => {
        router.push(`/runresults/${route.params.id}`);
    }

    const onFeedbackSave = () => {
        const updateData = queries.value.map(query => {
            return {
                humanFeedbackQueryID: query.HumanFeedbackQuery.ID,
                rects: query.HumanFeedbackRects.map(rect => {
                    return {
                        rectID: rect.ID,
                        selected: rect.Selected,
                    }
                })
            }
        });

        updateFeedback(null, updateData);
    }

    const onQueryRectChecked = (queryID, rectsSelected) => {
        const queryIndex = queries.value.findIndex(query => query.ID == queryID);

        if (queryIndex > -1) {
            queries.value[queryIndex].HumanFeedbackRects.forEach((rect, rectIndex) => {
                const selectedRectIndex = rectsSelected.findIndex(selectedRect => selectedRect.ID == rect.ID)

                if (selectedRectIndex > -1) {
                    queries.value[queryIndex].HumanFeedbackRects[rectIndex].Selected = true
                } else {
                    queries.value[queryIndex].HumanFeedbackRects[rectIndex].Selected = false
                }
            })
        }
    }

</script>

<template>
    <LayoutAuthenticated>
        <SectionMain>
            <loading v-model:active="isLoading" :is-full-page="false" />
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiChartTimelineVariant" :title="pageTitle" main>
                <BaseButtons style="float:right">
                    <BaseButton :label="$t('buttons.submit')" color="success" @click="onFeedbackSave" />
                    <BaseButton :label="$t('buttons.cancel')" color="danger" @click="onFeedbackCancel" />
                </BaseButtons>
            </SectionTitleLineWithButton>
            <Carousel :value="queries" :numVisible="1" :numScroll="1" :responsiveOptions="responsiveOptions">
                <template #item="slotProps">
                    <FeedbackQuery :query="slotProps.data" @checked="onQueryRectChecked" />
                </template>
            </Carousel>
        </SectionMain>
        <Toast />
    </LayoutAuthenticated>
</template>