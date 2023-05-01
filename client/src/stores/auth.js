
import { watchEffect } from 'vue';
import jwt_decode from 'jwt-decode';
import { doRequest } from '../util';
import { useRouter } from 'vue-router';
import { useStorage } from '@vueuse/core'

import { defineStore, storeToRefs } from "pinia";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    userName: useStorage('userName', null),
    userEmail: null,
    userAvatar: null,
    accessToken: useStorage('accessToken', null),
    refreshToken: useStorage('refreshToken', null),

    error: null,

    onAuthRoute: '/',
    requireAuthRoute: '/login',
    publicRoutePaths: ['/signup', '/login']
  }),
  actions: {
    async logIn(username, password, router, redirectURL) {
        const { data, error } = await authenticate(username, password, '/api/user/login')

        if (error) {
            this.error = error;
            return;
        }

        this.userName = username;
        this.error = null;
        const { accessToken, refreshToken } = data.tokens;
        this.accessToken = accessToken.idToken;
        this.refreshToken = refreshToken.idToken;
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
        
        this.userName = this.accessToken = this.refreshToken = null;
        removeTokens();
    }
  },
});

const initializeUser = async () => {
    state.isLoading = true;
    state.error = null;

    const [accessToken, refreshToken] = getTokens();

    const accessTokenClaims = isTokenValid(idToken);
    const refreshTokenClaims = isTokenValid(refreshToken);

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

    const updatedaccessTokenClaims = isTokenValid(tokens.accessToken);

    state.currentUser = updatedAccessTokenClaims.user;
    state.accessToken = tokens.accessToken;
};

export const useAuth = () => {

    const store = useAuthStore();
    const router = useRouter();

    const isAccessTokenValid = isTokenValid(store.accessToken);
    const isRefreshTokenValid = isTokenValid(store.refreshToken);

    if (!isAccessTokenValid && isRefreshTokenValid) {
        const { data, error } = getNewAccessToken(this.refreshToken);

        if (error) {
            store.accessToken = store.refreshToken = store.userName = null;
            removeTokens();
        }

        const { accessToken, refreshToken } = data.tokens;
        store.accessToken = accessToken.idToken;
        store.refreshToken = refreshToken.idToken;
    } else {
        store.accessToken = store.refreshToken = store.userName = null;
        removeTokens();
    }

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

    return {
        data,
        error
    }
};

const removeTokens = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("userName");
};

const isTokenValid = (token) => {
    
    if (!token) {
        return false;
    }

    const payload = jwt_decode(token);

    if (Date.now() / 1000 >= payload.exp) {
        return false;
    }

    return true;
};

const getNewAccessToken = async (refreshToken) => {
    return await doRequest({
        url: '/api/user/tokens',
        method: 'POST',
        data: {
            refreshToken
        },
    });
}