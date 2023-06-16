<script setup>
    import { TrainModel } from '@/pipelines/steps';
</script>

<template>
    <GroupElement 
        :name="TrainModel.name" 
        :label="TrainModel.label" 
        :conditions="TrainModel.conditions">
        <TextElement 
            :name="TrainModel.fields.modelDirectory.name" 
            :label="TrainModel.fields.modelDirectory.label" 
            :rules="TrainModel.fields.modelDirectory.rules" />
        <TextElement 
            :name="TrainModel.fields.epochs.name" 
            :label="TrainModel.fields.epochs.label"
            input-type="number" 
            :rules="TrainModel.fields.epochs.rules" 
            autocomplete="off" />
    </GroupElement>
</template>