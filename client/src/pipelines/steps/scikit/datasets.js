import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, customDelay } from '@/util';
import { getNode, createMessage } from '@formkit/core';

export const datasetConfigFields = [
    {
        $formkit: 'text',
        name: "filePath",
        label: "File Path",
        validation: 'required',
        if: '$isActiveDataset("scikitLoadFile")'
    },
    {
        $formkit: 'number',
        name: "lowerXRangeIndex",
        label: "Lower X Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "upperXRangeIndex",
        label: "Upper X Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "lowerYRangeIndex",
        label: "Lower Y Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "upperYRangeIndex",
        label: "Upper Y Range Index",
        min: 0,
    },
]

export const datasetConfigSection = {
    $formkit: 'group',
    id: 'stepConfig',
    name: 'stepConfig',
    children: datasetConfigFields,
}

export const scikitDatasetOptions = [
    { id: 0, value: "scikitBreastCancer", label: "Breast Cancer Dataset" },
    { id: 1, value: "scikitDiabetes", label: "Diabetes Dataset" },
    { id: 2, value: "scikitDigits", label: "Digits Dataset" },
    { id: 3, value: "scikitIris", label: "Iris Dataset" },
    { id: 4, value: "scikitLinerrud", label: "Linnerud Dataset" },
    { id: 5, value: "scikitWine", label: "Wine Dataset" },
    { id: 6, value: "scikitLoadFile", label: "Load Daset From File" },
]

export const scikitDatasetForm = function (data, onSubmit) {
    const activeStep = ref('');
    const activeDataset = ref(scikitDatasetOptions[0].value);
    const steps = reactive({});
    const visitedSteps = ref([]); // track visited steps
    const showIsFirstStep = data.nameAndType == null || data?.nameAndType?.isFirstStep == false;

    // watch our activeStep and store visited steps
    // to know when to show errors
    watch(activeStep, (newStep, oldStep) => {
        if (oldStep && !visitedSteps.value.includes(oldStep)) {
            visitedSteps.value.push(oldStep)
        }
        // trigger showing validation on fields
        // within all visited steps
        visitedSteps.value.forEach((step) => {
            const node = getNode(step)
            node.walk((n) => {
                n.store.set(
                    createMessage({
                        key: 'submitted',
                        value: true,
                        visible: false,
                    })
                )
            })
        })
    })

    const setStep = (delta) => {
        const stepNames = Object.keys(steps);
        const currentIndex = stepNames.indexOf(activeStep.value);
        activeStep.value = stepNames[currentIndex + delta];
    }

    // pushes the steps (group nodes - $formkit: 'group') into the steps object
    const stepPlugin = (node) => {
        if (node.props.type == "group") {
            // builds an object of the top-level groups
            steps[node.name] = steps[node.name] || {}

            node.on('created', () => {
                // use 'on created' to ensure context object is available
                steps[node.name].valid = toRef(node.context.state, 'valid')
            })

            // listen for changes in error count and store it
            node.on('count:errors', ({ payload: count }) => {
                steps[node.name].errorCount = count
            })

            // listen for changes in count of blocking validations messages
            node.on('count:blocking', ({ payload: count }) => {
                steps[node.name].blockingCount = count
            })

            // set the active tab to the 1st tab
            if (activeStep.value === '') {
                activeStep.value = node.name
            }

            // Stop plugin inheritance to descendant nodes
            return false
        }
    }

    const formkitData = reactive({
        steps,
        visitedSteps,
        activeStep,
        showIsFirstStep,
        activeDataset,
        plugins: [
            stepPlugin
        ],
        isActiveDataset: (dataset) => {
            return activeDataset.value == dataset;
        },
        setActiveDataset: (changeEvent) => {
            activeDataset.value = changeEvent.target.value;
        },
        setStep: target => () => {
            setStep(target)
        },
        setActiveStep: stepName => () => {
            activeStep.value = stepName
        },
        showStepErrors: stepName => {
            return (steps[stepName].errorCount > 0 || steps[stepName].blockingCount > 0) && (visitedSteps.value && visitedSteps.value.includes(stepName))
        },
        stepIsValid: stepName => {
            return steps[stepName].valid && steps[stepName].errorCount === 0
        },
        submitForm: async (formData, node) => {
            try {
                await customDelay(formData);
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
            }
        },
        stringify: (value) => JSON.stringify(value, null, 2),
        camel2title
    })

    const formSchema = [
        {
            $cmp: 'FormKit',
            props: {
                type: 'form',
                id: 'form',
                onSubmit: '$submitForm',
                plugins: '$plugins',
                actions: false,
                value: { ...data }
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
                                style: {
                                    if: '$activeNodeType == ""',
                                    then: 'display: none;'
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
                                            name: 'name',
                                            label: 'Step Name',
                                            placeholder: 'Step Name',
                                            validation: 'required'
                                        },
                                        {
                                          $formkit: 'select',
                                          name: 'dataset',
                                          label: 'Dataset',
                                          placeholder: "",
                                          options: scikitDatasetOptions,
                                          validation: 'required',
                                          onChange: "$setActiveDataset",
                                        },
                                        {
                                            $formkit: 'checkbox',
                                            name: 'isFirstStep',
                                            label: 'Is First Step?',
                                            if: '$showIsFirstStep == true',
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
                                datasetConfigSection
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
                                    children: 'Back'
                                },
                                {
                                    $formkit: 'button',
                                    disabled: '$activeStep === "stepConfig"',
                                    onClick: '$setStep(1)',
                                    children: 'Next'
                                }
                            ]
                        },
                    ]
                },
                {
                    $el: 'div',
                    attrs: {
                        class: 'formkit-bottom-buttons'
                    },
                    children: [
                        {
                            $formkit: 'button',
                            label: 'Cancel',
                            id: 'cancel-create-step-button'
                        },
                        {
                            $formkit: 'submit',
                            label: 'Submit',
                            disabled: '$get(form).state.valid !== true',
                        },
                    ]
                },
            ]
        },
    ];

    return { formSchema, formkitData }
}