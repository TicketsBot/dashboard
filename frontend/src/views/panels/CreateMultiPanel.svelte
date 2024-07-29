<main>
    <a href="/manage/{guildId}/panels" class="link">
        <i class="fas fa-arrow-left"></i>
        Back to Panels
    </a>
    <Card footer={false}>
        <span slot="title">Create Multi-Panel</span>
        <div slot="body" class="card-body">
            <p>Note: The panels which you wish to combine into a multi-panel must already exist</p>

            {#if !$loadingScreen}
                <div style="margin-top: 10px">
                    <MultiPanelCreationForm {guildId} {channels} {panels} bind:data={multiPanelCreateData}/>

                    <div class="submit-wrapper">
                        <Button icon="fas fa-paper-plane" fullWidth={true} on:click={createMultiPanel}>Create
                        </Button>
                    </div>
                </div>
            {/if}
        </div>
    </Card>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        padding: 2% 10% 4% 10%;
        width: 100%;
        row-gap: 1vh;
    }

    main > a {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 18px;
    }

    .card-body {
        display: flex;
        flex-direction: column;
        width: 100%;
    }

    .submit-wrapper {
        margin: 1vh auto auto;
        width: 30%;
    }
</style>

<script>
    import {loadingScreen} from "../../js/stores";
    import MultiPanelCreationForm from "../../components/manage/MultiPanelCreationForm.svelte";
    import Button from "../../components/Button.svelte";
    import Card from "../../components/Card.svelte";
    import {onMount} from "svelte";
    import {notifyError, removeBlankEmbedFields, setBlankStringsToNull, withLoadingScreen} from "../../js/util";
    import {loadChannels, loadPanels} from "../../js/common";
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {navigateTo} from "svelte-router-spa";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let channels = [];
    let panels = [];

    let multiPanelCreateData;

    async function createMultiPanel() {
        const data = structuredClone(multiPanelCreateData);

        setBlankStringsToNull(data);
        removeBlankEmbedFields(data);

        const res = await axios.post(`${API_URL}/api/${guildId}/multipanels`, data);
        if (res.status !== 200) {
            notifyError(res.data.error);
        } else {
            navigateTo(`/manage/${guildId}/panels?created=true`)
        }
    }

    onMount(async () => {
        await withLoadingScreen(async () => {
            await Promise.all([
                loadChannels(guildId).then(r => channels = r).catch(e => notifyError(e)),
                loadPanels(guildId).then(r => panels = r).catch(e => notifyError(e)),
            ])
        });
    });
</script>