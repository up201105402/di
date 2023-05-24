<script setup>
  import { ref, computed, reactive, markRaw } from 'vue';
  import stepTypes from '@/pipelines/steps'
  import CheckoutRepositoryForm from '@/pipelines/steps/components/CheckoutRepositoryForm.vue';
  import LoadTrainingDatasetForm from '@/pipelines/steps/components/LoadTrainingDatasetForm.vue';
  import TrainModelForm from '@/pipelines/steps/components/TrainModelForm.vue';
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
    data: {
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
      formkitData.activeStep = stepName
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
        emit("onSubmit", { id: props.nodeId, data: formData })
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
        value: { ...props.data }
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
                      options: {
                        checkoutRepo: 'Checkout Repo',
                        trainModel: 'Train Model'
                      },
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
                  if: '$activeStep !== "stepConfig"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'stepConfig',
                  name: 'stepConfig',
                  children: [
                    {
                      $formkit: 'date',
                      label: 'Date of incorporation',
                      name: 'date_inc',
                      validation: 'required',
                      attrs: {
                        style: {
                          if: '$nameAndType !== "stepConfig"',
                          then: 'display: none;'
                        }
                      },
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
                  disabled: '$activeStep === "stepConfig"',
                  onClick: '$setStep(1)',
                  children: 'Next Step'
                }
              ]
            },
          ]
        },
        {
          $formkit: 'submit',
          label: 'Submit',
          disabled: '$get(form).state.valid !== true',
        }
      ]
    },
  ];

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
  <FormKitSchema :schema="schema" :data="formkitData" />
</template>

<style>
  @import "https://cdn.formk.it/web-assets/multistep-form.css";
</style>