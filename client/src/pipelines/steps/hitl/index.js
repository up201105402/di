import { reactive, toRef, ref, watch } from 'vue';
import { camel2title, i18nFromStepName } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { stepTabs, cancelAndSubmitButtons, getFormBody } from '@/pipelines/steps/formBasics';
import { i18n } from '@/i18n';

const { t } = i18n.global;

export const datesetOptions = [
    { id: 0, value: "scikitBreastCancer", label: t('pages.pipelines.edit.dialog.nameAndType.dataset.options.scikitBreastCancer') },
    { id: 1, value: "scikitDiabetes", label: t('pages.pipelines.edit.dialog.nameAndType.dataset.options.scikitDiabetes') },
    { id: 2, value: "scikitDigits", label: t('pages.pipelines.edit.dialog.nameAndType.dataset.options.scikitDigits') },
];

const nameAndTypeGroupChildren = [
    {
        $formkit: 'text',
        name: 'name',
        label: t('pages.pipelines.edit.dialog.nameAndType.name'),
        placeholder: t('pages.pipelines.edit.dialog.nameAndType.name'),
        validation: 'required'
    },
];

export const stepConfigGroupChildren = [
    {
        $formkit: 'text',
        name: 'data_dir',
        label: 'pages.pipelines.edit.dialog.stepConfig.data_dir',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'models_dir',
        label: 'pages.pipelines.edit.dialog.stepConfig.models_dir',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'epochs',
        label: 'pages.pipelines.edit.dialog.stepConfig.epochs',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'tr_fraction',
        label: 'pages.pipelines.edit.dialog.stepConfig.tr_fraction',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'val_fraction',
        label: 'pages.pipelines.edit.dialog.stepConfig.val_fraction',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'train_desc',
        label: 'pages.pipelines.edit.dialog.stepConfig.train_desc',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'sampling',
        label: 'pages.pipelines.edit.dialog.stepConfig.sampling',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'entropy_thresh',
        label: 'pages.pipelines.edit.dialog.stepConfig.entropy_thresh',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'nr_queries',
        label: 'pages.pipelines.edit.dialog.stepConfig.nr_queries',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'isOversampled',
        label: 'pages.pipelines.edit.dialog.stepConfig.isOversampled',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'start_epoch',
        label: 'pages.pipelines.edit.dialog.stepConfig.start_epoch',
        validation: 'required'
    },
    {
        $formkit: 'text',
        name: 'dataset',
        label: 'pages.pipelines.edit.dialog.stepConfig.dataset',
        validation: 'required'
    }
]

export const hitlForm = function (data, onSubmit, editable = true) {
    const activeStep = ref('');
    const activeDataset = ref(data?.nameAndType?.dataset || datesetOptions[0].value);
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
                node.clearErrors()
                onSubmit(formData);
            } catch (err) {
                node.setErrors(err.formErrors, err.fieldErrors)
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