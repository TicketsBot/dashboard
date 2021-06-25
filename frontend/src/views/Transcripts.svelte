<div class="content">
  <div class="col">
    <Card footer="{false}" dropdown="{true}" ref="filter-card">
      <span slot="title">
        <i class="fas fa-filter"></i>
        Filter Logs
      </span>

      <div slot="body" class="body-wrapper">
        <form class="form-wrapper" on:submit|preventDefault={filter}>
          <div class="row">
            <Input col3=true label="Ticket ID" placeholder="Ticket ID"
                   on:input={handleInputTicketId} bind:value={filterSettings.ticketId}/>

            <Input col3=true label="Username" placeholder="Username" on:input={handleInputUsername}
                   bind:value={filterSettings.username}/>

            <Input col3=true label="User ID" placeholder="User ID" on:input={handleInputUserId}
                   bind:value={filterSettings.userId}/>
          </div>
          <div class="row centre">
            <Button icon="fas fa-search">Filter</Button>
          </div>
        </form>
      </div>
    </Card>

    <div style="margin: 2% 0;">
      <Card footer="{false}">
        <span slot="title">
          Transcripts
        </span>

        <div slot="body" class="main-col">
          <table class="nice">
            <thead>
            <tr>
              <th>Ticket ID</th>
              <th>Username</th>
              <th>Close Reason</th>
              <th>Transcript</th>
            </tr>
            </thead>
            <tbody>
            {#each transcripts as transcript}
              <tr>
                <td>{transcript.ticket_id}</td>
                <td>{transcript.username}</td>
                <td>{transcript.close_reason || 'No reason specified'}</td>
                <td>
                  <Navigate to="{`/manage/${guildId}/transcripts/view/${transcript.ticket_id}`}" styles="link">
                    <Button>View</Button>
                  </Navigate>
                </td>
              </tr>
            {/each}
            </tbody>
          </table>

          <div class="nav" class:nav-margin={transcripts.length === 0}>
            <i class="fas fa-chevron-left" class:hidden={page === 1} on:click={loadPrevious}></i>
            <span>Page {page}</span>
            <i class="fas fa-chevron-right"
               class:hidden={transcripts.length < pageLimit || transcripts[transcripts.length - 1].ticket_id === 1}
               on:click={loadNext}></i>
          </div>
        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import Card from '../components/Card.svelte'
    import Input from '../components/form/Input.svelte'
    import Button from '../components/Button.svelte'

    import {notifyError, withLoadingScreen} from '../js/util'
    import {onMount} from "svelte";
    import {dropdown} from "../js/stores";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {Navigate} from 'svelte-router-spa'

    setDefaultHeaders();

    export let currentRoute;
    let guildId = currentRoute.namedParams.id

    let filterSettings = {};
    let transcripts = [];

    const pageLimit = 30;
    let page = 1;

    let handleInputTicketId = () => {
        filterSettings.username = undefined;
        filterSettings.userId = undefined;
    };

    let handleInputUsername = () => {
        filterSettings.ticketId = undefined;
        filterSettings.userId = undefined;
    };

    let handleInputUserId = () => {
        filterSettings.ticketId = undefined;
        filterSettings.username = undefined;
    };

    async function loadPrevious() {
        if (page === 1) {
            return;
        }

        let paginationSettings = {
            after: transcripts[0].ticket_id,
        };

        if (await loadData(paginationSettings)) {
            page--;
        }
    }

    async function loadNext() {
        if (transcripts.length < pageLimit || transcripts[transcripts.length - 1].ticket_id === 1) {
            return;
        }

        let paginationSettings = {
            before: transcripts[transcripts.length - 1].ticket_id,
        };

        if (await loadData(paginationSettings)) {
            page++;
        }
    }

    function buildQuery(paginationSettings) {
        let query = new URLSearchParams();
        if (paginationSettings['before'] !== undefined) {
            query.append('before', paginationSettings['before']);
        }

        if (paginationSettings['after'] !== undefined) {
            query.append('after', paginationSettings['after']);
        }

        if (filterSettings['ticketId'] !== undefined) {
            query.append('ticketid', filterSettings.ticketId);
        }

        if (filterSettings['username'] !== undefined) {
            query.append('username', filterSettings.username);
        }

        if (filterSettings['userId'] !== undefined) {
            query.append('userid', filterSettings.userId);
        }

        return query;
    }

    async function filter() {
        await loadData({});
        page = 1;
    }

    async function loadData(paginationSettings) {
        let query = buildQuery(paginationSettings);
        const res = await axios.get(`${API_URL}/api/${guildId}/transcripts?${query.toString()}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return false;
        }

        transcripts = res.data;
        return true;
    }

    onMount(async () => {
        dropdown.set(false)
    })

    withLoadingScreen(async () => {
        await loadData({})
    })
</script>

<style>
    .content {
        display: flex;
        justify-content: center;
        height: 100%;
        width: 100%;
    }

    .col {
        display: flex;
        flex-direction: column;
        width: 95%;
        height: 100%;
        margin-top: 30px;
    }

    .main-col {
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
    }

    .centre {
        justify-content: center !important;
    }

    .form-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    :global([ref=filter-card]) {
        min-height: 110px !important;
    }

    :global(table.nice) {
        width: 100%;
        border-collapse: collapse;
    }

    :global(table.nice > thead > tr > th) {
        text-align: left;
        font-weight: normal;
        border-bottom: 1px solid #dee2e6;
        padding-left: 10px;
    }

    :global(table.nice > thead > tr, table.nice > tbody > tr) {
        border-bottom: 1px solid #dee2e6;
    }

    :global(table.nice > tbody > tr:last-child) {
        border-bottom: none;
    }

    :global(table.nice > tbody > tr > td) {
        padding: 10px 0 10px 10px;
    }

    .nav {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
    }

    .nav > i {
        color: #1dc7ea;
        cursor: pointer;
    }

    .nav > span {
        margin: 0 5px;
    }

    .nav-margin {
        margin-top: 15px;
    }

    .hidden {
        color: #6c757d !important;
        cursor: default !important;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }

        :global([ref=filter-card]) {
            min-height: 252px !important;
        }
    }
</style>
