<script setup>
    import { ref } from "vue";
    import FeedbackRectsTable from '@/components/FeedbackRectsTable.vue';
    import FeedbackCanvas from '@/components/FeedbackCanvas.vue';

    const props = defineProps({
        query: {
            type: Object,
            required: true
        }
    });

    const emit = defineEmits(['checked']);

    const getSelectedRects = (rects) => {
        return rects.filter(rect => rect.Selected);
    }

    const rectsSelected = ref(getSelectedRects(props.query.HumanFeedbackRects));

    const onImageClicked = (queryID, ptX, ptY) => {
        const rect = props.query.HumanFeedbackRects.find(rect => {
            return (ptX >= rect.X1) && (ptX <= rect.X2) && (ptY >= rect.Y1) && (ptY <= rect.Y2);
        });

        if (rect) {
            const selectedIndex = rectsSelected.value.findIndex(selectedRect => selectedRect.ID == rect.ID);

            if (selectedIndex > -1) {
                rectsSelected.value.splice(selectedIndex, 1);
            } else {
                rectsSelected.value.push(rect);
            }

            emit('checked', props.query.ID, rectsSelected.value);
        }
    }

    const onRectChecked = (queryID, rect) => {
        const selectedIndex = rectsSelected.value.findIndex(selectedRect => selectedRect.ID == rect.ID);

        if (selectedIndex > -1) {
            rectsSelected.value.splice(selectedIndex, 1);
        } else {
            rectsSelected.value.push(rect);
        }

        emit('checked', props.query.ID, rectsSelected.value);
    }

</script>

<template>
    <div class="border-1 surface-border border-round m-2 text-center py-5 px-3">
        <FeedbackCanvas :id="query.HumanFeedbackQuery.ID" :imageURL="query.ImageURL" :drawnRectangles="rectsSelected" @mouseClick="onImageClicked" />
        <div>
            <h4 class="mb-1">Query {{ query.HumanFeedbackQuery.QueryID + 1 }}</h4>
            <h4 class="mb-1">{{ query.RunStepStatus.Name }} ({{ query.RunStepStatus.ID }})</h4>
        </div>
        <FeedbackRectsTable :items="query.HumanFeedbackRects" :queryID="query.HumanFeedbackQuery.QueryID" :checkedRows="rectsSelected" @checked="onRectChecked" checkable />
    </div>
</template>