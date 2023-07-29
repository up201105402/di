<script setup>
    import { ref, computed, onMounted } from 'vue';
    import Editor from 'primevue/editor';

    const props = defineProps(['modelValue', 'class'])
    const emit = defineEmits(['modelValueUpdate']) // this is need instead of update:modelValue because FormKit does not seem to support v-model

    const quillModules = {
        "toolbar": false,
    }

    const html = props.modelValue ? '<p>' + props.modelValue.replaceAll('\n', '</p><p>') + '</p>' : '';
    const text = ref(html.replace('<p></p>', '<p>&nbsp;</p>'), '');

    const onModelValueUpdate = (e) => {
        // emit('modelValueUpdate', e);
    }

    const onTextChange = (e) => {
        emit('modelValueUpdate', e);
    }

</script>

<template>
    <div>
        <label class="formkit-label" for="script-editor">Script</label>
        <Editor id="formkit-label" :pt="{ content: { class: props.class } }" v-model="text" @text-change="onTextChange" @update:modelValue="onModelValueUpdate" editorStyle="height: 320px; margin-bottom:10px" :modules="quillModules" />
    </div>
</template>

<style>
    .ql-editor {
        tab-size: 20;
        -moz-tab-size: 20;
        -o-tab-size: 20;
    }
</style>