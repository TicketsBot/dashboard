<div class="parent">
  <div class="content">
    <Card footer={false}>
      <span slot="title">Support Teams</span>
      <div slot="body" class="body-wrapper">
        <div class="section">
          <h2 class="section-title">Close Ticket</h2>

          <form on:submit|preventDefault={closeTicket}>
            <div class="row" style="max-height: 63px; align-items: flex-end"> <!-- hacky -->
              <div class="col-3" style="margin-bottom: 0 !important;">
                <Input label="Close Reason" placeholder="No reason specified" col1={true} bind:value={closeReason}/>
              </div>
              <div class="col-1">
                <div style="margin-left: 30px">
                  <Button danger={true} icon="fas fa-lock">Close Ticket</Button>
                </div>
              </div>
            </div>
          </form>
        </div>
        <div class="section">
          <h2 class="section-title">View Ticket</h2>
          <DiscordMessages {ticketId} {isPremium} {messages} bind:container on:send={sendMessage} />
        </div>
      </div>
    </Card>
  </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, notifyRatelimit, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders, getToken} from '../includes/Auth.svelte'
    import Input from "../components/form/Input.svelte";
    import {navigateTo} from "svelte-router-spa";
    import DiscordMessages from "../components/DiscordMessages.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let ticketId = parseInt(currentRoute.namedParams.ticketid);

    let closeReason = '';
    let messages = [];
    let isPremium = false;
    let container;

    let WS_URL = env.WS_URL || 'ws://172.26.50.75:3000';

    function scrollContainer() {
        container.scrollTop = container.scrollHeight;
    }

    async function closeTicket() {
        let data = {
            reason: closeReason,
        };

        const res = await axios.delete(`${API_URL}/api/${guildId}/tickets/${ticketId}`, data);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/tickets`);
    }

    async function sendMessage(e) {
        let data = {
            message: e.detail,
        };

        const res = await axios.post(`${API_URL}/api/${guildId}/tickets/${ticketId}`, data);
        if (res.status !== 200) {
            if (res.status === 429) {
                notifyRatelimit();
            } else {
                notifyError(res.data.error);
            }
        }
    }

    function connectWebsocket() {
        const ws = new WebSocket(`${WS_URL}/webchat`);

        ws.onopen = () => {
            ws.send(JSON.stringify({
                "type": "auth",
                "data": {
                    "guild_id": guildId,
                    "ticket_id": ticketId,
                    "token": getToken(),
                }
            }));
        };

        ws.onmessage = (evt) => {
            const data = JSON.parse(evt.data);
            messages = [...messages, data];
            scrollContainer();
        };
    }

    async function loadMessages() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tickets/${ticketId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        messages = res.data.messages;
    }

    async function loadPremium() {
        const res = await axios.get(`${API_URL}/api/${guildId}/premium`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        isPremium = res.data.premium;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadPremium();
        await loadMessages();

        scrollContainer();

        if (isPremium) {
            connectWebsocket();
        }
    });
</script>

<style>
    .parent {
        display: flex;
        justify-content: center;
        width: 100%;
        height: 100%;
    }

    .content {
        display: flex;
        justify-content: space-between;
        width: 96%;
        height: 100%;
        margin-top: 30px;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        padding: 1%;
    }

    .section {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .section:not(:first-child) {
        margin-top: 2%;
    }

    .section-title {
        font-size: 36px;
        font-weight: bolder !important;
    }

    h3 {
        font-size: 28px;
        margin-bottom: 4px;
    }

    .row {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
    }

    @media only screen and (max-width: 576px) {
        .row {
            flex-direction: column;
        }
    }
</style>
