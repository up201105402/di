<script setup>

    import { ref, computed, watch, onMounted } from "vue";
    import $ from 'jquery';

    const props = defineProps({
        id: {
            type: Number,
            required: true,
        },
        imageURL: {
            type: String,
            required: true,
        }
    });

    const emit = defineEmits(["mouseClick"])

    const mouseClick = (e) => {
        const $canvas = $("#" + e.target.id);
        const ptX = e.pageX - $canvas.offset().left;
        const ptY = e.pageY - $canvas.offset().top;
        emit("mouseClick", props.id, ptX, ptY)
    }

    onMounted(() => {
        const canvas = document.getElementById(`query-${props.id}-canvas`);
        const ctx = canvas.getContext('2d');
        const imageObj = new Image();

        imageObj.onload = function (e) {
            ctx.drawImage(imageObj, 0, 0, imageObj.width, imageObj.height);
        };
        imageObj.src = props.imageURL;

        // ctx.globalAlpha = 0.5;

        canvas.addEventListener('click', mouseClick, false);
    })
</script>
<template>
    <canvas :id="`query-${id}-canvas`" class="mx-auto" width="300" height="300"></canvas>
</template>