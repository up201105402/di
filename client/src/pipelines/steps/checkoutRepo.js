
export const checkoutRepoStepConfigFields = [
    {
        $formkit: 'text',
        label: 'URL',
        name: 'repoURL',
        validation: 'required|url',
        if: '$isActiveNodeType("CheckoutRepo")',
    },
    {
        $formkit: 'text',
        label: 'File Path',
        name: 'filePath',
        validation: 'required',
        if: '$isActiveNodeType("CheckoutRepo")',
    },
]