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
            <Input col3=true label="Ticket ID" placeholder="Ticket ID" bind:value={filterSettings.ticketId}/>
            <Input col3=true label="Username" placeholder="Username" bind:value={filterSettings.username}/>
            <Input col3=true label="User ID" placeholder="User ID" bind:value={filterSettings.userId}/>
          </div>
          <div class="row centre">
            <Button icon="fas fa-search">Filter</Button>
          </div>
        </form>
      </div>
    </Card>

    <div style="margin-top: 2%">
      <Card footer="{false}">
        <span slot="title">
          Transcripts
        </span>

        <div slot="body" class="body-wrapper">
          <table>
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
            {/each}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import Card from '../components/Card.svelte'
    import Input from '../components/form/Input.svelte'
    import Button from '../components/Button.svelte'

    import {withLoadingScreen} from '../js/util'
    import {onMount} from "svelte";
    import {dropdown} from "../js/stores";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id

    let filterSettings = {};
    let transcripts = [];

    async function filter() {
        alert('aaaaaaaaaaaaaaaaaaaaa')
    }

    onMount(async () => {
        dropdown.set(false)
    })

    withLoadingScreen(async () => {

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

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }

        :global([ref=filter-card]) {
            min-height: 252px !important;
        }
    }
</style>
