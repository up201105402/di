<script setup>
    import { Handle, Position } from '@vue-flow/core';
    import { NodeToolbar } from '@vue-flow/node-toolbar';
    import { computed } from 'vue';

    const props = defineProps({
        id: {
            type: String,
            required: true,
        },
        label: {
            type: String,
            required: true,
        },
        data: {
            type: Object,
            required: true,
        },
    })

    const emit = defineEmits(['change', 'gradient'])

    function onSelect(color) {
        emit('change', color);
    }

    function onGradient() {
        emit('gradient');
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
    <NodeToolbar style="display: flex; gap: 0.5rem; align-items: center" :is-visible="data.toolbarVisible" :position="Position.Top">
        <button>Action1</button>
        <button>Action2</button>
        <button>Action3</button>
    </NodeToolbar>

    <div class="nodeType">
        <span class="nodeId">{{ parseInt(props.id) + 1 }}</span>
        <span class="nodeTypeLabel">Checkout Repo</span>
    </div>
    <div>
        <div>{{ props.label }}</div>

        <Handle v-if="props.data.isFirstStep === false" :id="props.id + '_input'" type="source"
            :position="Position.Left" :style="sourceHandleStyle" />
        <Handle :id="props.id + '_output'" type="target" :position="Position.Right" :style="outputHandleStyle" />
    </div>
</template>

<style>
    .nodeType {
        text-transform: none;
        color: blue;
        border-bottom: 2px solid;
    }

    .nodeId {
        color: blue;
        padding: 0 2px;
        border-right: 2px solid;
    }

    .nodeTypeLabel {
        color: blue;
        padding: 0 2px;
    }
</style>