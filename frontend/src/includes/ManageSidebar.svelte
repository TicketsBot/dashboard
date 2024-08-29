<section class="sidebar">
    <header>
        <img src="{iconUrl}" class="guild-icon" alt="Guild icon" width="50" height="50"/>
        {guild.name}
    </header>
    <nav>
        <ul class="nav-list">
            <ManageSidebarLink {currentRoute} title="← Back to servers" href="/" />

            {#if isAdmin}
                <ManageSidebarLink {currentRoute} title="Settings" icon="fa-cogs" href="/manage/{guildId}/settings" />
            {/if}

            <ManageSidebarLink {currentRoute} title="Transcripts" icon="fa-copy" href="/manage/{guildId}/transcripts" />

            {#if isAdmin}
                <ManageSidebarLink {currentRoute} routePrefix="/manage/{guildId}/panels" title="Ticket Panels" icon="fa-mouse-pointer" href="/manage/{guildId}/panels">

                    <SubNavigation>
                        <SubNavigationLink {currentRoute} href="/manage/{guildId}/panels" routePrefix="/manage/{guildId}/panels">blah</SubNavigationLink>
                        <SubNavigationLink {currentRoute}>blahhh</SubNavigationLink>
                    </SubNavigation>
                </ManageSidebarLink>

                <ManageSidebarLink {currentRoute} title="Forms" icon="fa-poll-h" href="/manage/{guildId}/forms" />
                <ManageSidebarLink {currentRoute} title="Staff Teams" icon="fa-users" href="/manage/{guildId}/teams" />
                <ManageSidebarLink {currentRoute} title="Integrations" icon="fa-robot" href="/manage/{guildId}/integrations" />
            {/if}

            <ManageSidebarLink {currentRoute} title="Tickets" icon="fa-ticket-alt" href="/manage/{guildId}/tickets" />
            <ManageSidebarLink {currentRoute} title="Blacklist" icon="fa-ban" href="/manage/{guildId}/blacklist" />
            <ManageSidebarLink {currentRoute} title="Tags" icon="fa-tags" href="/manage/{guildId}/tags" />
        </ul>
    </nav>
    <nav class="bottom">
        <hr/>
        <ul class="nav-list">
            <ManageSidebarLink {currentRoute} title="Documentation" icon="fa-book" href="https://docs.ticketsbot.net" newWindow />
            <ManageSidebarLink {currentRoute} title="Logout" icon="fa-sign-out-alt" href="/logout" />
        </ul>
    </nav>
</section>

<style>
    .sidebar {
        display: flex;
        flex-direction: column;
        align-self: flex-start;
        background-color: #272727;
        padding: 15px;
        width: 275px;
        border-radius: 6px;
        user-select: none;
    }

    header {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 10px;

        font-weight: bold;

        padding: 6px 10px;
        border-radius: 4px;

        background: linear-gradient(33.3deg, #873ef5 0%, #995DF3 100%);
        box-shadow: 0 6px 6px rgba(10, 10, 10, .1), 0 0 0 1px rgba(10, 10, 10, .1);
    }

    .guild-icon {
        width: 48px;
        height: 48px;
        border-radius: 50%;
    }

    nav > ul {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }

    nav hr {
        width: 40%;
        padding-left: 20px;
    }
</style>

<script>
    import {onMount} from "svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {notifyError, withLoadingScreen} from "../js/util";
    import ManageSidebarLink from "./ManageSidebarLink.svelte";
    import SubNavigation from "./SubNavigation.svelte";
    import SubNavigationLink from "./SubNavigationLink.svelte";

    export let currentRoute;
    export let permissionLevel;

    $: isAdmin = permissionLevel >= 2;

    let guildId = currentRoute.namedParams.id;

    let guild = {};
    let iconUrl = "";

    async function loadGuild() {
        const res = await axios.get(`${API_URL}/api/${guildId}/guild`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        guild = res.data;
    }

    function isAnimated() {
        if (guild.icon === undefined || guild.icon === "") {
            return false;
        } else {
            return guild.icon.startsWith('a_')
        }
    }

    function getIconUrl() {
        if (!guild.icon) {
            return `https://cdn.discordapp.com/embed/avatars/${Number((BigInt(guildId) >> BigInt(22)) % BigInt(6))}.png`
        }

        if (isAnimated()) {
            return `https:\/\/cdn.discordapp.com/icons/${guild.id}/${guild.icon}.gif?size=256`
        } else {
            return `https:\/\/cdn.discordapp.com/icons/${guild.id}/${guild.icon}.webp?size=256`
        }
    }

    onMount(async () => {
        await withLoadingScreen(async () => {
            await loadGuild();

            iconUrl = getIconUrl();
        })
    });
</script>