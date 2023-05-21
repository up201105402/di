<script setup>
    import { Handle, Position } from '@vue-flow/core';
    import { computed } from 'vue';

    const props = defineProps({
        data: {
            type: Object,
            required: true,
        },
        label: {
            type: String,
            required: false,
        }
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
    <div>
        <div>{{ props.label }}</div>

        <Handle id="a" type="source" :position="Position.Left" :style="sourceHandleStyle" />
        <Handle id="b" type="target" :position="Position.Right" :style="outputHandleStyle" />
    </div>
</template>