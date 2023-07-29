<script setup>
    import { ref, reactive } from "vue";
    import FeedbackRectsTable from '@/components/FeedbackRectsTable.vue';
    import FeedbackCanvas from '@/components/FeedbackCanvas.vue';
    import Tag from 'primevue/tag';

    const props = defineProps({
        query: {
            type: Object,
            required: true,
        },
        queryCount: {
            type: Number,
            required: true,
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

            emit('checked', props.query.HumanFeedbackQuery.ID, rectsSelected.value);
        }
    }

    const onRectChecked = (queryID, rect) => {
        const selectedIndex = rectsSelected.value.findIndex(selectedRect => selectedRect.ID == rect.ID);

        if (selectedIndex > -1) {
            rectsSelected.value.splice(selectedIndex, 1);
        } else {
            rectsSelected.value.push(rect);
        }

        emit('checked', props.query.HumanFeedbackQuery.ID, rectsSelected.value);
    }

    const getTagSeverity = (status) => {
        switch (status.ID) {
            case 1:
                return "danger";
            case 2:
                return "success";
            case 3:
                return "info";
        }
    }

</script>

<template>
    <div class="border-1 surface-border border-round m-2 text-center py-5 px-3">
        <FeedbackCanvas :id="query.HumanFeedbackQuery.ID" :imageURL="query.ImageURL" :drawnRectangles="rectsSelected" @mouseClick="onImageClicked" />
        <div>
            <h4 class="mb-1">Query {{ query.HumanFeedbackQuery.QueryID + 1 }}/{{ queryCount }}</h4>
            <h4 class="mb-1">Epoch {{ query.HumanFeedbackQuery.Epoch + 1 }}</h4>
            <Tag :severity="getTagSeverity(query.HumanFeedbackQuery.QueryStatus)" :value="query.HumanFeedbackQuery.QueryStatus.Name" />
        </div>
        <FeedbackRectsTable :items="query.HumanFeedbackRects" :queryID="query.HumanFeedbackQuery.QueryID" :checkedRows="rectsSelected" @checked="onRectChecked" checkable />
    </div>
</template>