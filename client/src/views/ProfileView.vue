<script setup>
  import { reactive, computed } from "vue";
  import { useAuthStore } from "@/stores/auth.js";
  import {
    mdiAccount,
    mdiMail,
    mdiAsterisk,
    mdiFormTextboxPassword,
    mdiGithub,
  } from "@mdi/js";
  import SectionMain from "@/components/SectionMain.vue";
  import CardBox from "@/components/CardBox.vue";
  import BaseDivider from "@/components/BaseDivider.vue";
  import FormField from "@/components/FormField.vue";
  import FormControl from "@/components/FormControl.vue";
  import FormFilePicker from "@/components/FormFilePicker.vue";
  import BaseButton from "@/components/BaseButton.vue";
  import BaseButtons from "@/components/BaseButtons.vue";
  import UserCard from "@/components/UserCard.vue";
  import LayoutAuthenticated from "@/layouts/LayoutAuthenticated.vue";
  import SectionTitleLineWithButton from "@/components/SectionTitleLineWithButton.vue";
  import Toast from 'primevue/toast';
  import { useToast } from 'primevue/usetoast';
  import { useRouter } from 'vue-router';
  import { i18n } from '@/i18n';
  
  const { t } = i18n.global;

  const authStore = useAuthStore();
  const router = useRouter();
  const toast = useToast();

  const profileForm = reactive({
    name: authStore.userName,
    email: authStore.userEmail,
  });

  const passwordForm = reactive({
    password_current: "",
    password: "",
    password_confirmation: "",
  });

  const submitProfile = () => {
    useAuthStore().editUsername(profileForm.name, router, '/pipelines');
  };

  const isNewPasswordValid = computed(() => passwordForm.password_confirmation != "" && passwordForm.password == passwordForm.password_confirmation)

  const submitPass = () => {
    if (isNewPasswordValid.value) {
      toast.add({ severity: 'error', summary: 'Error', detail: t('pages.profile.form.errors.newPassword.notEqual'), life: 3000 });
    } else {
      useAuthStore().editPassword(passwordForm.password, router, '/pipelines');
    }
  };
  
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <SectionTitleLineWithButton :hasButton="false" :icon="mdiAccount" :title="$t('pages.profile.header')" main />

      <UserCard class="mb-6" />

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <CardBox is-form @submit.prevent="submitProfile">

          <FormField :label="$t('pages.profile.form.name.label')" :help="$t('pages.profile.form.name.help')">
            <FormControl v-model="profileForm.name" :icon="mdiAccount" name="username" required autocomplete="username" />
          </FormField>

          <template #footer>
            <BaseButtons>
              <BaseButton color="success" type="submit" :label="$t('buttons.submit')" />
            </BaseButtons>
          </template>
        </CardBox>

        <CardBox is-form @submit.prevent="submitPass">
          <FormField :label="$t('pages.profile.form.currentPassword.label')" :help="$t('pages.profile.form.currentPassword.help')">
            <FormControl v-model="passwordForm.password_current" :icon="mdiAsterisk" name="password_current" type="password" required autocomplete="current-password" />
          </FormField>

          <BaseDivider />

          <FormField :label="$t('pages.profile.form.newPassword.label')" :help="$t('pages.profile.form.newPassword.help')">
            <FormControl v-model="passwordForm.password" :icon="mdiFormTextboxPassword" name="password" type="password" required autocomplete="new-password" />
          </FormField>

          <FormField :label="$t('pages.profile.form.confirmNewPassword.label')" :help="$t('pages.profile.form.confirmNewPassword.help')">
            <FormControl id="password-confirmation" v-model="passwordForm.password_confirmation" :icon="mdiFormTextboxPassword" name="password_confirmation" type="password" required autocomplete="new-password" />
          </FormField>

          <template #footer>
            <BaseButtons>
              <BaseButton color="success" type="submit" :label="$t('buttons.submit')" />
            </BaseButtons>
          </template>
        </CardBox>
      </div>
    </SectionMain>

    <Toast />
  </LayoutAuthenticated>
</template>