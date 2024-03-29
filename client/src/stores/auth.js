import jwt_decode from 'jwt-decode';
import { doRequest } from '@/util';
import { useRouter } from 'vue-router';
import { useStorage } from '@vueuse/core'
import { defineStore, storeToRefs } from "pinia";
import { i18n } from '@/i18n';

const { t } = i18n.global;

export const useAuthStore = defineStore("auth", {
    state: () => ({
        userName: useStorage('userName', null),
        userEmail: null,
        userAvatar: null,
        accessToken: useStorage('accessToken', null),
        refreshToken: useStorage('refreshToken', null),

        error: null,

        isLoading: false,

        onAuthRoute: '/',
        requireAuthRoute: '/login',
        publicRoutePaths: ['/signup', '/login']
    }),
    actions: {
        async logIn(username, password, router, redirectURL) {
            this.isLoading = true;

            const { data, error } = await authenticate(username, password, '/api/user/login')

            if (error) {
                this.error = error;
                this.isLoading = false;
                return;
            }

            this.userName = username;
            this.error = null;
            const { accessToken, refreshToken } = data.tokens;
            this.accessToken = accessToken.signedString;
            this.refreshToken = refreshToken.signedString;
            router.push(redirectURL);
            this.isLoading = false;
        },
        async editUsername(username, router, redirectURL) {
            this.isLoading = true;
            const { data, error } = await doRequest({
                url: '/api/user',
                method: 'POST',
                headers: {
                  Authorization: `${accessToken.value}`,
                },
                data: {
                    username: username,
                },
            });
              

            if (error) {
                toast.add({ severity: 'error', summary: 'Error', detail: error, life: 3000 });
                return;
            }

            this.error = null;
            const { accessToken, refreshToken } = data.tokens;
            this.accessToken = accessToken.signedString;
            this.refreshToken = refreshToken.signedString;
            toast.add({ severity: 'error', summary: t('messages.types.success'), detail: t('pages.profile.form.success.usernameChanged'), life: 3000 });
            this.isLoading = false;
        },
        async editPassword(password, toast) {
            this.isLoading = true;
            const { data, error } = await doRequest({
                url: '/api/user',
                method: 'POST',
                headers: {
                  Authorization: `${accessToken.value}`,
                },
                data: {
                    oldPassword: oldPassword,
                    password: password
                },
            });
              

            if (error) {
                toast.add({ severity: 'error', summary: t('messages.types.error'), detail: error, life: 3000 });
                return;
            }

            this.error = null;
            const { accessToken, refreshToken } = data.tokens;
            this.accessToken = accessToken.signedString;
            this.refreshToken = refreshToken.signedString;
            toast.add({ severity: 'error', summary: t('messages.types.success'), detail: t('pages.profile.form.success.passwordChanged'), life: 3000 });
            this.isLoading = false;
        },
        async signUp(username, password, router, redirectURL) {
            this.isLoading = true;
            const { data, error } = await authenticate(username, password, '/api/user/signup');

            if (error) {
                this.accessToken = this.refreshToken = this.userName = null;
                removeTokens();
                this.error = error;
                this.isLoading = false;
                return;
            }

            this.userName = username;
            this.error = null;
            const { accessToken, refreshToken } = data.tokens;
            this.accessToken = accessToken.signedString;
            this.refreshToken = refreshToken.signedString;
            router.push(redirectURL);
            this.isLoading = false;
        },
        async signOut(router, redirectURL) {
            this.isLoading = true;
            const { error } = await doRequest({
                url: '/api/user/signout',
                method: 'POST',
                headers: {
                    Authorization: `${this.accessToken}`,
                },
            });

            if (error) {
                this.error = error;
                this.isLoading = false;
                return;
            }

            this.userName = this.accessToken = this.refreshToken = null;
            removeTokens();
            router.push(redirectURL);
            this.isLoading = false;
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
            accessToken: store.accessToken.value,
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