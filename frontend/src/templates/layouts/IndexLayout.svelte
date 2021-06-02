<Head/>

<div class="wrapper">
  <Sidebar referralShow=true referralLink="https://www.digitalocean.com?refcode=371f56712ea4" name="{name}" avatar="{avatar}" />
  <div class="super-container">
    <LoadingScreen/>
    <NotifyModal/>
    <div class="content-container">
      <slot/>
    </div>
  </div>
</div>

<style>
    body {
        padding: 0 !important;
    }

    .wrapper {
        margin: 0 !important;
        padding: 0 !important;
        display: flex;
        width: 100%;
        height: 100%;
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
</style>

<script>
    import Head from '../includes/Head.svelte'
    import Sidebar from '../includes/Sidebar.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {notifyError} from "../js/util";
    import {redirectLogin} from '../includes/Auth.svelte'

    let name;
    let avatar;

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
    }

    loadData();
</script>