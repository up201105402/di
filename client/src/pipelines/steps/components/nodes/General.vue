<script setup>
    import { Handle, Position } from '@vue-flow/core';
    import { NodeToolbar } from '@vue-flow/node-toolbar';
    import { mdiCloseCircle } from "@mdi/js";
    import { computed } from 'vue';
    import BaseIcon from "@/components/BaseIcon.vue";
    import $ from 'jquery';
    import Tag from 'primevue/tag';
    import { getStatusTagSeverity } from '@/util';

    const props = defineProps({
        id: {
            type: String,
            required: true,
        },
        label: {
            type: String,
            required: true,
        },
        events: {
            type: Object,
            required: true,
        },
        data: {
            type: Object,
            required: true,
        },
    })

    defineOptions({
        inheritAttrs: false
    })

    const emit = defineEmits(['change', 'gradient'])

    const onDeleteClick = (e) => {
        $(document).trigger('onNodeDelete', { id: props.data.id });
    }

    const sourceHandleStyle = computed(() => ({
        backgroundColor: props.data.color,
        filter: 'invert(100%)',
    }))

    const outputHandleStyle = computed(() => ({
        backgroundColor: props.data.color,
        filter: 'invert(100%)',
    }))

</script>

<template>
    <NodeToolbar v-if="!props.data.readonly" style="display: flex; gap: 0.5rem; align-items: center"
        :is-visible="data.toolbarVisible" :position="Position.Top">
        <button class="lg" @click.prevent="onDeleteClick">
            <BaseIcon :path="mdiCloseCircle " />
        </button>
    </NodeToolbar>

    <div class="node-type">
        <span class="node-id">{{ parseInt(props.id) + 1 }}</span>
        <span class="node-type-label">{{ $t('pages.pipelines.steps.' + props.data.type) }}</span>
    </div>
    <div>
        <span>{{ props.label }}</span>
        <span class="node-status-tag" v-if="props.data.status"><Tag :severity="getStatusTagSeverity(props.data.status.id)" :value="props.data.status.Name" /></span>
        <Handle v-if="props.data.nameAndType.isFirstStep === false" :id="props.id + '_input'" type="source"
            :position="Position.Left" :style="sourceHandleStyle" />
        <Handle :id="props.id + '_output'" type="target" :position="Position.Right" :style="outputHandleStyle" />
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

    .node-type-label {
        color: blue;
        padding: 0 2px;
    }

    .node-config {
        text-transform: none;
    }
</style>