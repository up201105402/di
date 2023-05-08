<script setup>
import { stepTypes } from '@/pipelines/steps'
const CheckoutRepositoryForm = stepTypes[0].definedComponent;
</script>

<template>
  <!-- <Vueform>
        <template #empty>
            <FormSteps>
                <FormStep name="page0" label="Select Type" :elements="[
                'Step Name',
                'Step Type',
              ]" />
                <FormStep name="page1" label="Configure" :elements="stepTypes.map(step => step.name)" :labels="{
                next: 'Add',
              }" />
            </FormSteps>

            <FormElements>
                <TextElement name="Step Name" label="Name" :rules="['required']" />
                <SelectElement name="Step Type" :search="true" :native="false" label="Type" input-type="search"
                    autocomplete="off" :items="stepTypes.map(stepType => ({
                        'value': stepType.id,
                        'label': stepType.name,
                    }))" 
                    :rules="['required']" />
                <GroupElement name="Checkout repository" :conditions="[
                [
                  'Step Type',
                  'in',
                  [
                    '0',
                  ],
                ],
              ]" label="Checkout repository">
                    <TextElement name="Repository URL" label="URL" :rules="[
                  'required',
                  'regex:((git|ssh|http(s)?)|(git@[\\w\\.]+))(:(//)?)([\\w\\.@\\:/\\-~]+)(\\.git)(/)?',
                ]" />
                </GroupElement>
                <CheckoutRepositoryForm />
                <GroupElement name="Load Training Dataset" :conditions="[
                [
                  'Step  Type',
                  'in',
                  [
                    '1',
                  ],
                ],
              ]" label="Load Training Dataset">
                    <TextElement name="Training Dataset Directory" label="Directory" :rules="[
                  'required',
                  'regex:^(.+)\\/([^\\/]+)$',
                ]" />
                    <TextElement name="Fraction of the training data" label="Fraction of the training data to use"
                        :rules="[
                  'required',
                  'regex:^(.+)\\/([^\\/]+)$',
                ]" input-type="number" />
                </GroupElement>
                <GroupElement name="Train Model" label="Train Model" :conditions="[
                [
                  'Step  Type',
                  'in',
                  [
                    '2',
                  ],
                ],
              ]">
                    <TextElement name="Model Directory" label="Directory" :rules="[
                  'required',
                  'regex:^(.+)\\/([^\\/]+)$',
                ]" />
                    <TextElement name="Number of Epochs" input-type="number" :rules="[
                  'nullable',
                  'required',
                  'min:1',
                  'numeric',
                ]" autocomplete="off" label="Epochs" />
                </GroupElement>
            </FormElements>

            <FormStepsControls />
        </template>
    </Vueform> -->
  <FormKit type="multi-step">
    <FormKit type="step" name="stepOne">
      <!-- content for stepOne goes here! -->
    </FormKit>
  </FormKit>
</template>