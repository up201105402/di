<template>
  <router-view />
</template>

<script>
  import { defineComponent } from 'vue';
  import { useAuth } from "@/stores/auth.js";

  export default defineComponent({
    name: 'App',
    setup() {
      useAuth();
    }
  })
</script>