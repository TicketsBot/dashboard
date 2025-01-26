<Head/>
<SunsetBanner />

<div class="wrapper">
    <Navbar {guildId} {permissionLevel} {dropdown} />
    <div class="super-container" class:dropdown={$dropdown}>
        <LoadingScreen/>
        <div class="content-container" class:hide={$loadingScreen}>
            <ManageSidebar {currentRoute} {permissionLevel} />
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

    .super-container {
        padding: 30px;
    }

    .content-container {
        display: flex;
        gap: 30px;
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
    import {permissionLevelCache} from '../js/stores';
    import {get} from 'svelte/store';
    import ManageSidebar from "../includes/ManageSidebar.svelte";
    import SunsetBanner from "../includes/SunsetBanner.svelte";

    export let currentRoute;
    export let params = {};

    let guildId = currentRoute.namedParams.id;
    let permissionLevel = 0;

    setDefaultHeaders();

    async function loadPermissionLevel() {
        const cache = get(permissionLevelCache);
        if (cache && cache[guildId]) {
            const data = cache[guildId];
            if (data.last_updated) {
                const date = new Date(data.last_updated.getTime() + 60000);
                if (date > new Date()) {
                    permissionLevel = data.permission_level;
                    return;
                }
            }
        }

        const res = await axios.get(`${API_URL}/user/permissionlevel?guild=${guildId}`);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        permissionLevel = res.data.permission_level;

        permissionLevelCache.update(cache => {
            cache[guildId] = {
                permission_level: permissionLevel,
                last_updated: new Date(),
            };

            return cache;
        });
    }

    withLoadingScreen(async () => {
        await loadPermissionLevel();
    });
</script>