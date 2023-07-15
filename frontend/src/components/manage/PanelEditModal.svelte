<div class="modal" transition:fade bind:this={modal}>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Edit Panel</span>

      <div slot="body" class="body-wrapper">
        <PanelCreationForm {guildId} {channels} {roles} {emojis} {teams} {forms} bind:data={panel} seedDefault={false} />
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

<svelte:window on:keydown={handleKeydown}/>

<script>
    import {createEventDispatcher} from 'svelte';
    import {fade} from 'svelte/transition'
    import PanelCreationForm from "./PanelCreationForm.svelte";
    import Card from "../Card.svelte";
    import Button from "../Button.svelte";

    export let modal;

    export let guildId;
    export let panel = {};
    export let channels = [];
    export let forms = [];
    export let roles = [];
    export let emojis = [];
    export let teams = []
    export let isPremium = false;

    const dispatch = createEventDispatcher();

    function dispatchClose() {
        dispatch('close', {});
    }

    // Dispatch with data
    function dispatchConfirm() {
        let form_id = (panel.form_id === null  || panel.form_id === "null") ? null : parseInt(panel.form_id);
        let mapped = {...panel, form_id: form_id};
        dispatch('confirm', mapped);
    }

    function handleKeydown(e) {
      if (e.key === "Escape") {
        dispatchClose();
      }
    }
</script>

<style>
    .modal {
        position: absolute;
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
        margin: 2% auto auto auto;
        padding-bottom: 5%;
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