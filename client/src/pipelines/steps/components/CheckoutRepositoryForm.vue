<script setup>
    import { CheckoutRepo } from '@/pipelines/steps';
</script>

<template>
    <GroupElement :name="CheckoutRepo.name" :conditions="CheckoutRepo.conditions" :label="CheckoutRepo.label">
        <TextElement 
            :name="CheckoutRepo.fields.repoURL.name" 
            :label="CheckoutRepo.fields.repoURL.label" 
            :rules="CheckoutRepo.fields.repoURL.rules" />
    </GroupElement>
</template>