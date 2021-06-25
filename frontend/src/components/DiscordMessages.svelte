<div class="discord-container">
  <div class="channel-header">
    <span id="channel-name">#ticket-{ticketId}</span>
  </div>
  <div id="message-container">
    {#each messages as message}
      <div class="message">
        <b>{message.username}:</b> {message.content}
      </div>
    {/each}
  </div>
  <div class="input-container">
    <form on:submit|preventDefault={sendMessage}>
      <input type="text" class="message-input" bind:value={sendContent}>
    </form>
  </div>
</div>

<script>
    import {createEventDispatcher} from "svelte";

    export let ticketId;
    export let messages = [];

    const dispatch = createEventDispatcher();
    let sendContent = '';

    function sendMessage() {
        dispatch('send', sendContent);
        sendContent = '';
    }
</script>

<style>
    @font-face {
        font-family: Catamaran;
        font-weight: 300;
        src: url('/assets/fonts/Catamaran-Light.ttf') format('truetype')
    }

    @font-face {
        font-family: Catamaran;
        font-weight: 500;
        src: url('/assets/fonts/Catamaran-Regular.ttf') format('truetype')
    }

    @font-face {
        font-family: Catamaran;
        font-weight: 600;
        src: url('/assets/fonts/Catamaran-Medium.ttf') format('truetype')
    }

    @font-face {
        font-family: Catamaran;
        font-weight: 700;
        src: url('/assets/fonts/Catamaran-Bold.ttf') format('truetype')
    }

    .discord-container {
        background-color: #2e3136;
        border-radius: 4px;
        height: 80vh;
        max-height: 100vh;
        margin: 0;
        padding: 0;
        font-family: 'Catamaran', sans-serif !important;
    }

    .channel-header {
        background-color: #1e2124;
        height: 5vh;
        width: 100%;
        border-radius: 4px 4px 0 0;
        position: relative;

        text-align: center;
        display: flex;
        align-items: center;
    }

    #channel-name {
        color: white;
        padding-left: 20px;
    }

    #message-container {
        height: 70vh;
        max-height: 70vh;
        position: relative;
        overflow: scroll;
        overflow-x: hidden;
    }

    .message {
        color: white !important;
        padding-left: 20px;
    }

    .message-input {
        padding-top: 20px !important;
        border-color: #2e3136 !important;
        padding-left: 5px !important;
        padding-right: 5px !important;
        background-color: #2e3136 !important;
        color: white !important;
        height: 100%;
        max-height: 100%;
        min-height: 100%;
    }

    .message-input:focus {
        border-color: #2e3136 !important;
        background-color: #2e3136 !important;
        color: white !important;
    }
</style>