<div class="wrapper">
  <div class="col-main">
    <div class="row">
      <Card footer="{false}">
        <span slot="title">Reaction Panels</span>
        <div slot="body" class="card-body">
          <p>Your panel quota: <b>{panels.length} / {isPremium ? '∞' : '3'}</b></p>

          <table style="margin-top: 10px">
            <thead>
            <tr>
              <th>Channel</th>
              <th>Panel Title</th>
              <th class="category-col">Ticket Channel Category</th>
              <th>Edit</th>
              <th>Delete</th>
            </tr>
            </thead>
            <tbody>
            {#each panels as panel}
              <tr>
                <td>#{channels.find((c) => c.id === panel.channel_id)?.name ?? 'Unknown Channel'}</td>
                <td>{panel.title}</td>
                <td class="category-col">{channels.find((c) => c.id === panel.category_id)?.name ?? 'Unknown Category'}</td>
                <td>
                  <Button on:click={() => editPanel(panel.panel_id)}>Edit</Button>
                </td>
                <td>
                  <Button on:click={() => deletePanel(panel.panel_id)}>Delete</Button>
                </td>
              </tr>
            {/each}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
    <div class="row">
      <Card footer="{false}">
        <span slot="title">Create Panel</span>

        <div slot="body" class="body-wrapper">
          <PanelCreationForm guildId={guildId} bind:data={panelCreateData}
                             on:submit={() => console.log(panelCreateData)}/>
        </div>
      </Card>
    </div>
  </div>
  <div class="col-small">
    <div class="row">
      <Card footer="{false}">
        <span slot="title">Multi-Panels</span>
        <div slot="body">

        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, withLoadingScreen} from "../js/util";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Button from "../components/Button.svelte";
    import PanelCreationForm from "../components/manage/PanelCreationForm.svelte";

    export let currentRoute;
    export let params = {};

    setDefaultHeaders()

    let guildId = currentRoute.namedParams.id;

    let channels = [];
    let panels = [];
    let isPremium = false;

    let panelCreateData;

    async function editPanel(panelId) {

    }

    async function deletePanel(panelId) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/panels/${panelId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        panels = panels.filter((p) => p.panel_id !== panelId);
    }

    async function loadPremium() {
        const res = await axios.get(`${API_URL}/api/${guildId}/premium`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        isPremium = res.data.premium;
    }

    async function loadChannels() {
        const res = await axios.get(`${API_URL}/api/${guildId}/channels`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        channels = res.data;
    }

    async function loadPanels() {
        const res = await axios.get(`${API_URL}/api/${guildId}/panels`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        panels = res.data;
    }

    withLoadingScreen(async () => {
        await loadPremium();
        await loadChannels();
        await loadPanels();
    })
</script>

<style>
    .wrapper {
        display: flex;
        flex-direction: row;
        height: 100%;
        width: 100%;
        margin-top: 30px;
    }

    .body-wrapper {
        width: 100%;
    }

    .col-main {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 65%;
        height: 100%;
    }

    .col-small {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 35%;
        height: 100%;
    }

    .row {
        display: flex;
        width: 96%;
        height: 100%;
        margin-bottom: 2%;
    }

    .card-body {
        width: 100%;
    }

    @media only screen and (max-width: 1100px) {
        .wrapper {
            flex-direction: column;
        }

        .col-main, .col-small {
            width: 100%;
            margin-bottom: 4%;
        }
    }

    @media only screen and (max-width: 576px) {
        .category-col {
            display: none;
        }

        .row {
            width: 100%;
        }
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th {
        text-align: left;
        font-weight: normal;
        border-bottom: 1px solid #dee2e6;
        padding-left: 10px;
    }

    tr {
        border-bottom: 1px solid #dee2e6;
    }

    tr:last-child {
        border-bottom: none;
    }

    td {
        padding: 10px 0 10px 10px;
    }
</style>