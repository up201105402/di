<script setup>
    import { Handle, Position } from '@vue-flow/core';
    import { NodeToolbar } from '@vue-flow/node-toolbar';
    import { computed } from 'vue';
    import BaseIcon from "@/components/BaseIcon.vue";
    import { mdiCloseCircle, mdiAlertCircle, mdiCancel, mdiCheckCircleOutline, mdiDotsCircle, mdiPowershell  } from "@mdi/js";
    import $ from 'jquery';
    import { getStatusTagSeverity } from '@/util';

    const props = defineProps({
        data: {
            type: Object,
            required: true,
        },
        label: {
            type: String,
            required: true,
        },
    })

    defineOptions({
        inheritAttrs: false
    })

    const emit = defineEmits(['change', 'gradient'])

    function onSelect(color) {
        emit('change', color)
    }

    function onGradient() {
        emit('gradient')
    }

    const sourceHandleStyle = computed(() => ({
        backgroundColor: props.data.color,
        filter: 'invert(100%)',
    }))

    const outputHandleStyle = computed(() => ({
        backgroundColor: props.data.color,
        filter: 'invert(100%)',
    }))

    const getIcon = () => {
        if (props.data.status == 1) {
            return mdiCancel;
        } else if (props.data.status == 2) {
            return mdiDotsCircle;
        } else if (props.data.status == 3) {
            return mdiAlertCircle;
        } else if (props.data.status == 4) {
            return mdiCheckCircleOutline;
        }
    }

    const onDeleteClick = (e) => {
        $(document).trigger('onNodeDelete', { id: props.data.id });
    }

</script>

<template>
    <NodeToolbar v-if="!props.data.readonly" style="display: flex; gap: 0.5rem; align-items: center" :is-visible="data.toolbarVisible" :position="Position.Top">
        <button class="lg" @click.prevent="onDeleteClick">
            <BaseIcon :path="mdiCloseCircle " />
        </button>
    </NodeToolbar>

    <div class="node-type">
        <span class="node-id">{{ parseInt(props.data.id) + 1 }}</span>
        <BaseIcon class="scikit-logo" style="display: inline" src="/assets/Python-logo.png" :path="mdiPowershell " />
        <span class="node-type-label">{{ $t('pages.pipelines.steps.' + props.data.type) }}</span>
    </div>
    <div class="node-config" >
        <span>{{ props.label }}</span>
        <span class="node-status-tag" v-if="props.data.status"><Tag :severity="getStatusTagSeverity(props.data.status.id)" :value="props.data.status.Name" /></span>
        <Handle v-if="props.data.nameAndType.isFirstStep === false" id="a" type="source" :position="Position.Left" :style="sourceHandleStyle" />
        <Handle id="b" type="target" :position="Position.Right" :style="sourceHandleStyle" />
    </div>
</template>


<style>
    .node-type {
        text-transform: none;
        color: blue;
        border-bottom: 2px solid;
    }

    .node-id {
        color: blue;
        padding: 0 2px;
        border-right: 2px solid;
    }

    .scikit-logo {
        display: inline;
    }

    .node-type-label {
        color: blue;
        padding: 0 2px;
    }

    .node-config {
        text-transform: none;
    }
</style>