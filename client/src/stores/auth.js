
import { watchEffect } from 'vue';
import { storeTokens, getTokens, doRequest, getTokenPayload, removeTokens } from '../util';
import { useRouter } from 'vue-router';

import { defineStore, storeToRefs } from "pinia";
import axios from "axios";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    /* User */
    userName: null,
    userEmail: null,
    userAvatar: null,
    idToken: null,

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

    fetch(sampleDataKey) {
      axios
        .get(`data-sources/${sampleDataKey}.json`)
        .then((r) => {
          if (r.data && r.data.data) {
            this[sampleDataKey] = r.data.data;
          }
        })
        .catch((error) => {
          alert(error.message);
        });
    },

    async logIn(username, password, router, redirectURL) {
        const { data, error } = await authenticate(username, password, '/api/user/login')

        if (error) {
            this.error = error;
            return;
        }

        this.userName = username;
        this.error = null;
        this.idToken = data.token.idToken;
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
        this.idToken = data.idToken;
        router.push(redirectURL);
    },

    async signOut() {
        const { error } = await doRequest({
            url: '/api/user/signout',
            method: 'POST',
            headers: {
                Authorization: `${this.idToken}`,
            },
        });
    
        if (error) {
            // state.error = error;
            // state.isLoading = false;
            this.error = error;
            return;
        }
    
        // state.currentUser = null;
        // state.idToken = null;

        this.userName = null;
        this.idToken = null;
    
        removeTokens();
    }

  },
});

const initializeUser = async () => {
    state.isLoading = true;
    state.error = null;

    const [idToken, refreshToken] = getTokens();

    const idTokenClaims = getTokenPayload(idToken);
    const refreshTokenClaims = getTokenPayload(refreshToken);

    if (idTokenClaims) {
        state.idToken = idToken;
        state.currentUser = idTokenClaims.user;
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
    storeTokens(tokens.idToken, tokens.refreshToken);

    const updatedIdTokenClaims = getTokenPayload(tokens.idToken);

    state.currentUser = updatedIdTokenClaims.user;
    state.idToken = tokens.idToken;
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

    storeTokens('tokens.idToken_'  + username, 'tokens.refreshToken');

    //const tokenClaims = getTokenPayload(tokens.idToken);

    // set tokens to local storage with expiry (separate function)
    // state.idToken = 'tokens.idToken_' + username;
    //state.currentUser = tokenClaims.user;
    // state.currentUser = username;
    // state.isLoading = false;

    return {
        data,
        error
    }
};
