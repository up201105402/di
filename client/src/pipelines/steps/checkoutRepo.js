import { reactive, toRef, ref, watch, isReadonly } from 'vue';
import { camel2title, i18nFromStepName } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { stepTabs, cancelAndSubmitButtons, getFormBody } from '@/pipelines/steps/formBasics';
import { i18n } from '@/i18n';

const { t } = i18n.global;

export const nameAndTypeGroupChildren = [
    {
        $formkit: 'text',
        name: 'name',
        label: t('pages.pipelines.edit.dialog.nameAndType.name'),
        placeholder: t('pages.pipelines.edit.dialog.nameAndType.name'),
        validation: 'required'
    },
    {
        $formkit: 'checkbox',
        name: 'isFirstStep',
        label: t('pages.pipelines.edit.dialog.nameAndType.isFirstStep'),
        if: '$showIsFirstStep == true',
    },
]

export const stepConfigGroupChildren = [
    {
        $formkit: 'text',
        label: t('pages.pipelines.edit.dialog.stepConfig.repoUrl'),
        name: 'repoURL',
        validation: 'required|url',
    },
]

export const checkoutRepoForm = function (data, onSubmit, editable = true) {
    const activeStep = ref('');
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
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
                console.error(err);
            }
        },
        stringify: (value) => JSON.stringify(value, null, 2),
        camel2title,
        i18nFromStepName
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
                stepTabs,
                getFormBody(nameAndTypeGroupChildren, stepConfigGroupChildren),
                cancelAndSubmitButtons(editable),
            ]
        },
    ];

    return { formSchema, formkitData }
}