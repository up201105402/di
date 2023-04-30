
import { watchEffect } from 'vue';
import { storeTokens, getTokens, doRequest, getTokenPayload, removeTokens } from '../util';
import { useRouter } from 'vue-router';
import { useStorage } from '@vueuse/core'

import { defineStore, storeToRefs } from "pinia";
import axios from "axios";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    /* User */
    userName: null,
    userEmail: null,
    userAvatar: null,
    accessToken: useStorage('accessToken', null),
    refreshToken: useStorage('refreshToken', null),

    error: null,

    /* Field focus with ctrl+k (to register only once) */
    isFieldFocusRegistered: false,

    /* Sample data (commonly used) */
    clients: [],
    history: [],

    onAuthRoute: '/',
    requireAuthRoute: '/login',
    publicRoutePaths: ['/signup', '/login']
  }),
  actions: {
    setUser(payload) {
      if (payload.name) {
        this.userName = payload.name;
      }
      if (payload.email) {
        this.userEmail = payload.email;
      }
      if (payload.avatar) {
        this.userAvatar = payload.avatar;
      }
    },
    async logIn(username, password, router, redirectURL) {
        const { data, error } = await authenticate(username, password, '/api/user/login')

        if (error) {
            this.error = error;
            return;
        }

        this.userName = username;
        this.error = null;
        this.accessToken = data.token.idToken;
        router.push(redirectURL);
    },
    
    async signUp(username, password, router, redirectURL) {
        const { data, error } = await authenticate(username, password, '/api/user/signup');

        if (error) {
            this.error = error;
            return;
        }

        this.userName = username;
        this.error = null;
        this.accessToken = data.idToken;
        router.push(redirectURL);
    },

    async signOut() {
        const { error } = await doRequest({
            url: '/api/user/signout',
            method: 'POST',
            headers: {
                Authorization: `${this.accessToken}`,
            },
        });
    
        if (error) {
            this.error = error;
            return;
        }
        
        this.userName = null;
        this.accessToken = null;
    
        removeTokens();
    }

  },
});

const initializeUser = async () => {
    state.isLoading = true;
    state.error = null;

    const [accessToken, refreshToken] = getTokens();

    const accessTokenClaims = getTokenPayload(idToken);
    const refreshTokenClaims = getTokenPayload(refreshToken);

    if (accessTokenClaims) {
        state.accessToken = accessToken;
        state.currentUser = accessTokenClaims.user;
    }

    state.isLoading = false;

    // silently refresh tokens in local storage
    // if we for some reason don't have refresh token (e.g., if the user deleted it manually)
    // then we don't proceed
    if (!refreshTokenClaims) {
        return;
    }

    const { data, error } = await doRequest({
        url: '/api/user/tokens',
        method: 'POST',
        data: {
            refreshToken,
        },
    });

    if (error) {
        console.error('Error refreshing tokens\n', error);
        return;
    }

    const { tokens } = data;
    storeTokens(tokens.accessToken, tokens.refreshToken);

    const updatedaccessTokenClaims = getTokenPayload(tokens.accessToken);

    state.currentUser = updatedAccessTokenClaims.user;
    state.accessToken = tokens.accessToken;
};

export const useAuth = () => {

    const store = useAuthStore();

    if (!store) {
        throw new Error('Main store has not been initialized!');
    }

    const router = useRouter();

    watchEffect(() => {
        const currentRoutePath = router.currentRoute.value.path;

        const { userName, requireAuthRoute } = storeToRefs(store)

        if (!store.publicRoutePaths.find(elem => elem === currentRoutePath)) {
            if (!userName.value && requireAuthRoute.value) {
                router.push(requireAuthRoute.value);
            }
        }
    });

    return store;
}

const authenticate = async (username, password, url) => {

    const { data, error } = await doRequest({
        url,
        method: 'POST',
        data: {
            username,
            password,
        },
    });

    // if (error) {
    //     state.error = error;
    //     state.isLoading = false;
    //     return;
    // }

    // const { tokens } = data;

    storeTokens('tokens.accessToken_'  + username, 'tokens.refreshToken');

    //const tokenClaims = getTokenPayload(tokens.accessToken);

    // set tokens to local storage with expiry (separate function)
    // state.accessToken = 'tokens.accessToken_' + username;
    //state.currentUser = tokenClaims.user;
    // state.currentUser = username;
    // state.isLoading = false;

    return {
        data,
        error
    }
};
