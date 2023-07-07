
export const checkoutRepoStepConfigFields = [
    {
        $formkit: 'text',
        label: 'URL',
        name: 'repoURL',
        validation: 'required|url',
    },
    {
        $formkit: 'text',
        label: 'File Path',
        name: 'filePath',
        validation: 'required',
    },
]

export const checkoutRepoConfigSection = {
    $formkit: 'group',
    id: 'checkoutRepoConfig',
    name: 'stepConfig',
    children: checkoutRepoStepConfigFields,
}