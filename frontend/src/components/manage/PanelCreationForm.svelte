{#if welcomeMessageBuilder}
    <EmbedBuilder data={data.welcome_message}
                  on:close={closeWelcomeMessageBuilder}
                  on:confirm={handleWelcomeMessageUpdate}/>
{/if}

<form class="settings-form" on:submit|preventDefault>
    <div class="row">
        <div class="col-1-3">
            <Input label="Panel Title" placeholder="Open a ticket!" col1=true bind:value={data.title}/>
        </div>
        <div class="col-2-3">
            <Textarea col1=true label="Panel Content" placeholder="By clicking the button, a ticket will be opened for you."
                bind:value={data.content}/>
        </div>
    </div>
    <div class="row">
        <Colour col4=true label="Panel Colour" on:change={updateColour} bind:value={tempColour}/>
        <ChannelDropdown label="Panel Channel" col4=true {channels} bind:value={data.channel_id}/>
        <CategoryDropdown label="Ticket Category" col4=true {channels} bind:value={data.category_id}/>
        <Dropdown col4=true label="Form" bind:value={data.form_id}>
            <option value=null>None</option>
            {#each forms as form}
                <option value={form.form_id}>{form.title}</option>
            {/each}
        </Dropdown>
    </div>
    <div class="row">
        <Dropdown col4=true label="Button Style" bind:value={data.button_style}>
            <option value="1">Blue</option>
            <option value="2">Grey</option>
            <option value="3">Green</option>
            <option value="4">Red</option>
        </Dropdown>

        <Input col4={true} label="Button Text" placeholder="Open a ticket!" bind:value={data.button_label} />

        <div class="col-2" style="z-index: 1">
            <label for="emoji-pick-wrapper" class="form-label">Button Emoji</label>
            <div id="emoji-pick-wrapper" class="row">
                <div class="col-2">
                    <label class="form-label" style="margin-bottom: 0 !important;">Custom Emoji</label>
                    <Toggle hideLabel
                            toggledColor="#66bb6a"
                            untoggledColor="#ccc"
                            bind:toggled={data.use_custom_emoji}
                            on:toggle={handleEmojiTypeChange} />
                </div>
                {#if data.use_custom_emoji}
                    <div class="multiselect-super">
                        <Select items={emojis}
                                Item={EmojiItem}
                                selectedValue={data.emote}
                                optionIdentifier="id"
                                getSelectionLabel={emojiNameMapper}
                                getOptionLabel={emojiNameMapper}
                                placeholderAlwaysShow={true}
                                on:select={handleCustomEmojiChange} />
                    </div>
                {:else}
                    <EmojiInput col1=true bind:value={data.emote}/>
                {/if}
            </div>
        </div>
    </div>
    <div class="row" style="justify-content: center">
        <div class="col-3">
            <Button icon="fas fa-sliders-h" fullWidth=true type="button"
                    on:click={toggleAdvancedSettings}>Toggle Advanced Settings
            </Button>
        </div>
    </div>
    <div class="row advanced-settings" class:advanced-settings-show={advancedSettings}
         class:advanced-settings-hide={!advancedSettings} class:show-overflow={overflowShow}>
        <div class="inner" class:inner-show={advancedSettings} class:absolute={advancedSettings && !overflowShow} >
            <div class="row">
                <div class="col-2">
                    <label class="form-label">Welcome Message</label>
                    <div class="row" style="justify-content: flex-start; gap: 10px">
                        <Button icon="fas fa-brush" on:click={openWelcomeMessageBuilder}>Open Editor</Button>
                        <Button icon="fas fa-trash-can" danger
                                on:click={() => data.welcome_message = null}>Clear</Button>
                    </div>
                </div>
                <div class="col-2">
                    <label for="naming-scheme-wrapper" class="form-label">Naming Scheme</label>
                    <div class="row" id="naming-scheme-wrapper">
                        <div>
                            <label class="form-label">Use Server Default</label>
                            <Toggle hideLabel
                                    toggledColor="#66bb6a"
                                    untoggledColor="#ccc"
                                    bind:toggled={data.use_server_default_naming_scheme} />
                        </div>
                        <div class="col-fill">
                            {#if !data.use_server_default_naming_scheme}
                                <Input label="Naming Scheme"
                                       bind:value={data.naming_scheme}
                                       placeholder="ticket-%id%"
                                       tooltipText="Click here for the full placeholder list"
                                       tooltipLink="https://docs.ticketsbot.net" />
                            {/if}
                        </div>
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col-2">
                    <label class="form-label">Mention On Open</label>
                    <div class="multiselect-super">
                        <Select items={mentionItems}
                                bind:selectedValue={selectedMentions}
                                on:select={updateMentions}
                                optionIdentifier="id"
                                getSelectionLabel={mentionNameMapper}
                                getOptionLabel={mentionNameMapper}
                                placeholderAlwaysShow={true}
                                isMulti={true} />
                    </div>
                </div>
                <div class="col-2">
                    <label class="form-label">Support Teams</label>
                    <div class="multiselect-super">
                        <Select items={teamsWithDefault}
                                bind:selectedValue={selectedTeams}
                                on:select={updateTeams}
                                isSearchable={false}
                                optionIdentifier="id"
                                getSelectionLabel={nameMapper}
                                getOptionLabel={nameMapper}
                                isMulti={true} />
                    </div>
                </div>
            </div>
            <div class="row">
                <Input col2={true} label="Large Image URL" bind:value={data.image_url} placeholder="https://example.com/image.png" />
                <Input col2={true} label="Small Image URL" bind:value={data.thumbnail_url} placeholder="https://example.com/image.png" />
            </div>
        </div>
    </div>
</form>

<script>
    import Input from "../form/Input.svelte";
    import Textarea from "../form/Textarea.svelte";
    import Colour from "../form/Colour.svelte";
    import Button from "../Button.svelte";
    import ChannelDropdown from "../ChannelDropdown.svelte";
    import EmbedBuilder from "../EmbedBuilder.svelte";

    import {createEventDispatcher, onMount} from 'svelte';
    import {colourToInt, intToColour} from "../../js/util";
    import CategoryDropdown from "../CategoryDropdown.svelte";
    import EmojiInput from "../form/EmojiInput.svelte";
    import EmojiItem from "../EmojiItem.svelte";
    import Select from 'svelte-select';
    import Dropdown from "../form/Dropdown.svelte";
    import Toggle from "svelte-toggle";

    export let guildId;
    export let seedDefault = true;

    const dispatch = createEventDispatcher();

    let tempColour = '#2ECC71';

    export let data = {};

    export let channels = [];
    export let roles = [];
    export let emojis = [];
    export let teams = [];
    export let forms = [];

    let advancedSettings = false;
    let overflowShow = false;

    let teamsWithDefault = [];
    let mentionItems = [];

    let selectedTeams = seedDefault ? [{id: 'default', name: 'Default'}] : [];
    let selectedMentions = [];

    let welcomeMessageBuilder = false;

    function openWelcomeMessageBuilder() {
        welcomeMessageBuilder = true;
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    function closeWelcomeMessageBuilder() {
        welcomeMessageBuilder = false;
    }

    // Replace spaces with dashes in naming scheme as the user types
    $: if (data.naming_scheme !== undefined && data.naming_scheme !== null && data.naming_scheme.includes(' ')) {
        data.naming_scheme = data.naming_scheme.replaceAll(' ', '-');
    }

    function updateMentions() {
        if (selectedMentions === undefined) {
            selectedMentions = [];
        }

        data.mentions = selectedMentions.map((option) => option.id);
    }

    function updateTeams() {
        if (selectedTeams === undefined) {
            selectedTeams = [];

            data.default_team = false;
            data.teams = [];
        } else {
            data.default_team = selectedTeams.find((option) => option.id === 'default') !== undefined;
            data.teams = selectedTeams
                .filter((option) => option.id !== 'default')
                .map((option) => parseInt(option.id));
        }
    }

    const nameMapper = (team) => team.name;
    const emojiNameMapper = (emoji) => `:${emoji.name}:`;

    function mentionNameMapper(role) {
        if (role.id === "user") {
            return role.name;
        } else {
            return `@${role.name}`;
        }
    }

    function toggleAdvancedSettings() {
        advancedSettings = !advancedSettings;
        if (advancedSettings) {
            setTimeout(() => {
                overflowShow = true;
            }, 300);
        } else {
            overflowShow = false;
        }
    }

    function handleWelcomeMessageUpdate(e) {
        data.welcome_message = e.detail;
        closeWelcomeMessageBuilder();
    }

    function handleEmojiTypeChange(e) {
        let isCustomEmoji = e.detail;
        if (isCustomEmoji) {
            data.emote = undefined;
        } else {
            data.emote = '📩';
        }
    }

    function handleCustomEmojiChange(e) {
        let emoji = e.detail;
        data.emote = {
            id: emoji.id,
            name: emoji.name
        };
    }

    function updateColour() {
        data.colour = colourToInt(tempColour);
    }

    function updateMentionValues() {
        mentionItems = [{id: 'user', name: 'Ticket Opener'}, ...roles];
    }

    function updateTeamsItems() {
        teamsWithDefault = [{id: 'default', name: 'Default'}, ...teams];
    }

    function applyOverrides() {
        if (data.default_team === true) {
            $: selectedTeams.push({id: 'default', name: 'Default'});
        }

        if (data.teams) {
            $: data.teams
                .map((id) => teams.find((team) => team.id === id))
                .forEach((team) => selectedTeams.push(team));
        }

        if (data.mentions) {
            $: data.mentions
                .map((id) => mentionItems.find((role) => role.id === id))
                .forEach((mention) => selectedMentions.push(mention));
        }

        data.emote = data.emote;

        tempColour = intToColour(data.colour);
    }

    onMount(() => {
        updateMentionValues();
        updateTeamsItems();

        if (seedDefault) {
            data = {
              //title: 'Open a ticket!',
              //content: 'By clicking the button, a ticket will be opened for you.',
              colour: 0x2ECC71,
              use_custom_emoji: false,
              emote: '📩',
              welcome_message: null,
              mentions: [],
              default_team: true,
              teams: [],
              button_style: "1",
              form_id: "null",
              channel_id: channels.find((c) => c.type === 0)?.id,
              category_id: channels.find((c) => c.type === 4)?.id,
              use_server_default_naming_scheme: true,
            };
        } else {
            applyOverrides();
        }
    })
</script>

<style>
    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        margin-bottom: 10px;
    }

    form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .col-fill {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
    }

    :global(.col-1-3) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 32%;
        height: 100%;
    }

    :global(.col-2-3) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 64%;
        height: 100%;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
            justify-content: center;
        }
        :global(.col-1-3, .col-2-3) {
            width: 100% !important;
        }
    }

    :global(.advanced-settings) {
        transition: min-height .3s ease-in-out, margin-top .3s ease-in-out, margin-bottom .3s ease-in-out;
        position: relative;
        overflow: hidden;
    }

    :global(.advanced-settings-hide) {
        height: 0;
        visibility: hidden;
        margin: 0;
        flex: unset;
        min-height: 0 !important;
    }

    .advanced-settings-show {
        visibility: visible;
        min-height: 297px;
        margin-bottom: 10px;
    }

    :global(.show-overflow) {
        overflow: visible;
    }

    .inner {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: flex-start;
        /*position: absolute;*/
        height: 100%;
        width: 100%;
    }

    .absolute {
        position: absolute;
    }

    #naming-scheme-wrapper {
        gap: 10px;
    }

    :global(.multiselect-super) {
        display: flex;
        width: 100%;
        height: 100%;
        --background: #2e3136;
        --border: #2e3136;
        --borderRadius: 4px;
        --itemHoverBG: #121212;
        --listBackground: #2e3136;
        --itemColor: white;
        --multiItemBG: #272727;
        --multiItemActiveBG: #272727;
        --multiClearFill: #272727;
        --multiClearHoverFill: #272727;
        --inputColor: white;
        --inputFontSize: 16px;
    }

    :global(.multiselect-super > .selectContainer) {
        width: 100%;
    }

    :global(.selectContainer > .multiSelect, .selectContainer > .multiSelect > input) {
        cursor: pointer;
    }
</style>
