{#if tagSelectorModal}
    <div class="modal" transition:fade>
        <div class="modal-wrapper">
            <Card footer footerRight fill={false}>
                <span slot="title">Send Tag</span>

                <div slot="body" class="modal-inner">
                    <Dropdown col2 label="Select a tag..." bind:value={selectedTag}>
                        {#each Object.keys(tags) as tag}
                            <option value={tag}>{tag}</option>
                        {/each}
                    </Dropdown>
                </div>

                <div slot="footer" style="gap: 12px">
                    <Button danger icon="fas fa-times" on:click={() => tagSelectorModal = false}>Close</Button>
                    <Button icon="fas fa-paper-plane" on:click={sendTag}>Send</Button>
                </div>
            </Card>
        </div>
    </div>

    <div class="modal-backdrop" transition:fade>
    </div>
{/if}

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
            {#if isPremium}
                <i class="fas fa-paper-plane send-button" on:click={sendMessage}/>
                <div class="tag-selector">
                    <Button type="button" noShadow on:click={openTagSelector}>Select Tag</Button>
                </div>
            {/if}
        </form>
    </div>
</div>

<script>
    import {createEventDispatcher} from "svelte";
    import {fade} from "svelte/transition";
    import Button from "./Button.svelte";
    import Card from "./Card.svelte";
    import Dropdown from "./form/Dropdown.svelte";

    export let ticketId;
    export let isPremium = false;
    export let messages = [];
    export let container;

    export let tags = [];

    const dispatch = createEventDispatcher();
    let sendContent = '';
    let selectedTag;

    function sendMessage() {
        dispatch('send', {
            type: 'message',
            content: sendContent
        });
        sendContent = '';
    }

    let tagSelectorModal = false;

    function openTagSelector() {
        tagSelectorModal = true;
        window.scrollTo({top: 0, behavior: 'smooth'});
    }

    function sendTag() {
        tagSelectorModal = false;

        dispatch('send', {
            type: 'tag',
            tag_id: selectedTag
        });

        selectedTag = undefined;
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
        flex: 1;

        font-size: 16px;
        line-height: 24px;
        height: 40px;
        padding: 8px;

        border-color: #2e3136 !important;
        background-color: #2e3136 !important;
        color: white !important;
    }

    .message-input:focus, .message-input:focus-visible {
        outline-width: 0;
    }

    form {
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    .send-button {
        margin-right: 8px;
        cursor: pointer;
    }

    .tag-selector {
        margin-right: 4px;
    }

    /** modal **/
    .modal {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 999;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-wrapper {
        display: flex;
        width: 60%;
        margin: 10% auto auto auto;
    }

    .modal-inner {
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
        gap: 2%;
        width: 100%;
    }

    @media only screen and (max-width: 1280px) {
        .modal-wrapper {
            width: 96%;
        }
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 500;
        background-color: #000;
        opacity: .5;
    }
</style>
