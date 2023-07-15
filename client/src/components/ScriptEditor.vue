<script setup>
    import { ref, computed, onMounted } from 'vue';
    import Editor from 'primevue/editor';

    const props = defineProps(['modelValue'])
    const emit = defineEmits(['modelValueUpdate']) // this is need instead of update:modelValue because FormKit does not seem to support v-model

    const quillModules = {
        "toolbar": false
    }

    const text = ref(props.modelValue);

    const onModelValueUpdate = (e) => {
        emit('modelValueUpdate', e);
    }

</script>

<template>
    <div>
        <label class="formkit-label" for="script-editor">Script</label>
        <Editor id="formkit-label" v-model="text" @update:modelValue="onModelValueUpdate" editorStyle="height: 320px; margin-bottom:10px" :modules="quillModules" />
    </div>
</template>