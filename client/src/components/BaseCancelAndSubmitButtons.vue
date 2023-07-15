<script setup>
    import BaseButtons from '@/components/BaseButtons.vue';
    import BaseButton from '@/components/BaseButton.vue';

    const props = defineProps({
        cancelButtonId: String,
        submitEnabled: Boolean,
    })
</script>
<template>
    <BaseButtons>
        <BaseButton :id="cancelButtonId" :label="'Cancel'" color="danger" />
        <BaseButton :disabled="!submitEnabled" label="Submit" color="success" type="submit" />
    </BaseButtons>
</template>