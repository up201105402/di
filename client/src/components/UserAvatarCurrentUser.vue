<script setup>
import { useAuthStore } from "@/stores/auth.js";
import UserAvatar from "@/components/UserAvatar.vue";

const authStore = useAuthStore();
</script>

<template>
  <UserAvatar :username="authStore.userName" :avatar="authStore.userAvatar">
    <slot />
  </UserAvatar>
</template>
