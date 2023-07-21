<script setup>

import { watch, onMounted } from "vue";
import $ from 'jquery';

const props = defineProps({
    id: {
        type: Number,
        required: true,
    },
    imageURL: {
        type: String,
        required: true,
    },
    drawnRectangles: {
        type: Array,
        default: [],
    }
});

const emit = defineEmits(["mouseClick"])

const imageObj = new Image();

const mouseClick = (e) => {
    const $canvas = $("#" + e.target.id);
    const ptX = e.pageX - $canvas.offset().left;
    const ptY = e.pageY - $canvas.offset().top;
    emit("mouseClick", props.id, ptX, ptY)
}

const drawRectangles = (ctx, rects) => {
    rects.forEach(rect => {
        ctx.beginPath();
        ctx.lineWidth = "3";
        ctx.strokeStyle = "green";
        ctx.rect(rect.X1, rect.Y1, rect.X2 - rect.X1, rect.Y2 - rect.Y1);
        ctx.stroke();
    })
}

const drawCanvas = () => {
    const canvas = document.getElementById(`query-${props.id}-canvas`);
    
    const ctx = canvas.getContext('2d');
    ctx.reset();

    imageObj.onload = () => {
        ctx.drawImage(imageObj, 0, 0, imageObj.width, imageObj.height);
        drawRectangles(ctx, props.drawnRectangles)
    };

    if (!imageObj.src) {
        imageObj.src = props.imageURL;
    } else {
        ctx.drawImage(imageObj, 0, 0, imageObj.width, imageObj.height);
        drawRectangles(ctx, props.drawnRectangles)
    }

    canvas.addEventListener('click', mouseClick, false);
}

watch(props.drawnRectangles, drawCanvas)

onMounted(drawCanvas)

</script>
<template>
    <canvas :id="`query-${id}-canvas`" class="mx-auto" width="300" height="300"></canvas>
</template>