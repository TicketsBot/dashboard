<div class="guild-badge" on:click={goto(guild.id)} class:disabled={guild.permission_level === 0}>
    <div class="guild-icon-bg">
        {#if guild.icon === undefined || guild.icon === ""}
            <i class="fas fa-question guild-icon-fa" class:disabled={guild.permission_level === 0}></i>
        {:else}
            <img class="guild-icon" src="{getIconUrl()}" alt="Guild Icon"
                 class:disabled={guild.permission_level === 0}/>
        {/if}
    </div>

    <div class="text-wrapper" class:disabled={guild.permission_level === 0}>
        <span class="guild-name">
          {guild.name}
        </span>
        <span class="no-permission" class:disabled={guild.permission_level > 0}>
        No permission
        <Tooltip tip="You do not have permission to manage this server." top color="#121212">
            <a href="https://docs.ticketsbot.net/miscellaneous/dashboard-no-permission" target="_blank">
                <i class="fas fa-circle-question form-label tooltip-icon"></i>
            </a>
        </Tooltip>
    </span>
    </div>
</div>

<script>
    import Tooltip from "svelte-tooltip";

    export let guild;

    function isAnimated() {
        if (guild.icon === undefined || guild.icon === "") {
            return false;
        } else {
            return guild.icon.startsWith('a_')
        }
    }

    function getIconUrl() {
        if (isAnimated()) {
            return `https:\/\/cdn.discordapp.com/icons/${guild.id}/${guild.icon}.gif?size=256`
        } else {
            return `https:\/\/cdn.discordapp.com/icons/${guild.id}/${guild.icon}.webp?size=256`
        }
    }

    async function goto(guildId) {
        if (guild.permission_level === 2) {
            window.location.href = `/manage/${guildId}/settings`;
        } else if (guild.permission_level === 1) {
            window.location.href = `/manage/${guildId}/transcripts`;
        } else {
            return;
        }
    }
</script>

<style>
    :global(.guild-badge) {
        display: flex;
        align-items: center;
        box-shadow: 0 4px 4px rgba(0, 0, 0, 0.25);

        width: 33%;
        background-color: #121212;
        height: 100px;
        margin-bottom: 10px;
        border-radius: 10px;
        cursor: pointer;
    }

    .guild-badge.disabled {
        cursor: default;
    }

    @media (max-width: 950px) {
        :global(.guild-badge) {
            width: 100%;
        }
    }

    :global(.guild-icon-bg) {
        height: 80px;
        width: 80px;
        background-color: #272727;
        border-radius: 50%;
        margin-left: 10px;
    }

    :global(.guild-icon) {
        height: 80px;
        width: 80px;
        border-radius: 50%;
    }

    :global(.guild-icon-fa) {
        border-radius: 50%;
        color: white;
        font-size: 60px !important;
        width: 80px;
        height: 80px;
        text-align: center;
        margin-top: 10px;
    }

    :global(.guild-name) {
        color: white !important;
    }

    .text-wrapper.disabled > .guild-name {
        opacity: 45%;
    }

    .guild-icon-bg > *.disabled {
        opacity: 25%;
    }

    .text-wrapper {
        display: flex;
        flex-direction: column;
        padding-left: 10px;
    }

    .text-wrapper > .no-permission {
        opacity: 75%;
    }

    .text-wrapper > .no-permission.disabled {
        visibility: hidden;
    }
</style>
