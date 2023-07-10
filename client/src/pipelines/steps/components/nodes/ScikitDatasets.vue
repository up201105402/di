<script setup>
    import { Handle, Position, NodeIdInjection } from '@vue-flow/core';
    import { NodeToolbar } from '@vue-flow/node-toolbar';
    import { computed } from 'vue';
    import BaseIcon from "@/components/BaseIcon.vue";
    import { mdiCloseCircle, mdiAlertCircle, mdiCancel, mdiCheckCircleOutline, mdiDotsCircle } from "@mdi/js";
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

</script>

<template>
    <NodeToolbar v-if="!props.data.readonly" style="display: flex; gap: 0.5rem; align-items: center" :is-visible="data.toolbarVisible"
        :position="Position.Top">
        <button class="lg" @click.prevent="onDeleteClick">
            <BaseIcon :path="mdiCloseCircle " />
        </button>
    </NodeToolbar>

    <div class="node-type">
        <span class="node-id">{{ parseInt(props.data.id) + 1 }}</span>
        <img class="scikit-logo" style="display: inline" src="/assets/1200px-scikit_learn_logo.png" width="40"
            height="22">
        <span class="node-type-label">{{ camel2title(props.data.type) }}</span>
    </div>
    <div>
        <span>{{ props.label }}</span>
        <span v-if="props.data.status"><BaseIcon :path="getIcon()" /></span>
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
</style>