import { i18n } from "@/i18n";

const { t } = i18n.global;

export const cancelAndSubmitButtons = (isSubmitEnabled = true) => {
    return {
        $el: 'div',
        attrs: {
            class: 'formkit-bottom-buttons'
        },
        children: [
            {
                $cmp: 'BaseCancelAndSubmitButtons',
                props: {
                    cancelButtonId: 'cancel-create-step-button',
                    submitEnabled: isSubmitEnabled ? '$get(form).state.valid == true' : false,
                }
            }
        ]
    }
};

export const stepTabs = {
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
                '$i18nFromStepName($stepName)'
            ]
        }
    ]
};

export const getFormBody = (nameAndTypeSectionChildren, stepConfigSectionChildren) => {
    return {
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
                        children: nameAndTypeSectionChildren
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
                        children: stepConfigSectionChildren,
                    }
                ]
            },
            stepNav,
        ]
    }
};

export const stepNav = {
    $el: 'div',
    attrs: {
        class: 'step-nav'
    },
    children: [
        {
            $cmp: 'BaseButton',
            onClick: '$setStep(-1)',
            children: t('buttons.back'),
            props: {
                label: t('buttons.back'),
                color: 'info',
                onClick: '$setStep(-1)',
                disabled: '$activeStep === "nameAndType"',
            }
        },
        {
            $cmp: 'BaseButton',
            onClick: '$setStep(1)',
            children: t('buttons.next'),
            props: {
                label: t('buttons.next'),
                color: 'info',
                onClick: '$setStep(1)',
                disabled: '$activeStep === "stepConfig"',
            }
        }
    ]
};
