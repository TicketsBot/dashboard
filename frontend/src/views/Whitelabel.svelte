<div class="wrapper">
    <div class="col">
        {#if active}
            <Card footer="{false}" fill="{false}">
                <h4 slot="title">Manage Bot</h4>
                <div slot="body" class="full-width">
                    <p>Your whitelabel bot <b>{bot.username}</b> is active.</p>

                    <div class="buttons">
                        <Button icon="fas fa-plus" on:click={invite}>
                            Generate Invite Link
                        </Button>

                        <Button icon="fas fa-paper-plane" on:click={createSlashCommands}>
                            Re-create Slash Commands
                        </Button>

                        <Button icon="fas fa-trash-can" on:click={disable} danger>
                            Disable Whitelabel
                        </Button>
                    </div>
                </div>
            </Card>

            <Card footer="{false}" fill="{false}">
                <h4 slot="title">Custom Status</h4>
                <div slot="body" class="full-width">
                    <form class="form-wrapper full-width" on:submit|preventDefault>
                        <div class="row">
                            <Dropdown col3 label="Status Type" bind:value={bot.status_type}>
                                <option value="0">Playing</option>
                                <option value="2">Listening</option>
                                <option value="3">Watching</option>
                            </Dropdown>

                            <div class="col-2-3">
                                <Input col1 label="Status Text" placeholder="/help" bind:value={bot.status}/>
                            </div>
                        </div>

                        <div class="buttons">
                            <Button icon="fas fa-paper-plane" on:click={updateStatus} fullWidth="{true}">
                                Submit
                            </Button>
                        </div>
                    </form>
                </div>
            </Card>
        {:else}
            <Card footer="{false}" fill="{false}">
                <h4 slot="title">Bot Token</h4>
                <div slot="body" class="full-width">
                    <form class="full-width" on:submit|preventDefault>
                        <label class="form-label">Bot Token</label>

                        <input name="token" type="text" bind:value={token} class="form-input full-width"
                               placeholder="xxxxxxxxxxxxxxxxxxxxxxxx.xxxxxx.xxxxxxxxxxxxxxxxxxxxxxxxxxx">
                        <p>Note: You will not be able to view the token after submitting it</p>

                        <div class="buttons">
                            <Button icon="fas fa-paper-plane" on:click={submitToken} fullWidth="{true}">Submit
                            </Button>
                        </div>
                    </form>
                </div>
            </Card>
        {/if}
    </div>
    <div class="col">
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

<style>
    .wrapper {
        display: flex;
        flex-direction: row;
        height: 100%;
        width: 100%;
        padding: 2%;
        gap: 2%;
    }

    .col {
        display: flex;
        flex-direction: column;
        width: 50%;
        gap: 2%;
    }

    @media only screen and (max-width: 1180px) {
        .wrapper {
            flex-direction: column;
        }

        .col {
            width: 100%;
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
        gap: 12px;
        margin-top: 12px;
    }

    @media only screen and (max-width: 576px) {
        .buttons {
            flex-direction: column;
        }
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

    .form-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        gap: 10px;
    }
</style>

<script>
    import {notifyError, notifyRatelimit, notifySuccess, withLoadingScreen} from '../js/util'
    import axios from "axios";
    import Card from '../components/Card.svelte'
    import Button from '../components/Button.svelte'
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Dropdown from "../components/form/Dropdown.svelte";
    import Input from "../components/form/Input.svelte";

    setDefaultHeaders()

    let active = false;
    let token;
    let bot = {};
    let errors = [];

    async function invite() {
        const res = await axios.get(`${API_URL}/user/whitelabel/`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        const inviteUrl = 'https://discord.com/oauth2/authorize?client_id=' + res.data.id + '&scope=bot+applications.commands&permissions=805825784';
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

        await loadBot();
        notifySuccess(`Started tickets whitelabel on ${res.data.bot.name}`);
    }

    async function updateStatus() {
        const data = {
            status: bot.status,
            status_type: bot.status_type,
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
        if (res.status !== 200) {
            if (res.status === 402) {
                window.location.replace("https://ticketsbot.net/premium");
                return false;
            }

            if (res.status !== 404) {
                notifyError(res.data.error);
            }

            return true;
        }

        bot = res.data;

        active = true;

        return true;
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

    async function createSlashCommands() {
        const opts = {
            timeout: 20 * 1000
        };

        const res = await axios.post(`${API_URL}/user/whitelabel/create-interactions`, {}, opts);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess('Slash commands have been created. Please note, they may take a few minutes before they are visible.');
    }

    async function disable() {
        const res = await axios.delete(`${API_URL}/user/whitelabel/`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        active = false;
        notifySuccess('Whitelabel has been disabled');
    }

    withLoadingScreen(async () => {
        if (await loadBot()) {
            await Promise.all([
                loadErrors()
            ]);
        }
    });
</script>
