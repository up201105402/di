import { createApp } from 'vue';
import { createPinia  } from "pinia";
import { i18n } from '@/i18n';
import App from '@/App.vue';

import router from "@/router";
import { useAuthStore } from "@/stores/auth";
import { useStyleStore } from "@/stores/style.js";
import { darkModeKey, styleKey } from "@/config.js";
import '@/css/main.css';

import { plugin as formKitPlugin } from '@formkit/vue';
import formkitConfig from './../formkit.config'
import '@formkit/themes/genesis'
import '@formkit/addons/css/multistep'

// PrimeVue
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
//theme
import "primevue/resources/themes/lara-light-indigo/theme.css";     
//core
import "primevue/resources/primevue.min.css";

/* Init Pinia */
const pinia = createPinia();
const authStore = useAuthStore(pinia);
const styleStore = useStyleStore(pinia);

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

const spa = createApp(App);
spa.use(authStore);
spa.use(router);
spa.use(i18n);
spa.use(PrimeVue);
spa.use(ToastService);
spa.use(formKitPlugin, formkitConfig);
spa.mount('#app')
