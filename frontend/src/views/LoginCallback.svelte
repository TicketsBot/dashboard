<script>
    import axios from "axios";
    import {redirectLogin, setToken} from '../includes/Auth.svelte'
    import {API_URL} from "../js/constants";
    import {errorPage} from '../js/util'
    import {navigateTo} from "svelte-router-spa";

    export let currentRoute;
    let state = currentRoute.queryParams.state;

    async function process() {
        const code = new URLSearchParams(window.location.search).get('code')
        if (code === null) {
            redirectLogin()
        }

        axios.defaults.validateStatus = false
        axios.defaults.headers.common['x-tickets'] = 'true'
        const res = await axios.post(`${API_URL}/callback?code=${code}`)
        if (res.status !== 200) {
            errorPage(res.data.error)
            return
        }

        setToken(res.data.token);
        window.localStorage.setItem('user_data', JSON.stringify(res.data.user_data));
        if (res.data.guilds) {
            window.localStorage.setItem('guilds', JSON.stringify(res.data.guilds));
        }

        let path = '/';

        try {
            if (state !== undefined && state.length > 0) {
                path = atob(state);

                if (path === '/callback') {
                    path = '/';
                }

                try {
                    new URL(path);
                    path = '/';
                } catch (e) {}
            }
        } catch (e) {
            console.log(`Error parsing state: ${e}`)
        } finally {
            navigateTo(path);
        }
    }

    process()
</script>

<style>
    body {
        background-color: #121212;
    }
</style>