<Head/>

<div class="wrapper">
  <AdminSidebar />
  <div class="super-container">
    <LoadingScreen/>
    <NotifyModal/>
    <div class="content-container" class:hide={$loadingScreen}>
      <Route {currentRoute} {params}/>
    </div>
  </div>
</div>

<script>
    import {navigateTo, Route} from 'svelte-router-spa'
    import Head from '../includes/Head.svelte'
    import Sidebar from '../includes/Sidebar.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import axios from "axios";
    import {API_URL} from '../js/constants'
    import {notifyError} from '../js/util'
    import {loadingScreen} from "../js/stores"
    import {redirectLogin, setDefaultHeaders} from '../includes/Auth.svelte'
    import AdminSidebar from "../includes/AdminSidebar.svelte";
    import {onMount} from "svelte";

    export let currentRoute;
    export let params = {};

    setDefaultHeaders();

    onMount(() => {
        let isAdmin = false;
        try {
            const userData = JSON.parse(window.localStorage.getItem('user_data'));
            isAdmin = userData.admin;
        } finally {
            if (!isAdmin) {
                navigateTo(`/`);
            }
        }
    });
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