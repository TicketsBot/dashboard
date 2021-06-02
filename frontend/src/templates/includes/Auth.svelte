<script context="module">
    import axios from 'axios';
    import {API_URL, OAUTH} from "../js/constants";

    const _tokenKey = 'token';

    export function getToken() {
        let token = window.localStorage.getItem(_tokenKey);
        if (token == null) {
            redirectLogin();
            return;
        }

        return token;
    }

    export function setToken(token) {
        window.localStorage.setItem(_tokenKey, token);
    }

    export function redirectLogin() {
        // TODO: State
        window.location.href = `https://discordapp.com/oauth2/authorize?response_type=code&redirect_uri=${OAUTH.redirectUri}&scope=identify%20guilds&client_id=${OAUTH.clientId}&state=`
    }

    export function clearLocalStorage() {
        window.localStorage.clear();
    }

    export function setDefaultHeaders() {
        axios.defaults.headers.common['Authorization'] = getToken();
        axios.defaults.headers.common['x-tickets'] = 'true'; // arbitrary header name and value
        axios.defaults.validateStatus = false;
    }

    function addRefreshInterceptor() {
        axios.interceptors.response.use(async (res) => { // we set validateStatus to false
            if (res.status === 401) {
                await _refreshToken();
            }
            return res;
        }, async (err) => {
            if (err.response.status === 401) {
                await _refreshToken();
            }
            return err.response;
        });
    }
</script>