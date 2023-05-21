<script setup>
  import { ref, computed, reactive } from 'vue';
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
        const res = await customaxios.post(formData)
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
                  if: '$activeStep !== "contactInfo"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'contactInfo',
                  name: 'contactInfo',
                  children: [
                    {
                      $formkit: 'text',
                      name: 'full_name',
                      label: '*Full Name',
                      placeholder: 'First Last',
                      validation: 'required'
                    },
                    {
                      $formkit: 'email',
                      name: 'email',
                      label: '*Email address',
                      placeholder: 'email@domain.com',
                      validation: 'required|email'
                    },
                    {
                      $formkit: 'tel',
                      name: 'tel',
                      label: '*Telephone',
                      placeholder: 'xxx-xxx-xxxx',
                      help: 'Phone number must be in the xxx-xxx-xxxx format.',
                      validation: 'required|matches:/^[0-9]{3}-[0-9]{3}-[0-9]{4}$/'
                    }
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
              $el: 'section',
              attrs: {
                style: {
                  if: '$activeStep !== "application"',
                  then: 'display: none;'
                }
              },
              children: [
                {
                  $formkit: 'group',
                  id: 'application',
                  name: 'application',
                  children: [
                    {
                      $formkit: 'checkbox',
                      label: '*I\'m not a previous grant recipient',
                      help: 'Have you received a grant from us before?',
                      name: 'not_previous_recipient',
                      validation: 'required|accepted',
                      validationMessages: {
                        accepted: 'We can only give one grant per organization.'
                      }
                    },
                    {
                      $formkit: 'textarea',
                      label: '*How will you use the money?',
                      name: 'how_money',
                      help: 'Must be between 20 and 500 characters.',
                      placeholder: 'Describe how the grant will accelerate your efforts.',
                      validation: 'required|length:20,500'
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
                  disabled: '$activeStep === "contactInfo"',
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
            {
              $el: 'details',
              children: [
                {
                  $el: 'summary',
                  children: 'Form data'
                },
                {
                  $el: 'pre',
                  children: '$stringify( $get(form).value )'
                }
              ]
            },
          ]
        },
        {
          $formkit: 'submit',
          label: 'Submit Application',
          disabled: '$get(form).state.valid !== true'
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
  <FormKitSchema :schema="schema" :data="formkitData" />
</template>

<style>
  @import "https://cdn.formk.it/web-assets/multistep-form.css";
</style>