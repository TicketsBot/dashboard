<Head/>

<NotifyModal/>
<Route {currentRoute} {params}/>

<style>
    body {
        padding: 0 !important;
        overflow-y: clip;
    }
</style>

<script>
    export let currentRoute;
    export let params = {};

    import Head from '../includes/Head.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import {Route} from 'svelte-router-spa'
</script>