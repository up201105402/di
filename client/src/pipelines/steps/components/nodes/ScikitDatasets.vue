<script setup>
    import { Handle, Position, NodeIdInjection } from '@vue-flow/core';
    import { NodeToolbar } from '@vue-flow/node-toolbar';
    import { computed } from 'vue';
    import BaseIcon from "@/components/BaseIcon.vue";
    import { mdiCloseCircle } from "@mdi/js";
    import { camel2title } from '@/util';

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
    
</script>

<template>
    <NodeToolbar style="display: flex; gap: 0.5rem; align-items: center" :is-visible="data.toolbarVisible"
        :position="Position.Top">
        <button class="lg" @click.prevent="onDeleteClick">
            <BaseIcon :path="mdiCloseCircle " />
        </button>
    </NodeToolbar>

    <div class="node-type">
        <span class="node-id">{{ parseInt(props.id) + 1 }}</span>
        <span class="node-type-label">Checkout Repo</span>
        <span class="scikit-logo"><img src="/assets/1200px-scikit_learn_logo.png" width="20px" height="11px"></span>
    </div>
    <div>
        <div>{{ props.label }}</div>

        <Handle v-if="props.data.isFirstStep === false" id="a" type="source" :position="Position.Left"
            :style="sourceHandleStyle" />
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

    .node-type-label {
        color: blue;
        padding: 0 2px;
    }
</style>