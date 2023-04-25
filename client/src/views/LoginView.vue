<script setup>
  import { reactive, computed } from "vue";
  import { useRouter } from "vue-router";
  import { mdiAccount, mdiAsterisk, mdiClose } from "@mdi/js";
  import SectionFullScreen from "@/components/SectionFullScreen.vue";
  import CardBox from "@/components/CardBox.vue";
  import OverlayLayer from "@/components/OverlayLayer.vue";
  import FormCheckRadio from "@/components/FormCheckRadio.vue";
  import FormField from "@/components/FormField.vue";
  import FormControl from "@/components/FormControl.vue";
  import BaseButton from "@/components/BaseButton.vue";
  import BaseButtons from "@/components/BaseButtons.vue";
  import LayoutGuest from "@/layouts/LayoutGuest.vue";
  import SectionTitle from "@/components/SectionTitle.vue";
  import CardBoxComponentTitle from "@/components/CardBoxComponentTitle.vue";
  import ErrorModal from '@/components/ErrorModal.vue';
  import { storeToRefs } from "pinia";

  import { useAuthStore } from '@/stores/auth.js';

  const { error } = storeToRefs(useAuthStore());

  const form = reactive({
    username: "",
    password: "",
    remember: true,
  });

  let displayError = computed({
    get: () => (error.value != null && error.value !== ''),
    set: (value) => !value ? error.value = null : error.value = error.value,
  })

  const router = useRouter();

  const signInAndRedirect = () => {
    useAuthStore().signIn(form.username, form.password, router, '/pipelines');
  };

  const acknowledge = () => {
    displayError = false;
  };

  window.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && displayError) {
      acknowledge();
    }
  });

</script>

<template>
  <LayoutGuest>
    <SectionFullScreen v-slot="{ cardClass }" bg="purplePink">
      <CardBox :class="cardClass" is-form @submit.prevent="signInAndRedirect">
        <SectionTitle>{{ $t('pages.login.name') }}</SectionTitle>

        <FormField label="Username" help="Please enter your username">
          <FormControl v-model="form.username" :icon="mdiAccount" name="login" autocomplete="username"
            placeholder="Username" />
        </FormField>

        <FormField label="Password" help="Please enter your password">
          <FormControl v-model="form.password" :icon="mdiAsterisk" type="password" name="password"
            autocomplete="current-password" placeholder="Password" />
        </FormField>

        <!-- <FormCheckRadio v-model="form.remember" name="remember" label="Remember" :input-value="true" /> -->

        <template #footer>
          <BaseButtons>
            <BaseButton type="submit" color="info" :label="$t('pages.login.submit')" />
            <BaseButton to="/signup" color="info" outline label="Sign Up" />
          </BaseButtons>
        </template>
      </CardBox>
      <ErrorModal :title="'Error loggin in'" v-model="displayError" :errorMessage="error || ''"></ErrorModal>
    </SectionFullScreen>
  </LayoutGuest>
</template>