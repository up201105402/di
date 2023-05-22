<script setup>
  import { ref, computed, reactive, markRaw } from 'vue';
  import stepTypes from '@/pipelines/steps'
  import CheckoutRepositoryForm from '@/pipelines/steps/components/CheckoutRepositoryForm.vue';
  import LoadTrainingDatasetForm from '@/pipelines/steps/components/LoadTrainingDatasetForm.vue';
  import TrainModelForm from '@/pipelines/steps/components/TrainModelForm.vue';
  import BaseButton from "@/components/BaseButton.vue";
  import BaseButtons from "@/components/BaseButtons.vue";
  import { watch } from 'vue';
  import { cleanObject } from '@/util';
  import { camel2title, customAxios } from '@/util'
  import useSteps from '@/pipelines/steps'

  let { steps, visitedSteps, activeStep, setStep, stepPlugin } = useSteps();

  const props = defineProps({
    nodeId: {
      type: String,
      required: false,
      default: null
    },
    nodeData: {
      type: Object,
      required: false,
      default: null,
    },
  });

  const data = computed({
    get: () => props.nodeData,
    set: (value) => { },
  });

  const emit = defineEmits(["onSubmit"]);

  const onSubmit = (e) => {
    cleanObject(e.data);
    emit("onSubmit", { id: props.nodeId, data: e.data })
  }

  const library = markRaw({
    SubmitButton: BaseButton
  });

  const formkitData = reactive({
    steps,
    visitedSteps,
    activeStep,
    plugins: [
      stepPlugin
    ],
    setStep: target => () => {
      setStep(target)
    },
    setActiveStep: stepName => () => {
      data.activeStep = stepName
    },
    showStepErrors: stepName => {
      return (steps[stepName].errorCount > 0 || steps[stepName].blockingCount > 0) && (visitedSteps.value && visitedSteps.value.includes(stepName))
    },
    stepIsValid: stepName => {
      return steps[stepName].valid && steps[stepName].errorCount === 0
    },
    submitApp: async (formData, node) => {
      try {
        const res = await customAxios.post(formData)
        node.clearErrors()
        alert('Your application was submitted successfully!')
      } catch (err) {
        node.setErrors(err.formErrors, err.fieldErrors)
      }
    },
    stringify: (value) => JSON.stringify(value, null, 2),
    camel2title
  })

  const schema = [
    {
      $cmp: 'FormKit',
      props: {
        type: 'form',
        id: 'form',
        onSubmit: '$submitApp',
        plugins: '$plugins',
        actions: false,
      },
      children: [
        {
          $el: 'ul',
          attrs: {
            class: "steps"
          },
          children: [
            {
              $el: 'li',
              for: ['step', 'stepName', '$steps'],
              attrs: {
                class: {
                  'step': true,
                  'has-errors': '$showStepErrors($stepName)'
                },
                onClick: '$setActiveStep($stepName)',
                'data-step-active': '$activeStep === $stepName',
                'data-step-valid': '$stepIsValid($stepName)'
              },
              children: [
                {
                  $el: 'span',
                  if: '$showStepErrors($stepName)',
                  attrs: {
                    class: 'step--errors'
                  },
                  children: '$step.errorCount + $step.blockingCount'
                },
                '$camel2title($stepName)'
              ]
            }
          ]
        },
        {
          $el: 'div',
          attrs: {
            class: 'form-body'
          },
          children: [
            {
              $el: 'section',
              attrs: {
                style: {
                  if: '$activeStep !== "nameAndType"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'nameAndType',
                  name: 'nameAndType',
                  children: [
                    {
                      $formkit: 'text',
                      name: 'nodeName',
                      label: 'Step Name',
                      placeholder: 'Step Name',
                      validation: 'required'
                    },
                    {
                      $formkit: 'select',
                      name: 'nodeType',
                      label: 'Node Type',
                      options: [
                        'Checkout Repo',
                        'Train Model'
                      ],
                      validation: 'required'
                    },
                  ]
                }
              ]
            },
            {
              $el: 'section',
              attrs: {
                style: {
                  if: '$activeStep !== "organizationInfo"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'organizationInfo',
                  name: 'organizationInfo',
                  children: [
                    {
                      $formkit: 'text',
                      label: '*Organization name',
                      name: 'org_name',
                      placeholder: 'MyOrg, Inc.',
                      help: 'Enter your official organization name.',
                      validation: 'required|length:3'
                    },
                    {
                      $formkit: 'date',
                      label: 'Date of incorporation',
                      name: 'date_inc',
                      validation: 'required'
                    }
                  ]
                }
              ]
            },
            {
              $el: 'div',
              attrs: {
                class: 'step-nav'
              },
              children: [
                {
                  $formkit: 'button',
                  disabled: '$activeStep === "nameAndType"',
                  onClick: '$setStep(-1)',
                  children: 'Previous Step'
                },
                {
                  $formkit: 'button',
                  disabled: '$activeStep === "application"',
                  onClick: '$setStep(1)',
                  children: 'Next Step'
                }
              ]
            },
          ]
        },
        {
          $formKit: 'submit',
          label: 'Submit',
          disabled: '$get(form).state.valid !== true',
          children: [
            {
              $cmp: 'SubmitButton',
              props: {
                label: 'Submit',
                disabled: '$get(form).state.valid !== true',
                color: 'success',
                onClick: '$submitApp'
              }
            }
          ],
        }
      ]
    },
  ]

</script>

<template>
  <!-- <Vueform v-model="data" sync :endpoint="false" @submit="onSubmit">
    <template #empty>
      <FormSteps>
        <FormStep name="step1" label="Select Type" :elements="['stepName', 'stepType']" />
        <FormStep name="step2" label="Configure" :elements="Object.values(stepTypes).map(stepType => stepType.name)" :labels="{ next: 'Submit' }" />
      </FormSteps>

      <FormElements>
        <TextElement name="stepName" label="Name" :rules="['required']" />
        <SelectElement name="stepType" :search="true" :native="false" label="Type" input-type="search"
          autocomplete="off" :items="Object.values(stepTypes).map(stepType => ({
            'value': stepType.name,
            'label': stepType.label,
          }))" :rules="['required']" />
        <CheckoutRepositoryForm />
        <LoadTrainingDatasetForm />
        <TrainModelForm />
      </FormElements>

      <FormStepsControls />
    </template>
  </Vueform> -->
  <FormKitSchema :library="library" :schema="schema" :data="formkitData" />
</template>

<style>
  @import "https://cdn.formk.it/web-assets/multistep-form.css";
</style>