<script setup>
    import { computed } from "vue";
    import { mdiClose } from "@mdi/js";
    import BaseButton from "@/components/BaseButton.vue";
    import BaseButtons from "@/components/BaseButtons.vue";
    import CardBox from "@/components/CardBox.vue";
    import OverlayLayer from "@/components/OverlayLayer.vue";
    import CardBoxComponentTitle from "@/components/CardBoxComponentTitle.vue";
    import { useAuthStore } from '@/stores/auth.js';

    const props = defineProps({
        title: {
            type: String,
            required: true,
        },
        errorMessage: {
            type: String,
            required: true,
        },
        submitLabel: {
            type: String,
            default: "Acknowledge",
        },
        modelValue: {
            type: [String, Number, Boolean],
            default: null,
        },
    });

    const emit = defineEmits(["update:modelValue", "acknowledge"]);

    const value = computed({
        get: () => props.modelValue,
        set: (value) => emit("update:modelValue", value),
    });

    const acknowledge = () => {
        value.value = false;
        emit("acknowledge");
    };

    window.addEventListener("keydown", (e) => {
        if (e.key === "Escape" && value.value) {
            acknowledge();
        }
    });

</script>

<template>
    <OverlayLayer v-show="value" z-index="z-30">
        <CardBox v-show="value" class="shadow-lg max-h-modal w-11/12 md:w-3/5 lg:w-2/5 xl:w-4/12 z-50" is-modal>
            <CardBoxComponentTitle :title="title">
                <BaseButton :icon="mdiClose" color="whiteDark" small rounded-full @click.prevent="acknowledge" />
            </CardBoxComponentTitle>

            <div class="space-y-3">
                <p>{{ errorMessage }}</p>
            </div>

            <template #footer>
                <BaseButtons>
                    <BaseButton :label="submitLabel" :color="'danger'" @click="acknowledge" />
                </BaseButtons>
            </template>
        </CardBox>
    </OverlayLayer>
</template>