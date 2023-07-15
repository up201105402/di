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

export const watchActiveStep = (newStep, oldStep) => {
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
};