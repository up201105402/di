import { createApp } from 'vue';
import { createPinia } from "pinia";
import { createI18n } from 'vue-i18n';
import { useMainStore } from './stores/main';
import { messages } from './i18n';
import App from './App.vue';

import router from "./router";
import { useAuthStore } from "@/stores/auth";
import { usePipelinesStore } from "@/stores/pipelines";
import { useStyleStore } from "@/stores/style.js";
import { darkModeKey, styleKey } from "@/config.js";
import './css/main.css';

const i18n = createI18n({
  legacy: false, // you must set `false`, to use Composition API
  locale: 'en', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages
})

/* Init Pinia */
const pinia = createPinia();

const mainStore = useMainStore(pinia);
const authStore = useAuthStore(pinia);
const pipelinesStore = usePipelinesStore(pinia);
const styleStore = useStyleStore(pinia);

/* Fetch sample data */
mainStore.fetch("clients");
mainStore.fetch("history");
pipelinesStore.fetch("pipelines");

/* App style */
styleStore.setStyle(localStorage[styleKey] ?? "basic");

/* Dark mode */
if (
  (!localStorage[darkModeKey] &&
    window.matchMedia("(prefers-color-scheme: dark)").matches) ||
  localStorage[darkModeKey] === "1"
) {
  styleStore.setDarkMode(true);
}

/* Default title tag */
const defaultDocumentTitle = "DI";

/* Set document title from route meta */
router.afterEach((to) => {
  document.title = to.meta?.title
    ? `${to.meta.title} — ${defaultDocumentTitle}`
    : defaultDocumentTitle;
});

const spa = createApp(App).use(authStore).use(router).use(i18n);
spa.mount('#app')
