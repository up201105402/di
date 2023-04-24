import { defineStore } from "pinia";
import axios from "axios";

export const usePipelinesStore = defineStore("pipelines", {
  state: () => ({
    pipelines: [],
  }),
  actions: {
    fetch(sampleDataKey) {
      axios
        .get(`data-sources/${sampleDataKey}.json`)
        .then((r) => {
          if (r.data && r.data.data) {
            this[sampleDataKey] = r.data.data;
          }
        })
        .catch((error) => {
          alert(error.message);
        });
    },
  },
});
