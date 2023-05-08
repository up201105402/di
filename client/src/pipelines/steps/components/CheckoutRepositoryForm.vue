<template>
    <GroupElement name="Checkout repository" :conditions="[
        [
            'Step Type',
            'in',
            [
                '0',
            ],
        ],
    ]" label="Checkout repository">
        <TextElement name="Repository URL" label="URL" :rules="[
            'required',
            'regex:((git|ssh|http(s)?)|(git@[\\w\\.]+))(:(//)?)([\\w\\.@\\:/\\-~]+)(\\.git)(/)?',
        ]" />
    </GroupElement>
</template>