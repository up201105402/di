<script setup>
    import { computed, ref, onMounted } from "vue";
    import { useMainStore } from "@/stores/main";
    import {
      mdiAccountMultiple,
      mdiCartOutline,
      mdiChartTimelineVariant,
      mdiMonitorCellphone,
      mdiReload,
      mdiGithub,
      mdiChartPie,
    } from "@mdi/js";
    import * as chartConfig from "@/components/Charts/chart.config.js";
    import LineChart from "@/components/Charts/LineChart.vue";
    import SectionMain from "@/components/SectionMain.vue";
    import CardBoxWidget from "@/components/CardBoxWidget.vue";
    import CardBox from "@/components/CardBox.vue";
    import TableSampleClients from "@/components/TableSampleClients.vue";
    import NotificationBar from "@/components/NotificationBar.vue";
    import BaseButton from "@/components/BaseButton.vue";
    import CardBoxTransaction from "@/components/CardBoxTransaction.vue";
    import CardBoxClient from "@/components/CardBoxClient.vue";
    import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
    import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
    import PipelinesTable from "@/components/PipelinesTable.vue";
  
    const mainStore = useMainStore();
  
    const clientBarItems = computed(() => mainStore.clients.slice(0, 4));

  </script>
  
  <template>
    <LayoutAuthenticated>
      <SectionMain>
        <SectionTitleLineWithButton :icon="mdiChartTimelineVariant" :title="$t('pages.pipelines.name')" main>
          <span />
        </SectionTitleLineWithButton>

        <PipelinesTable checkable />
        
      </SectionMain>
    </LayoutAuthenticated>
  </template>