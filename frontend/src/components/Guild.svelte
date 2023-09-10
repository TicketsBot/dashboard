<div class="guild-badge" on:click={goto(guild.id)}>
  <div class="guild-icon-bg">
    {#if guild.icon === undefined || guild.icon === ""}
      <i class="fas fa-question guild-icon-fa"></i>
    {:else}
      <img class="guild-icon" src="{getIconUrl()}" alt="Guild Icon"/>
    {/if}
  </div>

  <div>
    <span class="guild-name">
      {guild.name}
    </span>
  </div>
</div>

<script>
    import axios from 'axios';
    import {API_URL} from "../js/constants";
    import {notifyError} from "../js/util";

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
        } else {
            window.location.href = `/manage/${guildId}/transcripts`;
        }
    }

    async function getPermissionLevel(guildId) {
        const res = await axios.get(`${API_URL}/user/permissionlevel?guild=${guildId}`);
        if (res.status !== 200 || !res.data.success) {
            notifyError(res.data.error);
            return;
        }

        return res.data.permission_level;
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
        padding-left: 10px;
    }
</style>
