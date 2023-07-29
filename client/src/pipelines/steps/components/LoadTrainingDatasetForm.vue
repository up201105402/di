<script setup>
    import { loadTrainingDataset } from '@/pipelines/steps';
</script>

<template>
    <GroupElement :name="loadTrainingDataset.name" :conditions="loadTrainingDataset.conditions" :label="loadTrainingDataset.label">
        <TextElement 
            :name="loadTrainingDataset.fields.trainingDatasetDirectory.name" 
            :label="loadTrainingDataset.fields.trainingDatasetDirectory.label" 
            :rules="loadTrainingDataset.fields.trainingDatasetDirectory.rules" />
        <TextElement 
            :name="loadTrainingDataset.fields.fraction.name" 
            :label="loadTrainingDataset.fields.fraction.label"
            input-type="number" />
    </GroupElement>
</template>