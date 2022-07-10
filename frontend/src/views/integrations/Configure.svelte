{#if data !== undefined}
  <IntegrationEditor editMode {guildId} {data} on:delete={deleteIntegration}
                     on:submit={(e) => editIntegration(e.detail)} on:makePublic={(e) => makePublic(e.detail)}/>
{/if}

<script>
    import IntegrationEditor from "../../components/IntegrationEditor.svelte";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {notifyError, notifySuccess, withLoadingScreen} from "../../js/util";
    import axios from "axios";
    import {navigateTo} from "svelte-router-spa";
    import {API_URL} from "../../js/constants";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let integrationId = currentRoute.namedParams.integration;
    let freshlyCreated = currentRoute.queryParams.created === "true";

    let data;

    async function makePublic(data) {
        const res = await axios.post(`${API_URL}/api/integrations/${integrationId}/public`, data);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess("Your request to make this integration public has been submitted! It will be reviewed over the next few days.");
        data.public = true;
    }

    async function editIntegration(data) {
        const res = await axios.patch(`${API_URL}/api/integrations/${integrationId}`, data);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess("Integration updated");
        await loadIntegration();
    }

    async function deleteIntegration() {
        const res = await axios.delete(`${API_URL}/api/integrations/${integrationId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/integrations`);
    }

    async function loadIntegration() {
        const res = await axios.get(`${API_URL}/api/integrations/view/${integrationId}/detail`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadIntegration();

        if (freshlyCreated) {
            notifySuccess("Your integration has been created successfully!");
        }
    });
</script>