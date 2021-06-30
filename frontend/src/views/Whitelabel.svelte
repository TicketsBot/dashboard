<div class="wrapper">
  <div class="content">
    <div class="content-col">
      <Card footer="{false}" fill="{false}">
        <h4 slot="title">Bot Token</h4>
        <div slot="body" class="full-width">
          <form class="full-width" onsubmit="return false;">
            <label class="form-label">Bot Token</label>
            <input name="token" type="text" bind:value={token} class="form-input full-width"
                   placeholder="xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxx.xxxxxxxxxxxxxxxxxxxxxxxxxxx">
            <p>Note: You will not be able to view the token after submitting it</p>

            <div class="buttons">
              <div class="col">
                <Button icon="fas fa-paper-plane" on:click={submitToken} fullWidth="{true}">Submit</Button>
              </div>
              <div class="col">
                <Button icon="fas fa-plus" on:click={invite} fullWidth="{true}" disabled="{bot.id === undefined}">
                  Generate Invite Link
                </Button>
              </div>
            </div>
          </form>
        </div>
      </Card>
    </div>
    <div class="content-col">
      <Card footer="{false}" fill="{false}">
        <h4 slot="title">Slash Commands</h4>
        <div slot="body" class="full-width">
          <form class="full-width" onsubmit="return false;">
            <label class="form-label">Interactions Endpoint URL</label>
            <input name="token" type="text" bind:value={interactionUrl} class="form-input full-width" readonly>

            <label class="form-label">Public Key</label>
            <input name="token" type="text" bind:value={publicKey} class="form-input full-width">

            <div class="buttons">
              <div class="col">
                <Button icon="fas fa-paper-plane" on:click={updatePublicKey} fullWidth="{true}"
                        disabled="{bot.id === undefined}">Submit Key
                </Button>
              </div>
              <div class="col">
                <Button icon="fas fa-paper-plane" on:click={createSlashCommands} fullWidth="{true}"
                        disabled="{!publicKeyOk}">Create Slash Commands
                </Button>
              </div>
            </div>
          </form>
        </div>
      </Card>
    </div>
  </div>
  <div class="content">
    <div class="content-col">
      <Card footer="{false}" fill="{false}">
        <h4 slot="title">Custom Status</h4>
        <div slot="body" class="full-width">
          <form class="full-width" onsubmit="return false;">
            <label class="form-label">Status</label>
            <input name="token" type="text" bind:value={bot.status} class="form-input full-width" placeholder="/help">

            <div class="buttons">
              <Button icon="fas fa-paper-plane" on:click={updateStatus} fullWidth="{true}"
                      disabled="{bot.id === undefined}">Submit
              </Button>
            </div>
          </form>
        </div>
      </Card>
    </div>
    <div class="content-col">
      <Card footer="{false}" fill="{false}">
        <h4 slot="title">Error Log</h4>
        <div slot="body" class="full-width">
          <table class="error-log">
            <thead>
            <tr style="border-bottom: 1px solid #dee2e6;">
              <th class="table-col">Error</th>
              <th class="table-col">Time</th>
            </tr>
            </thead>
            <tbody id="error_body">
            {#each errors as error}
              <tr class="table-row table-border">
                <td class="table-col">{error.message}</td>
                <td class="table-col">{error.time.toLocaleString()}</td>
              </tr>
            {/each}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  </div>
</div>

<style>
    .wrapper {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
        align-items: center;
    }

    .content {
        display: flex;
        justify-content: space-around;
        flex-direction: row;
        width: 95%;

        margin-top: 2%;
    }

    .col {
        width: 48%;
        height: 100%;
    }

    .content-col {
        width: 48%;
        height: 100%;
    }

    @media only screen and (max-width: 900px) {
        .content {
            flex-direction: column;
        }

        .content-col {
            width: 100%;
            margin-top: 2%;
        }
    }

    /* TODO: Move to central stylesheet*/
    :global(.form-label) {
        font-size: 12px;
        margin-bottom: 5px;
        color: #9a9a9a;
        text-transform: uppercase;
    }

    :global(.form-input), :global(.form-input:focus-visible) {
        border-color: #2e3136 !important;
        background-color: #2e3136 !important;
        color: white !important;
        outline: none;
        border-radius: 4px;
        padding: 8px 12px;
        margin: 0 0 0.5em 0;
        height: 40px;
    }

    .full-width {
        width: 100%;
    }

    .buttons {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        margin-top: 12px;
    }

    .error-log {
        width: 100%;
        border-collapse: collapse;
    }

    .table-col {
        width: 50%;
        text-align: left;
        padding: 5px 10px;
    }

    .table-border {
        border-top: 1px solid #dee2e6;
    }
</style>

<script>
    import {notifyError, notifyRatelimit, notifySuccess, withLoadingScreen} from '../js/util'
    import axios from "axios";
    import Card from '../components/Card.svelte'
    import Button from '../components/Button.svelte'
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'

    setDefaultHeaders()

    let token;
    let publicKey;
    let publicKeyOk = false;
    let interactionUrl;
    let bot = {};
    let errors = [];

    async function invite() {
        const res = await axios.get(`${API_URL}/user/whitelabel/`);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        const inviteUrl = 'https://discord.com/oauth2/authorize?client_id=' + res.data.id + '&scope=bot+applications.commands&permissions=805825648';
        window.open(inviteUrl, '_blank');
    }

    async function submitToken() {
        const data = {
            token: token
        };

        const res = await axios.post(`${API_URL}/user/whitelabel/`, data);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        $: token = '';

        await loadInteractionUrl();
        await loadBot();
        notifySuccess(`Started tickets whitelabel on ${res.data.bot.username}#${res.data.bot.discriminator}`);
    }

    async function updatePublicKey() {
        const data = {
            public_key: publicKey,
        };

        const res = await axios.post(`${API_URL}/user/whitelabel/public-key`, data);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        $: publicKeyOk = true;

        notifySuccess('Updated slash command settings successfully')
    }

    async function updateStatus() {
        const data = {
            status: bot.status,
        };

        const res = await axios.post(`${API_URL}/user/whitelabel/status`, data);
        if (res.status !== 200 || !res.data.success) {
            if (res.status === 429) {
                notifyRatelimit()
            } else {
                notifyError(res.data.error)
            }

            return;
        }

        notifySuccess('Updated status successfully')
    }

    async function loadBot() {
        const res = await axios.get(`${API_URL}/user/whitelabel/`);
        if (res.status !== 200 || !res.data.success) {
            if (res.status !== 404) {
                notifyError(res.data.error);
            }
            return;
        }

        bot = res.data;
    }

    async function loadErrors() {
        const res = await axios.get(`${API_URL}/user/whitelabel/errors`);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        // append errors
        if (res.data.errors !== null) {
            errors = res.data.errors.map(error => Object.assign({}, error, {time: new Date(error.time)}));
        }
    }

    async function loadPublicKey() {
        const res = await axios.get(`${API_URL}/user/whitelabel/public-key`);
        if (res.status === 404) {
            return;
        }

        if ((res.status !== 200 || !res.data.success)) {
            notifyError(res.data.error);
            return;
        }

        publicKey = res.data.key;
        $: publicKeyOk = true;
    }

    async function loadInteractionUrl() {
        const res = await axios.get(`${API_URL}/user/whitelabel/`);
        if (res.status === 404) {
            return;
        }

        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        interactionUrl = 'https://gateway.ticketsbot.net/handle/' + res.data.id;
    }

    async function createSlashCommands() {
        const opts = {
            timeout: 20 * 1000
        };

        const res = await axios.post(`${API_URL}/user/whitelabel/create-interactions`, {}, opts);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess('Slash commands have been created. Please note, Discord may take up to an hour to show them in your client');
    }

    async function checkPremium() {
        const res = await axios.get(`${API_URL}/api/${guildId}/premium`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        let isWhitelabel = res.data.tier >= 1;
        if (!isWhitelabel) {
            window.location.replace("https://ticketsbot.net/premium");
        }

        return isWhitelabel;
    }

    withLoadingScreen(async () => {
        if (await checkPremium()) {
            await loadBot();
            await loadErrors();
            await loadInteractionUrl();
            await loadPublicKey();
        }
    });
</script>
