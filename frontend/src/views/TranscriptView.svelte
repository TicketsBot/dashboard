<iframe srcdoc={html} style="border: none; width: 100%; height: 100%">
</iframe>

<script>
    import axios from "axios";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {errorPage, withLoadingScreen} from '../js/util'
    import {API_URL} from "../js/constants";

    export let currentRoute;
    export let params = {};

    let guildId = currentRoute.namedParams.id;
    let ticketId = currentRoute.namedParams.ticketid;

    setDefaultHeaders();

    let html = '';

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/${guildId}/transcripts/${ticketId}/render`);
        if (res.status !== 200) {
            errorPage(res.data.error);
            return;
        }

        html = res.data;
    }

    withLoadingScreen(loadData);
</script>

<style>
    body {
        margin: 0;
        display: flex;
        width: 100%;
        height: 100%;
        background-color: #2e3136;
    }

    .discord-container {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        overflow-y: scroll;

        background-color: #2e3136;
        font-family: 'Poppins', sans-serif !important;
        font-weight: 400 !important;
        font-size: 16px;
    }

    .channel-header {
        background-color: #1e2124;
        height: 50px;
        width: 100%;

        text-align: center;
        display: flex;
        align-items: center;
    }

    .channel-name {
        color: white;
        padding-left: 20px;
    }

    #message-container {
        display: flex;
        flex-direction: column;
        flex: 1;
        height: 100%;
    }

    .message {
        display: flex;
        align-items: center;
        color: white !important;
        word-wrap: break-word;
        margin: 4px 0 4px 12px;
    }

    .username, .content, .avatar {
        margin-left: 4px;
    }

    .avatar {
        height: 32px;
        width: 32px;
        border-radius: 50%;
    }

    .attachment {
        margin-left: 12px;
        color: white;
    }

    @media only screen and (max-width: 576px) {
        .timestamp {
            display: none;
        }
    }
</style>