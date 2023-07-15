import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, customDelay } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { watchActiveStep, stepNav, cancelAndSubmitButtons } from '@/pipelines/steps/formBasics';

export const scriptStepConfigFields = (pipelineID) => {
    return [
        {
            $cmp: 'ScriptEditor',
            name: 'script',
            label: 'Script',
            bind: '$editorBindingProps',
            if: '$scriptType == "inline"'
        },
        {
            $el: 'label',
            attrs: {
                class: 'formkit-label',
                for: 'script-file-upload',
            },
            children: 'File',
            if: '$scriptType == "file"'
        },
        {
            $cmp: 'FormFilePicker',
            name: 'file',
            label: 'File Upload',
            bind: '$filePickerProps',
            props: {
                id: 'script-file-upload',
                label: 'File Upload',
                url: `/api/pipeline/${pipelineID}/file`
            },
            if: '$scriptType == "file"'
        },
    ]
}

export const scriptConfigSection = (pipelineID) => {
    return {
        $formkit: 'group',
        id: 'stepConfig',
        name: 'stepConfig',
        children: scriptStepConfigFields(pipelineID),
    }
}

export const scriptForm = function (data, onSubmit) {
    const activeStep = ref('');
    const steps = reactive({});
    const visitedSteps = ref([]); // track visited steps
    const showIsFirstStep = data.nameAndType == null || data?.nameAndType?.isFirstStep == false;
    const editorBindingProps = {
        modelValue: data?.stepConfig?.script,
        onModelValueUpdate: (e) => getNode(activeStep.value).value.script = e
    }
    const filePickerProps = {
        filename: data?.stepConfig?.filename,
        onFileUpdated: (e) => getNode(activeStep.value).value.filename = e
    }
    const scriptType = ref(data?.nameAndType?.scriptType || 'inline')

    // watch our activeStep and store visited steps
    // to know when to show errors
    watch(activeStep, watchActiveStep)

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
        scriptType,
        editorBindingProps,
        filePickerProps,
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
        setActiveType: event => {
            scriptType.value = event.target.value;
        },
        submitForm: async (formData, node) => {
            try {
                if (formData.stepConfig.script) {
                    formData.stepConfig.script = formData.stepConfig.script.replaceAll('<p>', '').replaceAll('</p>', '\n')   
                }
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
                                        {
                                            $formkit: 'select',
                                            name: 'scriptType',
                                            label: 'Script Type',
                                            placeholder: "",
                                            options: [
                                                {
                                                    value: 'inline',
                                                    label: 'Inline'
                                                }, 
                                                {
                                                    value: 'file',
                                                    label: 'File'
                                                }],
                                            validation: 'required',
                                            onChange: "$setActiveType",
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
                                scriptConfigSection(data.pipelineID)
                            ]
                        },
                        stepNav,
                    ]
                },
                cancelAndSubmitButtons
            ]
        },
    ];

    return { formSchema, formkitData }
}