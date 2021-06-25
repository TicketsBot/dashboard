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
          <DiscordMessages {ticketId} {messages} on:send={(msg) => console.log(msg)} />
        </div>
      </div>
    </Card>
  </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, notifySuccess, notifyRatelimit, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Input from "../components/form/Input.svelte";
    import {navigateTo} from "svelte-router-spa";
    import DiscordMessages from "../components/DiscordMessages.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let ticketId = currentRoute.namedParams.ticketid;

    let closeReason = '';
    let messages = [];

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

    async function loadMessages() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tickets/${ticketId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        messages = res.data.messages;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadMessages();
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

    .manage {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-top: 12px;
    }

    .col {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 49%;
        height: 100%;
    }

    table.nice > tbody > tr:first-child {
        border-top: 1px solid #dee2e6;
    }
</style>
