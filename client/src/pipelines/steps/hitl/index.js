import { reactive, toRef, ref, watch } from 'vue';
import { i18nFromStepName } from '@/util';
import { getNode, createMessage } from '@formkit/core';
import { stepTabs, cancelAndSubmitButtons, getFormBody } from '@/pipelines/steps/formBasics';
import { i18n } from '@/i18n';

const { t } = i18n.global;

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
        label: t('pages.pipelines.edit.dialog.stepConfig.data_dir'),
        validation: 'required|dirPath',
    },
    {
        $formkit: 'text',
        name: 'models_dir',
        label: t('pages.pipelines.edit.dialog.stepConfig.models_dir'),
        validation: 'required|dirPath',
    },
    {
        $formkit: 'number',
        name: 'epochs',
        label: t('pages.pipelines.edit.dialog.stepConfig.epochs'),
        validation: 'required|min:0',
    },
    {
        $formkit: 'number',
        name: 'tr_fraction',
        label: t('pages.pipelines.edit.dialog.stepConfig.tr_fraction'),
        validation: 'required',
        step: 'any',
    },
    {
        $formkit: 'number',
        name: 'val_fraction',
        label: t('pages.pipelines.edit.dialog.stepConfig.val_fraction'),
        validation: 'required',
        step: 'any',
    },
    {
        $formkit: 'text',
        name: 'train_desc',
        label: t('pages.pipelines.edit.dialog.stepConfig.train_desc'),
        validation: 'required',
    },
    {
        $formkit: 'select',
        name: 'sampling',
        label: t('pages.pipelines.edit.dialog.stepConfig.sampling'),
        validation: 'required',
        options: ['low_entropy', 'high_entropy'],
    },
    {
        $formkit: 'number',
        name: 'entropy_thresh',
        label: t('pages.pipelines.edit.dialog.stepConfig.entropy_thresh'),
        validation: 'required',
        step: 'any',
    },
    {
        $formkit: 'number',
        name: 'nr_queries',
        label: t('pages.pipelines.edit.dialog.stepConfig.nr_queries'),
        validation: 'required|min:0',
    },
    {
        $formkit: 'checkbox',
        name: 'isOversampled',
        label: t('pages.pipelines.edit.dialog.stepConfig.isOversampled'),
    },
    {
        $formkit: 'number',
        name: 'start_epoch',
        label: t('pages.pipelines.edit.dialog.stepConfig.start_epoch'),
        validation: 'required|min:0',
    },
    {
        $formkit: 'select',
        name: 'dataset',
        label: t('pages.pipelines.edit.dialog.stepConfig.dataset'),
        validation: 'required',
        options: ['APTOS19', 'ISIC17','NCI','ROSEYoutu','PornographyXXX','Custom'],
    },
    {
        $formkit: 'select',
        name: 'pretrained_model',
        label: t('pages.pipelines.edit.dialog.stepConfig.pretrainedModel'),
        validation: 'required',
        options: [
            'alexnet', 
            'convnext_tiny', 
            'convnext_small', 
            'convnext_base', 
            'convnext_large',
            'densenet121',
            'densenet161',
            'densenet169',
            'densenet201', 
            'efficientnet_b0',
            'efficientnet_b1',
            'efficientnet_b2',
            'efficientnet_b3',
            'efficientnet_b4',
            'efficientnet_b5',
            'efficientnet_b6',
            'efficientnet_b7',
            'efficientnet_v2_s',
            'efficientnet_v2_m',
            'efficientnet_v2_l',
            'googlenet',
            'inception_v3',
            'mnasnet0_5',
            'mnasnet0_75',
            'mnasnet1_0',
            'mnasnet1_3',
            'mobilenet_v2',
            'mobilenet_v3_large',
            'mobilenet_v3_small',
            'regnet_y_400mf',
            'regnet_y_800mf',
            'regnet_y_1_6gf',
            'regnet_y_3_2gf',
            'regnet_y_8gf',
            'regnet_y_16gf',
            'regnet_y_32gf',
            'regnet_y_128gf',
            'regnet_x_400mf',
            'regnet_x_800mf',
            'regnet_x_1_6gf',
            'regnet_x_3_2gf',
            'regnet_x_8gf',
            'regnet_x_16gf',
            'regnet_x_32gf',
            'resnet18',
            'resnet34',
            'resnet50',
            'resnet101',
            'resnet152',
            'resnext50_32x4d',
            'resnext101_32x8d',
            'resnext101_64x4d',
            'shufflenet_v2_x0_5',
            'shufflenet_v2_x1_0',
            'shufflenet_v2_x1_5',
            'shufflenet_v2_x2_0',
            'squeezenet1_0',
            'squeezenet1_1',
            'vgg11',
            'vgg11_bn',
            'vgg13',
            'vgg13_bn',
            'vgg16',
            'vgg16_bn',
            'vgg19',
            'vgg19_bn',
            'vit_b_16',
            'vit_b_32',
            'vit_l_16',
            'vit_l_32',
            'vit_h_14',
            'swin_t',
            'swin_s',
            'swin_b',
            'swin_v2_t',
            'swin_v2_s',
            'swin_v2_b',
            'maxvit_t',
            'wide_resnet50_2',
            'wide_resnet101_2',
        ],
    },
    {
        $formkit: 'select',
        name: 'optimizer',
        label: t('pages.pipelines.edit.dialog.stepConfig.optimizer'),
        validation: 'required',
        options: [
            'Adadelta',
            'Adagrad',
            'Adam',
            'AdamW',
            'SparseAdam',
            'Adamax',
            'ASGD',
            'SGD',
            'RAdam',
            'Rprop',
            'RMSprop',
            'NAdam',
            'LBFGS',
        ],
    },
    {
        $formkit: 'number',
        name: 'learning_rate',
        label: t('pages.pipelines.edit.dialog.stepConfig.learningRate'),
        validation: 'required',
        value: 0.00001,
    },
];

export const hitlForm = function (data, onSubmit, editable = true) {
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
            }
        },
        stringify: (value) => JSON.stringify(value, null, 2),
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