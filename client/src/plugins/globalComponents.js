import BaseButton from "@/components/BaseButton.vue";

const GlobalComponents = {
  install(Vue) {
    Vue.component(BaseButton.name, BaseButton);
  }
};

export default GlobalComponents;
