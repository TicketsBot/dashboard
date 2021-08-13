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
      <ChannelDropdown col1={true} {channels} label="Panel Channel" bind:value={data.channel_id}/>
    </div>
  </div>
  <div class="row">
    <div class="col-1">
      <PanelDropdown label="Panels" bind:panels bind:selected={data.panels} />
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

    export let data = {};

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
        height: 100%;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }
    }
</style>
