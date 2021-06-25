<div class="modal" transition:fade>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Edit Panel</span>

      <div slot="body" class="body-wrapper">
        <PanelCreationForm {guildId} {channels} {roles} {teams} bind:data={panel} seedDefault={false} />
      </div>

      <div slot="footer">
        <Button danger={true} on:click={dispatchClose}>Cancel</Button>
        <div style="margin-left: 12px">
          <Button icon="fas fa-paper-plane" on:click={dispatchConfirm}>Submit</Button>
        </div>
      </div>
    </Card>
  </div>
</div>

<div class="modal-backdrop" transition:fade>
</div>

<script>
    import {createEventDispatcher} from 'svelte';
    import {fade} from 'svelte/transition'
    import PanelCreationForm from "./PanelCreationForm.svelte";
    import Card from "../Card.svelte";
    import Button from "../Button.svelte";

    export let guildId;
    export let panel = {};
    export let channels = [];
    export let roles = [];
    export let teams = [];

    const dispatch = createEventDispatcher();

    function dispatchClose() {
        dispatch('close', {});
    }

    // Dispatch with data
    function dispatchConfirm() {
        dispatch('confirm', panel);
    }
</script>

<style>
    .modal {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 501;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-wrapper {
        display: flex;
        width: 75%;
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