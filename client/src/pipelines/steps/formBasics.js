export const cancelAndSubmitButtons = {
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
                '$camel2title($stepName)'
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
};
