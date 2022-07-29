<form on:submit|preventDefault>
  <div class="row">
    <Input col1={true} label="Panel Title" placeholder="Click to open a ticket" bind:value={data.title}/>
  </div>
  <div class="row">
    <Textarea col1={true} label="Panel Content" bind:value={data.content}
              placeholder="Click on the button corresponding to the type of ticket you wish to open. Let users know which button responds to which category. You are able to use emojis here."/>
  </div>
  <div class="row">
    <div class="col-1-3">
      <Colour col1={true} label="Panel Colour" on:change={updateColour} bind:value={tempColour}/>
    </div>
    <div class="col-2-3">
      <ChannelDropdown col1 allowAnnouncementChannel {channels} label="Panel Channel" bind:value={data.channel_id}/>
    </div>
  </div>
  <div class="row">
    <div class="col-3-4" style="padding-right: 10px">
      <PanelDropdown label="Panels" {panels} bind:selected={data.panels} />
    </div>

    <div class="col-1-4">
      <Checkbox label="Use Select Menu" bind:value={data.select_menu} />
    </div>
  </div>

  <div class="row" style="justify-content: center; padding-top: 10px">
    <div class="col-1">
      <Button icon="fas fa-sliders-h" fullWidth=true type="button"
              on:click={toggleAdvancedSettings}>Toggle Advanced Settings
      </Button>
    </div>
  </div>
  <div class="row advanced-settings" class:advanced-settings-show={advancedSettings}
       class:advanced-settings-hide={!advancedSettings} class:show-overflow={overflowShow}>
    <div class="inner" class:inner-show={advancedSettings} class:absolute={advancedSettings && !overflowShow} >
      <div class="row">
        <Input col1={true} label="Large Image URL" bind:value={data.image_url} placeholder="https://example.com/image.png" />
      </div>
      <div class="row">
        <Input col1={true} label="Small Image URL" bind:value={data.thumbnail_url} placeholder="https://example.com/image.png" />
      </div>
    </div>
  </div>
</form>

<script>
    import Input from "../form/Input.svelte";
    import Textarea from "../form/Textarea.svelte";
    import Colour from "../form/Colour.svelte";
    import {colourToInt, intToColour} from "../../js/util";
    import ChannelDropdown from "../ChannelDropdown.svelte";
    import PanelDropdown from "../PanelDropdown.svelte";
    import {onMount} from "svelte";
    import Checkbox from "../form/Checkbox.svelte";
    import Button from "../Button.svelte";

    export let data;

    export let guildId;
    export let channels = [];
    export let panels = [];

    export let seedDefault = true;
    if (seedDefault) {
      data = {
        colour: 0x7289da,
        channels: channels[0].id,
        panels: [],
      }
    }

    let mounted = false;
    let advancedSettings = false;
    let overflowShow = false;

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

    let tempColour = '#7289da';

    function updateColour() {
        data.colour = colourToInt(tempColour);
    }

    function applyOverrides() {
        tempColour = intToColour(data.colour);
    }

    onMount(() => {
      if (!seedDefault) {
        applyOverrides();
      }
    })
</script>

<style>
    form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }
    }

    :global(.col-1-4) {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      width: 25%;
      height: 100%;
    }

    :global(.col-3-4) {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      width: 75%;
      height: 100%;
    }

    .inner {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      align-items: flex-start;
      height: 100%;
      width: 100%;

      margin-top: 10px;
    }

    .absolute {
      position: absolute;
    }

    .advanced-settings-show {
      visibility: visible;
      min-height: 142px;
    }
</style>
