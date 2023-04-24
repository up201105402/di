<template>
  <div class="box">
    <div class="banner_high">
      <img src="/assets/logo.svg" alt="">
    </div>

    <h1>{{ $t('pages.register.name') }}</h1>

    <div v-if="error" class="text-center my-2">
      <p class="text-red">{{ error }}</p>
    </div>

    <div class="in">
      <label for="name">{{ $t('pages.register.username.name') }}</label>

      <div>
        <input type="text" v-model="username" :placeholder="$t('pages.register.username.name')" />
        <img v-if="username.length >= 4" src="/assets/ok.svg" alt="">
      </div>
    </div>

    <div class="in">
      <label for="name">{{ $t('pages.register.password.name') }}</label>
      <div>
        <input v-model="password" type="password" :name="string"
          :placeholder="$t('pages.register.password.placeholder')" required />
        <img v-if="password.length >= 8" src="/assets/ok.svg" alt="">
      </div>
    </div>

    <div class="password_bar">
      <div :class="{'bar':true, 'green':(password.length > 1)}"></div>
      <div :class="{'bar':true, 'green':(password.length > 3)}"></div>
      <div :class="{'bar':true, 'green':(password.length > 5)}"></div>
      <div :class="{'bar':true, 'green':(password.length > 7)}"></div>
    </div>

    <button class="log" @click="signUpAndRedirect()">
      {{ $t('pages.register.submit') }}
    </button>
  </div>
</template>

<script>
  import { useAuth } from '../stores/auth.js';
  import { ref } from 'vue';
  import { useRouter } from 'vue-router';

  export default {
    data() {
      return {
        username: '',
        password: '',
      };
    },
    setup() {
      const { signUp, error } = useAuth();
      const router = useRouter();
      const username = ref('');
      const password = ref('');
      const err = ref('');

      const signUpAndRedirect = function () {
        signUp(username.value, password.value);
        if (!error.value) {
          router.push('/')
        }
      }

      return {
        signUpAndRedirect,
        username,
        password,
        error
      };
    }
  };
</script>