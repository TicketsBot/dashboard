<Head/>

<div class="wrapper">
  <Navbar guildId={guildId} bind:dropdown={$dropdown}/>
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
    export let currentRoute;
    export let params = {};

    let guildId = currentRoute.namedParams.id

    import Head from '../includes/Head.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import Navbar from '../includes/Navbar.svelte'

    import {Route} from 'svelte-router-spa'
    import {loadingScreen, dropdown} from '../js/stores'
</script>