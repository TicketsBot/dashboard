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
    <EmojiInput label="Button Emoji" col4=true bind:value={data.emote}/>
  </div>
  <div class="row">
    <Dropdown col4=true label="Button Style" bind:value={data.button_style}>
      <option value="1">Blue</option>
      <option value="2">Grey</option>
      <option value="3">Green</option>
      <option value="4">Red</option>
    </Dropdown>
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
    <div class="inner" class:inner-show={advancedSettings}>
      <div class="row">
      <Textarea col1=true bind:value={data.welcome_message} label="Welcome Message"
                placeholder="If blank, your server's default welcome message will be used"
                on:input={handleWelcomeMessageUpdate}/>
      </div>
      <div class="row">
        <div class="col-2">
          <label class="form-label">Mention On Open</label>
          <div class="multiselect-super">
            <Select items={mentionValues} bind:selectedValue={mentionsRaw} on:select={updateMentions} isMulti={true}/>
          </div>
        </div>
        <div class="col-2">
          <label class="form-label">Support Teams</label>
          <div class="multiselect-super">
            <Select items={teamsItems} bind:selectedValue={teamsRaw} on:select={updateTeams} isMulti={true}/>
          </div>
        </div>
      </div>
      <div class="row">
        <Input col2={true} label="Large Image URL" bind:value={data.image_url}/>
        <Input col2={true} label="Small Image URL" bind:value={data.thumbnail_url}/>
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

    import {createEventDispatcher, onMount} from 'svelte';
    import {colourToInt} from "../../js/util";
    import CategoryDropdown from "../CategoryDropdown.svelte";
    import EmojiInput from "../form/EmojiInput.svelte";
    import Select from 'svelte-select';
    import Dropdown from "../form/Dropdown.svelte";

    export let guildId;
    export let seedDefault = true;

    const dispatch = createEventDispatcher();

    let tempColour = '#2ECC71';

    export let data;
    if (seedDefault) {
        data = {
            //title: 'Open a ticket!',
            //content: 'By clicking the button, a ticket will be opened for you.',
            colour: 0x2ECC71,
            emote: '📩',
            welcome_message: null,
            mentions: [],
            default_team: true,
            teams: [],
            button_style: "1",
        };
    }

    export let channels = [];
    export let roles = [];
    export let teams = [];

    let advancedSettings = false;
    let overflowShow = false;

    // Oh my
    // TODO: Clean up
    let mentionValues = [{value: 'user', label: 'Ticket Opener'}];
    let mentionsRaw = [];

    function updateMentions() {
        if (mentionsRaw === undefined) {
            mentionsRaw = [];
        }

        data.mentions = mentionsRaw.map((option) => option.value);
    }

    let teamsItems = [{value: 'default', label: 'Default'}];
    let teamsRaw = [];
    if (seedDefault) {
        teamsRaw = [{value: 'default', label: 'Default'}];
    }

    function updateTeams() {
        if (teamsRaw === undefined) {
            data.teams = [];
        } else {
            data.default_team = teamsRaw.find((option) => option.value === 'default') !== undefined;
            data.teams = teamsRaw
                .filter((option) => option.value !== 'default')
                .map((option) => teams.find((team) => team.id == option.value));
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

    function handleWelcomeMessageUpdate() {
        if (data.welcome_message === "") {
            data.welcome_message = null;
        }
    }

    function updateColour() {
        data.colour = colourToInt(tempColour);
    }

    function updateMentionValues() {
        mentionValues = [{value: 'user', label: 'Ticket Opener'}];
        $: roles.forEach((role) => mentionValues.push({value: role.id, label: role.name}));
    }

    function updateTeamsItems() {
        teamsItems = [{value: 'default', label: 'Default'}];
        $: teams.forEach((team) => teamsItems.push({value: team.id, label: team.name}));
    }

    function applyOverrides() {
        if (data.default_team === true) {
            $: teamsRaw.push({value: 'default', label: 'Default'});
        }

        if (data.teams) {
            $: data.teams.forEach((team) => teamsRaw.push({value: team.id.toString(), label: team.name}));
        }

        if (data.mentions) {
            $: data.mentions.forEach((id) => mentionsRaw.push(mentionValues.find((val) => val.value === id)));
        }
    }

    onMount(() => {
        updateMentionValues();
        updateTeamsItems();

        if (seedDefault) {
            data.channel_id = channels.find((c) => c.type === 0).id;
            data.category_id = channels.find((c) => c.type === 4).id;
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

    .advanced-settings {
        transition: min-height .3s ease-in-out, margin-top .3s ease-in-out, margin-bottom .3s ease-in-out;
        position: relative;
        overflow: hidden;
    }

    .advanced-settings-hide {
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

    .show-overflow {
        overflow: visible;
    }

    .inner {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-items: flex-start;
        position: absolute;
        height: 100%;
        width: 100%;
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