import { watchEffect } from 'vue';
import jwt_decode from 'jwt-decode';
import { doRequest } from '@/util';
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
            this.accessToken = accessToken.signedString;
            this.refreshToken = refreshToken.signedString;
            router.push(redirectURL);
        },
        async signUp(username, password, router, redirectURL) {
            const { data, error } = await authenticate(username, password, '/api/user/signup');

            if (error) {
                store.accessToken = store.refreshToken = store.userName = null;
                removeTokens();
                return;
            }

            const { accessToken, refreshToken } = data.tokens;
            store.accessToken = accessToken.signedString;
            store.refreshToken = refreshToken.signedString;
            this.userName = username;
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

export const useAuth = async () => {

    const store = storeToRefs(useAuthStore());
    const router = useRouter();

    const isAccessTokenValid = isTokenValid(store.accessToken.value);
    const isRefreshTokenValid = isTokenValid(store.refreshToken.value);

    if (isRefreshTokenValid) {
        if (!isAccessTokenValid) {
            const result = await getNewAccessToken();

            if (!result) {
                redirectToAuthPage(router);
            }
        }
    } else {
        redirectToAuthPage(router);
    }

    return store;
}

const redirectToAuthPage = async (router) => {

    const store = storeToRefs(useAuthStore());

    store.accessToken.value = store.refreshToken.value = store.userName.value = null;
    removeTokens();

    const currentRoutePath = router.currentRoute.value.path;

    if (!store.publicRoutePaths.value.find(elem => elem === currentRoutePath)) {
        router.push(store.requireAuthRoute.value);
    }
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

export const getNewAccessToken = async () => {

    const store = storeToRefs(useAuthStore());
    const router = useRouter();

    const { data, status, error } = await doRequest({
        url: '/api/user/tokens',
        method: 'POST',
        data: {
            refreshToken: store.refreshToken.value,
        },
    }, false);

    if (error) {
        store.accessToken.value = store.refreshToken = store.userName = null;
        removeTokens();
        router.push(store.requireAuthRoute.value);
        return false;
    } else {
        const { accessToken, refreshToken } = data.tokens;
        store.accessToken.value = accessToken.signedString;
        store.refreshToken = refreshToken.signedString;
    }

    return true;
}