import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, customDelay } from '@/util';
import { getNode, createMessage } from '@formkit/core';

export const scriptStepConfigFields = [
    {
        $cmp: 'ScriptEditor',
        name: 'script',
        label: 'Script',
        bind: '$editorBindingProps',
    },
]

export const scriptConfigSection = {
    $formkit: 'group',
    id: 'stepConfig',
    name: 'stepConfig',
    children: scriptStepConfigFields,
}

export const scriptForm = function (data, onSubmit) {
    const activeStep = ref('');
    const steps = reactive({});
    const visitedSteps = ref([]); // track visited steps
    const showIsFirstStep = data.nameAndType == null || data?.nameAndType?.isFirstStep == false;
    const editorBindingProps = {
        modelValue: '',
        onModelValueUpdate: (e) => getNode(activeStep.value).value.script = e
    }

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
        editorBindingProps,
        plugins: [
            stepPlugin
        ],
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
                formData.stepConfig.script = formData.stepConfig.script.replaceAll('<p>', '').replaceAll('</p>', '\n')
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
                console.error(err);
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
                                scriptConfigSection
                            ]
                        },
                        {
                            $el: 'div',
                            attrs: {
                                class: 'step-nav'
                            },
                            children: [
                                {
                                    $cmp: 'BaseButton',
                                    onClick: '$setStep(-1)',
                                    children: 'Back',
                                    props: {
                                        label: 'Back',
                                        color: 'info',
                                        onClick: '$setStep(-1)',
                                        disabled: '$activeStep === "nameAndType"',
                                    }
                                },
                                {
                                    $cmp: 'BaseButton',
                                    onClick: '$setStep(1)',
                                    children: 'Next',
                                    props: {
                                        label: 'Next',
                                        color: 'info',
                                        onClick: '$setStep(1)',
                                        disabled: '$activeStep === "stepConfig"',
                                    }
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
                            $cmp: 'BaseCancelAndSubmitButtons',
                            props: {
                                cancelButtonId: 'cancel-create-step-button',
                                submitEnabled: '$get(form).state.valid == true'
                            }
                        }
                    ]
                },
            ]
        },
    ];

    return { formSchema, formkitData }
}