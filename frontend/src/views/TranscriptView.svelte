<div class="discord-container">
  <div>
    <div class="channel-header">
      <span class="channel-name">#ticket-{ticketId}</span>
    </div>
  </div>
  <div id="message-container">
    {#each messages as message}
      <div class="message">
        {#if message.timestamp > epoch}
          <span class="timestamp">
              [{message.timestamp.toLocaleTimeString([], dateFormatSettings)} {message.timestamp.toLocaleDateString()}]
          </span>
        {/if}

        <img src="https://cdn.discordapp.com/avatars/{message.author.id}/{message.author.avatar}.webp?size=256"
             class="avatar">
        <b class="username">{message.author.username}</b>
        <span class="content">{message.content}</span>

        {#if message.attachments !== undefined && message.attachments.length > 0}
          {#each message.attachments as attachment}
            <a href="{attachment.url}" target="_blank" title="{attachment.filename}" class="attachment"><i
                    class="far fa-file-alt fa-2x"></i></a>
          {/each}
        {/if}
      </div>
    {/each}
  </div>
</div>

<script>
    import axios from "axios";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {errorPage, withLoadingScreen} from '../js/util'
    import {API_URL} from "../js/constants";

    export let currentRoute;
    export let params = {};

    let guildId = currentRoute.namedParams.id;
    let ticketId = currentRoute.namedParams.ticketid;

    setDefaultHeaders()

    let messages = [];
    let epoch = new Date('2015');
    let dateFormatSettings = {
        hour: '2-digit',
        minute: '2-digit'
    };

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/${guildId}/transcripts/${ticketId}`);
        if (res.status !== 200) {
            errorPage(res.data.error);
            return;
        }

        messages = res.data.map(message => Object.assign({}, message, {timestamp: new Date(message.timestamp)}));
    }

    withLoadingScreen(loadData);
</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Noto+Sans&display=swap');

    body {
        margin: 0;
        display: flex;
        width: 100%;
        height: 100%;
        background-color: #2e3136;
    }

    .discord-container {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        overflow-y: scroll;

        background-color: #2e3136;
        font-family: 'Noto Sans', sans-serif !important;
        font-size: 16px;
    }

    .channel-header {
        background-color: #1e2124;
        height: 50px;
        width: 100%;

        text-align: center;
        display: flex;
        align-items: center;
    }

    .channel-name {
        color: white;
        padding-left: 20px;
    }

    #message-container {
        display: flex;
        flex-direction: column;
        flex: 1;
        height: 100%;
    }

    .message {
        display: flex;
        align-items: center;
        color: white !important;
        word-wrap: break-word;
        margin: 4px 0 4px 12px;
    }

    .username, .content, .avatar {
        margin-left: 4px;
    }

    .avatar {
        height: 32px;
        width: 32px;
        border-radius: 50%;
    }

    .attachment {
        margin-left: 12px;
        color: white;
    }

    @media only screen and (max-width: 576px) {
        .timestamp {
            display: none;
        }
    }
</style>