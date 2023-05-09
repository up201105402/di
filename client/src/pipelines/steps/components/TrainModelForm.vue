<script setup>
    import { trainModel } from '@/pipelines/steps';
</script>

<template>
    <GroupElement 
        :name="trainModel.name" 
        :label="trainModel.label" 
        :conditions="trainModel.conditions">
        <TextElement 
            :name="trainModel.fields.modelDirectory.name" 
            :label="trainModel.fields.modelDirectory.label" 
            :rules="trainModel.fields.modelDirectory.rules" />
        <TextElement 
            :name="trainModel.fields.epochs.name" 
            :label="trainModel.fields.epochs.label"
            input-type="number" 
            :rules="trainModel.fields.epochs.rules" 
            autocomplete="off" />
    </GroupElement>
</template>