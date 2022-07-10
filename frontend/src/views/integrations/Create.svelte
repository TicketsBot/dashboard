<IntegrationEditor {guildId} on:submit={(e) => createIntegration(e.detail)} />

<script>
    import IntegrationEditor from "../../components/IntegrationEditor.svelte";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {notifyError, withLoadingScreen} from "../../js/util";
    import axios from "axios";
    import {navigateTo} from "svelte-router-spa";
    import {API_URL} from "../../js/constants";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    async function createIntegration(data) {
        const res = await axios.post(`${API_URL}/api/integrations`, data);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/integrations/configure/${res.data.id}?created=true`);
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
    });
</script>