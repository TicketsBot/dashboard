<form class="settings-form" on:submit|preventDefault>
    <Collapsible defaultOpen>
        <span slot="header">Ticket Properties</span>
        <div slot="content" class="col-1">
            <div class="row">
                <div class="col-2">
                    <label class="form-label">Mention On Open</label>
                    <div class="col-1">
                        <WrappedSelect items={mentionItems}
                                       bind:selectedValue={selectedMentions}
                                       on:select={updateMentions}
                                       optionIdentifier="id"
                                       nameMapper={mentionNameMapper}
                                       placeholder="Select roles..."
                                       isMulti={true} />
                    </div>
                </div>
                <div class="col-2">
                    <label class="form-label">Support Teams</label>
                    <WrappedSelect items={teamsWithDefault}
                            bind:selectedValue={selectedTeams}
                            on:select={updateTeams}
                            optionIdentifier="id"
                            nameMapper={nameMapper}
                            placeholder="Select teams..."
                            isMulti={true}>
                        <div slot="item" let:item>{item.name}</div>
                        <div slot="selection" let:selection>{selection.name}</div>
                    </WrappedSelect>
                </div>
            </div>
            <div class="incomplete-row">
                <CategoryDropdown label="Ticket Category" col3 {channels} bind:value={data.category_id}/>

                <Dropdown col4 label="Form" bind:value={data.form_id}>
                    <option value=null>None</option>
                    {#each forms as form}
                        <option value={form.form_id}>{form.title}</option>
                    {/each}
                </Dropdown>

                <div>
                    <label for="naming-scheme-wrapper" class="form-label">Naming Scheme</label>
                    <div class="row" id="naming-scheme-wrapper">
                        <div>
                            <label class="form-label">Use Server Default</label>
                            <Toggle hideLabel
                                    toggledColor="#66bb6a"
                                    untoggledColor="#ccc"
                                    bind:toggled={data.use_server_default_naming_scheme} />
                        </div>
                    </div>
                </div>

                {#if !data.use_server_default_naming_scheme}
                    <Input col4
                           label="Naming Scheme"
                           bind:value={data.naming_scheme}
                           placeholder="ticket-%id%"
                           tooltipText="Click here for the full placeholder list"
                           tooltipLink="https://docs.ticketsbot.net" />
                {/if}
            </div>
            <div class="incomplete-row">
                <Dropdown col3 label="Exit Survey Form" premiumBadge={true} bind:value={data.exit_survey_form_id} disabled={!isPremium}>
                    <option value=null>None</option>
                    {#each forms as form}
                        <option value={form.form_id}>{form.title}</option>
                    {/each}
                </Dropdown>
            </div>
        </div>
    </Collapsible>

    <Collapsible defaultOpen>
        <span slot="header">Panel Message</span>
        <div slot="content" class="col-1">
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
                <ChannelDropdown label="Panel Channel" allowAnnouncementChannel col4 {channels} bind:value={data.channel_id}/>
                <div class="col-2">
                    <div class="row" style="justify-content: flex-start; gap: 10px">
                        <div style="white-space: nowrap">
                            <Checkbox label="Disable Panel" bind:value={data.disabled}></Checkbox>
                        </div>
                        {#if data.disabled}
                            <b style="display: flex; align-self: center">You will be unable to open any tickets with this panel</b>
                        {/if}
                    </div>
                </div>
            </div>

            <div class="row">
                <Dropdown col4=true label="Button Colour" bind:value={data.button_style}>
                    <option value="1">Blue</option>
                    <option value="2">Grey</option>
                    <option value="3">Green</option>
                    <option value="4">Red</option>
                </Dropdown>

                <Input col4={true} label="Button Text" placeholder="Open a ticket!" bind:value={data.button_label} />

                <div class="col-2" style="z-index: 1">
                    <label for="emoji-pick-wrapper" class="form-label">Button Emoji</label>
                    <div id="emoji-pick-wrapper" class="row" style="gap: 2%">
                        <div class="col">
                            <label class="form-label" style="margin-bottom: 0 !important; white-space: nowrap;">Custom Emoji</label>
                            <Toggle hideLabel
                                    toggledColor="#66bb6a"
                                    untoggledColor="#ccc"
                                    bind:toggled={data.use_custom_emoji}
                                    on:toggle={handleEmojiTypeChange} />
                        </div>
                        {#if data.use_custom_emoji}
                            <div class="col-fill">
                                <!--Item=EmojiItem-->
                                <WrappedSelect items={emojis}
                                        selectedValue={data.emote}
                                        optionIdentifier="id"
                                        nameMapper={emojiNameMapper}
                                        placeholderAlwaysShow={true}
                                        on:select={handleCustomEmojiChange} />
                            </div>
                        {:else}
                            <EmojiInput col1=true bind:value={data.emote}/>
                        {/if}
                    </div>
                </div>
            </div>


            <div class="row">
                <Input col2={true} label="Large Image URL" badge="Optional" bind:value={data.image_url} placeholder="https://example.com/image.png" />
                <Input col2={true} label="Small Image URL" badge="Optional" bind:value={data.thumbnail_url} placeholder="https://example.com/image.png" />
            </div>
        </div>
    </Collapsible>

    <Collapsible>
        <span slot="header">Welcome Message</span>
        <div slot="content" class="col-1">
            <div class="row">
                <EmbedForm bind:data={data.welcome_message} />
            </div>
        </div>
    </Collapsible>
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
    import Checkbox from "../form/Checkbox.svelte";
    import Collapsible from "../Collapsible.svelte";
    import EmbedForm from "../EmbedForm.svelte";
    import WrappedSelect from "../WrappedSelect.svelte";

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
    export let isPremium = false;

    let advancedSettings = false;
    let overflowShow = false;

    let teamsWithDefault = [];
    let mentionItems = [];

    let selectedTeams = seedDefault ? [{id: 'default', name: 'Default'}] : [];
    let selectedMentions = [];

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
              mentions: [],
              default_team: true,
              teams: [],
              button_style: "1",
              form_id: "null",
              channel_id: channels.find((c) => c.type === 0 || c.type === 5)?.id,
              category_id: channels.find((c) => c.type === 4)?.id,
              use_server_default_naming_scheme: true,
              welcome_message: {
                  fields: [],
                  colour: '#2ECC71',
                  author: {},
                  footer: {},
                  description: 'Thank you for contacting support.\nPlease describe your issue and wait for a response.'
              },
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

    .incomplete-row {
        display: flex;
        flex-direction: row;
        gap: 10px;
        width: 100%;
        margin-bottom: 10px;
    }

    form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .col {
        display: flex;
        flex-direction: column;
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

    :global(.show-overflow) {
        overflow: visible;
    }

    #naming-scheme-wrapper {
        gap: 10px;
    }
</style>
