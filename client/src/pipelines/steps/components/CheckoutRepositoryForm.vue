<script setup>
    import { checkoutRepo } from '@/pipelines/steps';
</script>

<template>
    <GroupElement :name="checkoutRepo.name" :conditions="checkoutRepo.conditions" :label="checkoutRepo.label">
        <TextElement 
            :name="checkoutRepo.fields.repoURL.name" 
            :label="checkoutRepo.fields.repoURL.label" 
            :rules="checkoutRepo.fields.repoURL.rules" />
    </GroupElement>
</template>