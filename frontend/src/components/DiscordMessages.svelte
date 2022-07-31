<div class="discord-container">
  <div class="channel-header">
    <span id="channel-name">#ticket-{ticketId}</span>
  </div>
  <div id="message-container" bind:this={container}>
    {#each messages as message}
      <div class="message">
        <b>{message.author.username}:</b> {message.content}
      </div>
    {/each}
  </div>
  <div class="input-container">
    <form on:submit|preventDefault={sendMessage}>
      <input type="text" class="message-input" bind:value={sendContent} disabled={!isPremium}
             placeholder="{isPremium ? `Message #ticket-${ticketId}` : 'Premium users can receive messages in real-time and respond to tickets through the dashboard'}">
    </form>
  </div>
</div>

<script>
    import {createEventDispatcher} from "svelte";

    export let ticketId;
    export let isPremium = false;
    export let messages = [];
    export let container;

    const dispatch = createEventDispatcher();
    let sendContent = '';

    function sendMessage() {
        dispatch('send', sendContent);
        sendContent = '';
    }
</script>

<style>
    @import url('https://fonts.googleapis.com/css2?family=Catamaran:wght@300;400;500;600;700;800&display=swap');

    .discord-container {
        display: flex;
        flex-direction: column;

        background-color: #2e3136;
        border-radius: 4px;
        height: 80vh;
        max-height: 100vh;
        margin: 0;
        padding: 0;
        font-family: 'Catamaran', sans-serif !important;
    }

    .channel-header {
        display: flex;
        align-items: center;

        background-color: #1e2124;
        height: 5vh;
        width: 100%;
        border-radius: 4px 4px 0 0;
        position: relative;

        text-align: center;
    }

    #channel-name {
        color: white;
        font-weight: bold;
        padding-left: 20px;
    }

    #message-container {
        display: flex;
        flex-direction: column;
        flex: 1;

        position: relative;
        overflow-y: scroll;
        overflow-wrap: break-word;
    }

    .message {
        color: white !important;
        padding-left: 20px;
    }

    #message-container:last-child {
        margin-bottom: 5px;
    }

    .message-input {
        display: flex;

        font-size: 16px;
        line-height: 24px;
        height: 40px;
        padding: 8px;

        border-color: #2e3136 !important;
        background-color: #2e3136 !important;
        color: white !important;
        width: 100%;
    }

    .message-input:focus, .message-input:focus-visible {
        outline-width: 0;
    }
</style>
