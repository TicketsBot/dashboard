<form class="settings-form" on:submit|preventDefault>
  <div class="row">
    <div class="col-1-3">
      <Input label="Panel Title" placeholder="Open a ticket!" col1=true bind:value={data.prefix}/>
    </div>
    <div class="col-2-3">
      <Textarea col1=true label="Panel Content" placeholder="By clicking the button, a ticket will be opened for you."
                bind:value={data.content}/>
    </div>
  </div>
  <div class="row">
    <Colour col4=true label="Panel Colour" on:change={updateColour} bind:value={tempColour}/>
    <ChannelDropdown label="Panel Channel" col4=true channels={channels} bind:value={data.channel_id}/>
    <CategoryDropdown label="Ticket Category" col4=true channels={channels} bind:value={data.category_id}/>
    <EmojiInput label="Button Emoji" col4=true bind:value={data.emote}/>
  </div>
  <div class="row" style="justify-content: center">
    <div class="col-3">
      <Button icon="fas fa-sliders-h" fullWidth=true type="button"
              on:click={() => advancedSettings = !advancedSettings}>Toggle Advanced Settings
      </Button>
    </div>
  </div>
  <div class="row advanced-settings" class:advanced-settings-show={advancedSettings}
       class:advanced-settings-hide={!advancedSettings}>
    <div class="inner">
      <div class="row">
      <Textarea col1=true bind:value={data.welcome_message} label="Welcome Message"
                placeholder="If blank, your server's default welcome message will be used"
                on:input={handleWelcomeMessageUpdate}/>
      </div>
      <div class="row">
        <div class="col-2">
          <MultiSelect bind:value={data.mentions}>
            <option value="user">Ticket Opener</option>
            {#each roles as role}
              <option value="{role.id}">{role.name}</option>
            {/each}
          </MultiSelect>
        </div>
      </div>
    </div>
  </div>
  <div class="row" style="justify-content: center">
    <div class="col-3">
      <Button icon="fas fa-paper-plane" fullWidth=true>Submit</Button>
    </div>
  </div>
</form>

<script>
    import Input from "../form/Input.svelte";
    import Textarea from "../form/Textarea.svelte";
    import Colour from "../form/Colour.svelte";
    import Button from "../Button.svelte";
    import ChannelDropdown from "../ChannelDropdown.svelte";

    import {colourToInt, notifyError, withLoadingScreen} from "../../js/util";
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from "../../includes/Auth.svelte";
    import CategoryDropdown from "../CategoryDropdown.svelte";
    import EmojiInput from "../form/EmojiInput.svelte";
    import MultiSelect from "../form/MultiSelect.svelte";

    export let guildId;

    let tempColour = '#2ECC71';
    export let data = {
        title: 'Open a ticket!',
        content: 'By clicking the button, a ticket will be opened for you.',
        colour: 0x2ECC71,
        emote: '📩',
        welcome_message: null
    };

    let channels = [];
    let roles = [];
    let advancedSettings = false;

    function handleWelcomeMessageUpdate() {
        if (data.welcome_message === "") {
            data.welcome_message = null;
        }
    }

    function updateColour() {
        data.colour = colourToInt(tempColour);
    }

    async function loadChannels() {
        const res = await axios.get(`${API_URL}/api/${guildId}/channels`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        channels = res.data;
    }

    async function loadRoles() {
        const res = await axios.get(`${API_URL}/api/${guildId}/roles`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        roles = res.data.roles;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadChannels();
        await loadRoles();
    })
</script>

<style>
    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-bottom: 10px;
    }

    form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .col-1-3 {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 32%;
        height: 100%;
    }

    .col-2-3 {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 64%;
        height: 100%;
    }

    .advanced-settings {
        transition: min-height .3s ease-in-out, margin-top .3s ease-in-out, margin-bottom .3s ease-in-out;
        position: relative;
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
        min-height: 500px;
        margin-bottom: 10px;
        overflow: hidden;
    }

    .inner {
        display: flex;
        flex-direction: column;
        position: absolute;
        height: 100%;
        width: 100%;
    }
</style>