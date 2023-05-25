<script setup>
    import { Handle, Position, NodeIdInjection } from '@vue-flow/core';
    import { computed } from 'vue';

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
        <div>{{ props.data.nameAndType.nodeName }}</div>

        <Handle v-if="props.data.isFirstStep === false" id="a" type="source" :position="Position.Left"
            :style="sourceHandleStyleA" />
        <Handle id="b" type="target" :position="Position.Right" :style="sourceHandleStyle" />
    </div>
</template>