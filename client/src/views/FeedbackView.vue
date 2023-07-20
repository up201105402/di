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
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import Loading from "vue-loading-overlay";
import { nodeTypes, menubarSteps } from "@/pipelines/steps";
import { camel2title } from '@/util';
import Editor from 'primevue/editor';
import { i18n } from '@/i18n';

const { accessToken, requireAuthRoute } = storeToRefs(useAuthStore());
const route = useRoute();
const elements = ref([]);
const log = ref('');
const runTitle = ref('');
const toast = useToast();
const { t } = i18n.global;

const quillModules = {
    "toolbar": false
}

// FETCH RUN

const { isLoading: isFetching, state: fetchResponse, isReady: isFetchFinished, execute: fetchRun } = useAsyncState(
    () => {
        return doRequest({
            url: `/api/runresults/${route.params.id}`,
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

const selectedStep = ref();

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

const parseRunDefinition = (run) => {
    try {
        return JSON.parse(run.Definition)
    } catch (e) {
        return [];
    }
}

const formatValue = (value) => {
    return '<p>' + value.replaceAll('<br>', '').replaceAll("\n", "</p><p>").replaceAll("\r", "</p><p>") + "</p>"
}

watch(isFetchFinished, () => {
    elements.value = fetchResponse.value?.data ? parseRunDefinition(fetchResponse.value.data.run) : [];
    elements.value.forEach(element => {
        element.data.readonly = true;
        fetchResponse.value.data.runStepStatuses.forEach(stepStatus => {
            if (element.data.id == stepStatus.StepID) {
                element.data.status = stepStatus.RunStatusID
            }
        })
    });
    runTitle.value = fetchResponse.value?.data ? t('pages.runs.results.header', { pipelineName: fetchResponse.value.data.run.Pipeline.name, runID: fetchResponse.value.data.run.ID }) : t('global.untitled');
    log.value = fetchResponse.value?.data ? formatValue(fetchResponse.value.data.log) : "";
})

const isLoading = ref(false) //computed(() => isFetching.value);

const isStepDialogActive = ref(false);
const formSchema = ref({});
const dialogTitle = ref('');
const stepData = ref({});
const editStepNodeId = ref("-1");
let count = 0;

const getSchemaFromType = (searchGroup, searchType) => {
    const group = steps.filter(e => e.type == searchGroup)[0]
    return group.steps.filter(e => e.type == searchType)[0];
}

const onNodeDoubleClick = (e) => {
    isStepDialogActive.value = true;
    editStepNodeId.value = e.node.id;
    stepData.value = { ...e.node.data }
    const formkitObject = deepFilterMenuBarSteps(menubarSteps, 'type', e.node.data.type).form({ ...e.node.data }, () => { }, false);
    formSchema.value = formkitObject.formSchema;
    const label = t('pages.pipelines.steps.' + e.node.data.type);
    dialogTitle.value = t('pages.runs.results.dialog.edit.header', { name: label });
    stepData.value = formkitObject.formkitData;
    count++;
}

const onCancel = () => {
    isStepDialogActive.value = false;
    count++;
}

const init = function () {
    canvas.addEventListener('click', mouseClick, false);
}

const mouseClick = (e) => {
    rect.startX = e.pageX - e.target.offsetLeft;
    rect.startY = e.pageY - e.target.offsetTop;
    console.error("rect.startX = " + rect.startX);
    console.error("rect.startY = " + rect.startY);
}

var canvas
var ctx
var imageObj
var rect = {}

onMounted(() => {
    canvas = document.getElementById('myCanvas');
    ctx = canvas.getContext('2d');
    imageObj = new Image();

    imageObj.onload = function (e) {
        ctx.drawImage(imageObj, 0, 0, imageObj.width, imageObj.height);
    };
    imageObj.src = 'http://localhost:8001/work/2/31/my_plot.png';

    // ctx.globalAlpha = 0.5;

    init();
});

</script>
  
<template>
    <LayoutAuthenticated>
        <SectionMain>
            <loading v-model:active="isLoading" :is-full-page="false" />
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiChartTimelineVariant" :title="runTitle" main />
            <canvas id="myCanvas" width="1000" height="500"></canvas>
            <SectionTitleLineWithButton :hasButton="false" :icon="mdiFileDocumentOutline" :title="$t('pages.runs.results.log.header')" />
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