<Head/>

<div class="wrapper">
  <Sidebar name="{name}" avatar="{avatar}" {isWhitelabel} />
  <div class="super-container">
    <LoadingScreen/>
    <NotifyModal/>
    <div class="content-container" class:hide={$loadingScreen}>
      <Route {currentRoute} {params}/>
    </div>
  </div>
</div>

<script>
    import {Route} from 'svelte-router-spa'
    import Head from '../includes/Head.svelte'
    import Sidebar from '../includes/Sidebar.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import axios from "axios";
    import {API_URL} from '../js/constants'
    import {notifyError} from '../js/util'
    import {loadingScreen} from "../js/stores"
    import {redirectLogin, setDefaultHeaders} from '../includes/Auth.svelte'

    export let currentRoute;
    export let params = {};

    setDefaultHeaders()

    let name;
    let avatar;

    let isWhitelabel = false;

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/session`);
        if (res.status !== 200) {
            if (res.data.auth === true) {
                redirectLogin();
            }

            notifyError(res.data.error);
            return;
        }

        name = res.data.username;
        avatar = res.data.avatar;
        isWhitelabel = res.data.whitelabel;
    }

    loadData();
</script>

<style>
    body {
        padding: 0 !important;
    }

    .wrapper {
        display: flex;
        width: 100%;
        height: 100%;
        margin: 0 !important;
        padding: 0 !important;
    }

    .super-container {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .content-container {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .hide {
        visibility: hidden;
    }

    @media (max-width: 950px) {
        .wrapper {
            flex-direction: column;
        }
    }
</style>