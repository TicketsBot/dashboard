<Head/>

<div class="wrapper">
    <Navbar {guildId} {permissionLevel} {dropdown} />
    <div class="super-container" class:dropdown={$dropdown}>
        <LoadingScreen/>
        <div class="content-container" class:hide={$loadingScreen}>
            <Route {currentRoute} {params}/>
        </div>
        <NotifyModal/>
    </div>
</div>

<style>
    body {
        padding: 0 !important;
    }

    .wrapper {
        margin: 0 !important;
        padding: 0 !important;
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
</style>

<script>
    import Head from '../includes/Head.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import Navbar from '../includes/Navbar.svelte'

    import {Route} from 'svelte-router-spa'
    import {dropdown, loadingScreen} from '../js/stores'
    import {notifyError, withLoadingScreen} from '../js/util';
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'

    export let currentRoute;
    export let params = {};

    let guildId = currentRoute.namedParams.id;
    let permissionLevel = 0;

    setDefaultHeaders();

    async function loadPermissionLevel() {
        const res = await axios.get(`${API_URL}/user/permissionlevel?guild=${guildId}`);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        permissionLevel = res.data.permission_level;
    }

    withLoadingScreen(async () => {
        await loadPermissionLevel();
    });
</script>