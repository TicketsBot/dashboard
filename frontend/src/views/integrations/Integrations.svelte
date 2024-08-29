<div class="content">
    <div class="container">
        <div class="spread">
            <h4 class="title">My Integrations</h4>
            <Button icon="fas fa-server" on:click={() => navigateTo(`/manage/${guildId}/integrations/create`)}>Create
                Integration
            </Button>
        </div>
        <div class="integrations my-integrations">
            {#each ownedIntegrations as integration}
                <div class="integration">
                    <Integration owned name={integration.name} {guildId} integrationId={integration.id}
                                 imageUrl={generateProxyUrl(integration)} guildCount={integration.guild_count}>
              <span slot="description">
                {integration.description}
              </span>
                    </Integration>
                </div>
            {/each}
        </div>
    </div>

    <div>
        <h4 class="title">Available Integrations</h4>
        <div class="integrations">
            <!-- Built in -->
            {#if page === 1}
                <div class="integration">
                    <Integration builtIn name="Bloxlink"
                                 imageUrl="https://dbl-static.b-cdn.net/9bbd1f9504ddefc89606b19b290e9a0f.png"
                                 viewLink="https://docs.ticketsbot.net/dashboard/settings/placeholders#bloxlink">
          <span slot="description">
            Our Bloxlink integration inserts the Roblox usernames, profile URLs and more of your users into
            ticket welcome messages automatically! This integration is automatically enabled in all servers, press the
            View button below to check out the full list of placeholders you can use!
          </span>
                    </Integration>
                </div>
            {/if}

            {#each availableIntegrations as integration}
                <div class="integration">
                    <Integration name={integration.name} {guildId} integrationId={integration.id}
                                 imageUrl={generateProxyUrl(integration)} ownerId={integration.owner_id}
                                 added={integration.added} guildCount={integration.guild_count} showAuthor
                                 author={integration.author} on:remove={() => removeIntegration(integration.id)}>
              <span slot="description">
                {integration.description}
              </span>
                    </Integration>
                </div>
            {/each}
        </div>
    </div>

    <div class="pagination">
        <i class="fas fa-chevron-left pagination-chevron" class:disabled-chevron={page === 1}
           on:click={previousPage}></i>
        <p>Page {page}</p>
        <i class="fas fa-chevron-right pagination-chevron" class:disabled-chevron={!hasNextPage}
           on:click={nextPage}></i>
    </div>
</div>

<script>
    import {notifyError, notifySuccess, withLoadingScreen} from '../../js/util'
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Integration from "../../components/manage/Integration.svelte";
    import Button from "../../components/Button.svelte";
    import {navigateTo} from "svelte-router-spa";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let freshlyRemoved = currentRoute.queryParams.removed === "true";
    let freshlyAdded = currentRoute.queryParams.added === "true";

    let ownedIntegrations = [];
    let availableIntegrations = [];

    const pageLimit = 20;
    const builtInIntegrationCount = 1;
    let page = 1;

    function previousPage() {
        if (page > 1) {
            page--;
            loadAvailableIntegrations();
        }
    }

    function nextPage() {
        if (hasNextPage) {
            page++;
            loadAvailableIntegrations();
        }
    }

    let hasNextPage = true;
    $: if (page === 1) {
        hasNextPage = availableIntegrations.length + builtInIntegrationCount >= pageLimit;
    } else {
        hasNextPage = availableIntegrations.length >= pageLimit;
    }

    function generateProxyUrl(integration) {
        if (integration.image_url === null || integration.proxy_token === undefined || integration.proxy_token === null) {
            return null;
        }

        return `https://image-cdn.ticketsbot.net/proxy?token=${integration.proxy_token}`
    }

    async function removeIntegration(integrationId) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/integrations/${integrationId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        await loadAvailableIntegrations();
        notifySuccess("Integration removed from server successfully");
    }

    async function loadAvailableIntegrations() {
        const res = await axios.get(`${API_URL}/api/${guildId}/integrations/available?page=${page}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return
        }

        availableIntegrations = res.data;
    }

    async function loadOwnedIntegrations() {
        let res = await axios.get(`${API_URL}/api/integrations/self`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        ownedIntegrations = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();

        await Promise.all([
            loadOwnedIntegrations(),
            loadAvailableIntegrations()
        ]);

        if (freshlyAdded) {
            notifySuccess("Integration added to server successfully!");
        } else if (freshlyRemoved) {
            notifySuccess("Integration removed from server successfully!");
        }
    });
</script>

<style>
    .content {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        row-gap: 4vh;
    }

    .container {
        display: flex;
        flex-direction: column;
    }

    .integrations {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 2%;
        row-gap: 4vh;
    }

    .integration {
        flex: 0 0 32%;
    }

    .my-integrations {
        margin-top: 2vh;
    }

    .title {
        color: white;
        font-size: 22px;
        font-weight: bolder;
    }

    .spread {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        gap: 10px;
    }

    .pagination {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        gap: 5px;
    }

    .pagination-chevron {
        cursor: pointer;
        color: #3472f7;
    }

    .disabled-chevron {
        color: #777 !important;
        cursor: default !important;
    }

    @media only screen and (max-width: 1200px) {
        .integration {
            flex: 0 0 49%;
        }
    }

    @media only screen and (max-width: 850px) {
        .integration {
            flex: 0 0 100%;
        }
    }
</style>
