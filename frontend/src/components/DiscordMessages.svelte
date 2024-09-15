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

<section class="discord-container">
    <div class="channel-header">
        <span class="channel-name">#ticket-{ticketId}</span>
    </div>
    <div class="message-container" bind:this={container}>
        {#each messages as message}
            <div class="message">
                <img class="avatar" src={getAvatarUrl(message.author.id, message.author.avatar)}
                     on:error={(e) => handleAvatarLoadError(e, message.author.id)} alt="Avatar"/>
                <div>
                    <div>
                        <span class="username">{message.author.global_name || message.author.username}</span>
                        <span class="timestamp">
                            {new Date() - new Date(message.timestamp) < 86400000 ? getRelativeTime(new Date(message.timestamp)) : new Date(message.timestamp).toLocaleString()}
                        </span>
                    </div>
                    <div class="content">
                        {#if message.content?.length > 0}
                            <span class="plaintext">{message.content}</span>
                        {/if}

                        {#if message.embeds?.length > 0}
                            <div class="embed-wrapper">
                                {#each message.embeds.filter(e => 'color' in e) as embed}
                                    <div class="embed">
                                        <div class="colour" style="background-color: #{embed.color.toString(16)}"></div>
                                        <div class="main">
                                            {#if embed.title}
                                                <b>{embed.title}</b>
                                            {/if}
                                            {#if embed.description}
                                                <span>{embed.description}</span>
                                            {/if}

                                            {#if embed.fields && embed.fields.length > 0}
                                                <div class="fields">
                                                    {#each embed.fields as field}
                                                        <div class="field" class:inline={field.inline}>
                                                            <span class="name">{field.name}</span>
                                                            <span class="value">{field.value}</span>
                                                        </div>
                                                    {/each}
                                                </div>
                                            {/if}

                                            {#if embed.image && embed.image.proxy_url}
                                                <img src={embed.image.proxy_url} alt="Embed Image"/>
                                            {/if}
                                        </div>
                                    </div>
                                {/each}
                            </div>
                        {/if}

                        {#if message.attachments?.length > 0}
                            <div class="attachment-wrapper">
                                {#each message.attachments.filter(a => isImage(a.filename)) as attachment}
                                    {@const proxyUrl = attachment.proxy_url.replaceAll("\u0026", "&")}
                                    <img src={proxyUrl} alt="{attachment.filename}"/>
                                {/each}

                                {#each message.attachments.filter(a => !isImage(a.filename)) as attachment}
                                    {@const directUrl = attachment.url.replaceAll("\u0026", "&")}
                                    {@const proxyUrl = attachment.proxy_url.replaceAll("\u0026", "&")}

                                    <div class="other">
                                        <div class="metadata">
                                            <span class="name">{attachment.filename}</span>
                                            <span class="size">{formatFileSize(attachment.size)}</span>
                                        </div>
                                        <a href="{isCdnUrl(directUrl) ? directUrl : proxyUrl}" target="_blank"
                                           download="{attachment.filename}">
                                            <i class="fa-solid fa-download"></i>
                                        </a>
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </div>
                </div>
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
</section>

<script>
    import {createEventDispatcher, onMount} from "svelte";
    import {fade} from "svelte/transition";
    import Button from "./Button.svelte";
    import Card from "./Card.svelte";
    import Dropdown from "./form/Dropdown.svelte";
    import {getAvatarUrl, getDefaultIcon} from "../js/icons";
    import {getRelativeTime} from "../js/util";

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

    let failed = [];

    function handleAvatarLoadError(e, userId) {
        if (!failed.includes(userId)) {
            failed.push(userId);
            e.target.src = getDefaultIcon(userId);
        }
    }

    function isImage(fileName) {
        const imageExtensions = ['png', 'jpg', 'jpeg', 'gif', 'gifv', 'webp'];
        return imageExtensions.includes(fileName.split('.').pop().toLowerCase());
    }

    function formatFileSize(size) {
        if (size < 1024) return `${size} B`;
        if (size < 1024 * 1024) return `${(size / 1024).toFixed(0)} KB`;
        if (size < 1024 * 1024 * 1024) return `${(size / 1024 / 1024).toFixed(0)} MB`;
        else return `${(size / 1024 / 1024 / 1024).toFixed(1)} GB`;
    }

    function isCdnUrl(url) {
        const parsed = new URL(url);
        return parsed.hostname === 'cdn.discordapp.com';
    }

    onMount(() => {
        messages = messages.map(message => {
            // Sort attachments; image first
            message.attachments = message.attachments.sort((a, b) => {
                if (isImage(a.filename) && !isImage(b.filename)) {
                    return -1;
                } else if (!isImage(a.filename) && isImage(b.filename)) {
                    return 1;
                } else {
                    return 0;
                }
            });

            return message;
        })
    });
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
        /*font-family: 'Catamaran', sans-serif !important;*/
        font-family: 'Poppins', sans-serif !important;
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

    .channel-name {
        color: white;
        font-weight: bold;
        padding-left: 20px;
    }

    .message-container {
        display: flex;
        flex-direction: column;
        flex: 1;
        gap: 10px;

        position: relative;
        overflow-y: scroll;
        overflow-wrap: break-word;

        padding-left: 10px;
    }

    .message {
        display: flex;
        flex-direction: row;
        gap: 10px;
    }

    .message:first-child {
        margin-top: 5px;
    }

    .avatar {
        width: 36px;
        height: 36px;
        border-radius: 50%;
    }

    .message > div {
        display: flex;
        flex-direction: column;
        line-height: 16px;
    }

    .username {
        color: white;
        font-weight: bold;
        font-size: 16px;
    }

    .timestamp {
        font-size: 11px;
        opacity: 0.6;
    }

    .plaintext {
        font-size: 14px;
    }

    .embed-wrapper {
        display: flex;
        flex-direction: column;
        gap: 5px;
        margin-top: 5px;
    }

    .embed {
        display: flex;
        flex-direction: row;
        gap: 5px;
        width: 100%;
        min-width: 300px;
        border-radius: 5px;
        background-color: #272727;
    }

    .embed > .colour {
        width: 4px;
        border-radius: 5px 0 0 5px;
    }

    .embed > .main {
        display: flex;
        flex-direction: column;
        padding: 10px 10px 10px 5px;
        width: 100%;
    }

    .fields {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 5px;
        margin-top: 10px;
    }

    .field {
        display: flex;
        flex-direction: column;
    }

    .field.inline {
        flex: 0 0 calc(33.3333% - 5px);
    }

    .field:not(.inline) {
        flex-basis: 100%;
    }

    .field > .name {
        font-size: 14px;
        font-weight: bold;
    }

    .field > .value {
        font-size: 14px;
    }

    .embed > .main > img {
        width: 100%;
        max-width: 300px;
        margin-top: 5px;
        border-radius: 3px;
    }

    .attachment-wrapper {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 5px;
        row-gap: 5px;
        margin-top: 5px;
    }

    .attachment-wrapper > img {
        box-sizing: border-box;
        width: 40%;
        min-width: 300px;
        border-radius: 5px;
    }

    .attachment-wrapper > .other {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 10px;
        padding: 5px 10px;
        border-radius: 5px;
        background-color: #272727;
    }

    .attachment-wrapper > .other > .metadata {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .attachment-wrapper > .other > .metadata > .name {
        font-size: 14px;
        font-weight: bold;
    }

    .attachment-wrapper > .other > .metadata > .size {
        font-size: 12px;
        opacity: 0.6;
    }

    .attachment-wrapper > .other i {
        font-size: 24px;
        color: white;
        opacity: 0.8;
        cursor: pointer;
    }

    .message-container:last-child {
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
