<script setup>
    import { computed } from "vue";
    import { useMainStore } from "@/stores/main";
    import {
      mdiChartTimelineVariant,
    } from "@mdi/js";
    import SectionMain from "@/components/SectionMain.vue";
    import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
    import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
    import FlowChart from "@/components/FlowChart.vue";
  
    const mainStore = useMainStore();
  
    const clientBarItems = computed(() => mainStore.clients.slice(0, 4));

  </script>
  
  <template>
    <LayoutAuthenticated>
      <SectionMain>
        <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
          <span />
        </SectionTitleLineWithButton>

        <FlowChart v-if="$route.params.id" />
        
      </SectionMain>
    </LayoutAuthenticated>
  </template>