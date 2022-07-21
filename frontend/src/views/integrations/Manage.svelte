<div class="parent">
  <div class="content">
    <div class="col-left">
      <Card footer footerRight>
        <span slot="title">Edit {integration.name} Settings</span>
        <div slot="body" class="body-wrapper">
          <h3>Secrets</h3>
          You cannot view previously submitted secrets for security reasons. The secret fields will show as empty,
          even if you have previously submitted them.

          {#if integration.secrets !== undefined}
            {#if integration.secrets.length === 0}
              <p>This integration does not require any secrets.</p>
            {:else}
              <p>This integration requires you to provide some secrets. These will be sent to the server controlled by
                the creator of {integration.name}, at: <code>{integration.webhook_url}</code></p>
              <p>Note, the integration creator may change the server at any time.</p>

              <div class="secret-container">
                {#each integration.secrets as secret}
                  <div class="secret-input">
                    <Input label="{secret.name}" placeholder="{secret.name}" bind:value={secretValues[secret.name]}/>
                    <p>{secret.description}</p>
                  </div>
                {/each}
              </div>
            {/if}
          {/if}
        </div>
        <div slot="footer" style="gap: 12px">
          <Button danger icon="fas fa-trash-can" on:click={removeIntegration}>Remove from server</Button>
          <Button icon="fas fa-floppy-disk" disabled={!allValuesFilled} on:click={saveSecrets}>Save Integration</Button>
        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import {notifyError, notifySuccess, withLoadingScreen} from '../../js/util'
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Card from "../../components/Card.svelte";
    import Button from "../../components/Button.svelte";
    import Input from "../../components/form/Input.svelte";
    import {navigateTo} from "svelte-router-spa";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let integrationId = currentRoute.namedParams.integration;

    let integration = {};
    let secretValues = {};

    let allValuesFilled = true;
    $: secretValues, updateAllValuesFilled();

    function updateAllValuesFilled() {
        if (integration.secrets === undefined) {
            return;
        }

        if (Object.keys(secretValues).length !== integration.secrets.length) {
            allValuesFilled = false;
            return;
        }

        for (let key in secretValues) {
            if (secretValues[key] === '') {
                allValuesFilled = false;
                return;
            }
        }

        allValuesFilled = true;
    }

    async function saveSecrets() {
        let data = {
            secrets: secretValues
        };

        const res = await axios.patch(`${API_URL}/api/${guildId}/integrations/${integrationId}`, data);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess('Integration settings saved successfully');
    }

    async function removeIntegration() {
        const res = await axios.delete(`${API_URL}/api/${guildId}/integrations/${integrationId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/integrations?removed=true`);
    }

    async function loadIntegration() {
        let res = await axios.get(`${API_URL}/api/integrations/view/${integrationId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        integration = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();

        await Promise.all([
            loadIntegration()
        ]);

        updateAllValuesFilled();
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
        flex-direction: row;
        justify-content: center;
        width: 96%;
        height: 100%;
        margin-top: 30px;
        padding-bottom: 5vh;

        gap: 2%;
    }

    .col-left {
        width: 60%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: 1vh;
    }

    .secret-container {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 2%;
    }

    .secret-input {
        flex: 0 0 49%;
    }

    @media only screen and (max-width: 950px) {
        .content {
            flex-direction: column;
            row-gap: 2vh;
        }

        .col-left {
            width: 100%;
        }
    }
</style>
