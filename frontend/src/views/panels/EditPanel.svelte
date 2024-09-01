<main>
    <a href="/manage/{guildId}/panels" class="link">
        <i class="fas fa-arrow-left"></i>
        Back to Panels
    </a>
    <Card footer="{false}">
        <span slot="title">Create Panel</span>

        <div slot="body" class="body-wrapper">
            {#if !$loadingScreen}
                <PanelCreationForm {guildId} {channels} {roles} {emojis} {teams} {forms} {isPremium}
                                   bind:data={panelData} seedDefault={false} />
                <div class="submit-wrapper">
                    <Button icon="fas fa-floppy-disk" fullWidth={true} on:click={editPanel}>Save</Button>
                </div>
            {/if}
        </div>
    </Card>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        width: 100%;
        row-gap: 1vh;
    }

    main > a {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 18px;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
    }

    .submit-wrapper {
        margin: 1vh auto auto;
        width: 30%;
    }
</style>

<script>
    import {loadingScreen} from "../../js/stores";
    import Button from "../../components/Button.svelte";
    import Card from "../../components/Card.svelte";
    import PanelCreationForm from "../../components/manage/PanelCreationForm.svelte";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {notifyError, notifySuccess, setBlankStringsToNull, withLoadingScreen} from "../../js/util";
    import {onMount} from "svelte";
    import {loadChannels, loadEmojis, loadForms, loadPanels, loadPremium, loadRoles, loadTeams} from "../../js/common";
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {Navigate, navigateTo} from "svelte-router-spa";

    setDefaultHeaders();

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let panelId = parseInt(currentRoute.namedParams.panelid);

    let channels = [];
    let roles = [];
    let emojis = [];
    let teams = [];
    let forms = [];
    let isPremium = false;

    let panelData;

    async function editPanel() {
        setBlankStringsToNull(panelData);

        const res = await axios.patch(`${API_URL}/api/${guildId}/panels/${panelId}`, panelData);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/panels?edited=true`);
    }

    onMount(async () => {
        await withLoadingScreen(async () => {
            let panels = [];

            await Promise.all([
                loadChannels(guildId).then(r => channels = r).catch(e => notifyError(e)),
                loadRoles(guildId).then(r => roles = r).catch(e => notifyError(e)),
                loadEmojis(guildId).then(r => emojis = r).catch(e => notifyError(e)),
                loadTeams(guildId).then(r => teams = r).catch(e => notifyError(e)),
                loadForms(guildId).then(r => forms = r).catch(e => notifyError(e)),
                loadPremium(guildId, false).then(r => isPremium = r).catch(e => notifyError(e)),
                loadPanels(guildId).then(r => panels = r).catch(e => notifyError(e))
            ]);

            panelData = panels.find(p => p.panel_id === panelId);
            if (!panelData) {
                navigateTo(`/manage/${guildId}/panels?notfound=true`);
            }
        });
    });
</script>