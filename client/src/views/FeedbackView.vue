<script setup>
import { storeToRefs } from "pinia";
import { ref, computed, watch } from "vue";
import { useAsyncState } from "@vueuse/core";
import { doRequest } from "@/util";
import { useAuthStore } from "@/stores/auth";
import { useRouter, useRoute } from 'vue-router';
import {
    mdiMessageAlertOutline,
} from "@mdi/js";
import SectionMain from "@/components/SectionMain.vue";
import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
import BaseButtons from "@/components/BaseButtons.vue";
import BaseButton from "@/components/BaseButton.vue";
import FeedbackQuery from '@/components/FeedbackQuery.vue';
import Carousel from 'primevue/carousel';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import Loading from "vue-loading-overlay";
import { i18n } from '@/i18n';

const { accessToken } = storeToRefs(useAuthStore());
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

const isLoading = computed(() => isFetching.value || isUpdating.value);

const onFeedbackCancel = () => {
    router.push(`/runresults/${route.params.id}`);
}

const onFeedbackSave = () => {
    const updateData = {
        humanFeedbackQueries: queries.value.map(query => {
            return {
                humanFeedbackQueryID: query.HumanFeedbackQuery.ID,
                rects: query.HumanFeedbackRects.map(rect => {
                    return {
                        rectID: rect.ID,
                        selected: rect.Selected,
                    }
                })
            }
        })
    };

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
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiMessageAlertOutline" :title="pageTitle" main>
                <BaseButtons style="float:right">
                    <BaseButton :label="$t('buttons.submit')" color="success" @click="onFeedbackSave" />
                    <BaseButton :label="$t('buttons.cancel')" color="danger" @click="onFeedbackCancel" />
                </BaseButtons>
            </SectionTitleLineWithButton>
            <Carousel :value="queries" :numVisible="1" :numScroll="1">
                <template #item="slotProps">
                    <FeedbackQuery :query="slotProps.data" :queryCount="queries.length" @checked="onQueryRectChecked" />
                </template>
            </Carousel>
        </SectionMain>
        <Toast />
    </LayoutAuthenticated>
</template>

<style>
.p-carousel-prev.p-link:not(.p-disabled) {
    background-color: rgb(37 99 235 / var(--tw-bg-opacity)) !important;
    color: white !important;
}

.p-carousel-next.p-link:not(.p-disabled) {
    background-color: rgb(37 99 235 / var(--tw-bg-opacity)) !important;
    color: white !important;
}
</style>