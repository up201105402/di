<script setup>
import { computed, ref } from "vue";
import { useAuthStore } from "@/stores/auth.js";
import { mdiCheckDecagram } from "@mdi/js";
import BaseLevel from "@/components/BaseLevel.vue";
import UserAvatarCurrentUser from "@/components/UserAvatarCurrentUser.vue";
import CardBox from "@/components/CardBox.vue";
import FormCheckRadio from "@/components/FormCheckRadio.vue";
import PillTag from "@/components/PillTag.vue";

const authStore = useAuthStore();

const userName = computed(() => authStore.userName);

const userSwitchVal = ref(false);
</script>

<template>
  <CardBox>
    <BaseLevel type="justify-around lg:justify-center">
      <UserAvatarCurrentUser class="lg:mx-12" />
      <div class="space-y-3 text-center md:text-left lg:mx-12">
        <h1 class="text-2xl" v-html="$t('pages.profile.user.greeting', { username: userName })" />
      </div>
    </BaseLevel>
  </CardBox>
</template>
