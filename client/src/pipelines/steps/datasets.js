
export const datasetConfigFields = [
    {
        $formkit: 'text',
        name: "filePath",
        label: "File Path",
        validation: 'required',
    },
    {
        $formkit: 'number',
        name: "lowerXRangeIndex",
        label: "Lower X Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "upperXRangeIndex",
        label: "Upper X Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "lowerYRangeIndex",
        label: "Lower Y Range Index",
        min: 0,
    },
    {
        $formkit: 'number',
        name: "upperYRangeIndex",
        label: "Upper Y Range Index",
        min: 0,
    },
]

export const datasetConfigSection = {
    $formkit: 'group',
    id: 'stepConfig',
    name: 'stepConfig',
    children: datasetConfigFields,
    if: '$isActiveNodeType("ScikitTrainingDataset") || $isActiveNodeType("ScikitTestingDataset")',
}

export const scikitDatasets = [
    { id: 0, value: "scikitBreastCancer", label: "Breast Cancer Dataset" },
    { id: 1, value: "scikitDiabetes", label: "Diabetes Dataset" },
    { id: 2, value: "scikitDigits", label: "Digits Dataset" },
    { id: 3, value: "scikitIris", label: "Iris Dataset" },
    { id: 4, value: "scikitLinerrud", label: "Linnerud Dataset" },
    { id: 5, value: "scikitWine", label: "Wine Dataset" },
    { id: 6, value: "scikitLoadFile", label: "Load Daset From File" },
]

export const scikitDatasetSelect = [
    {
        $formkit: 'select',
        name: 'dataset',
        label: 'Dataset',
        placeholder: "",
        options: scikitDatasets,
        validation: 'required',
        if: '$isActiveNodeType("ScikitTrainingDataset") || $isActiveNodeType("ScikitTestingDataset")',
        onChange: "$setSciKitDataset",
    },
]