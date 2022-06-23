<div class="modal" transition:fade>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Grant Permissions To Tickets Team</span>

      <div slot="body" class="body-wrapper">
        Grant permission for
        <Dropdown bind:value={timePeriod}>
          <option value="1">1 hour</option>
          <option value="6">6 hours</option>
          <option value="24">1 day</option>
          <option value="72">3 days</option>
        </Dropdown>
      </div>

      <div slot="footer" class="footer-wrapper">
        <Button danger={true} on:click={dispatchClose}>Cancel</Button>
        <div style="">
          <Button on:click={dispatchConfirm}>Confirm</Button>
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
    import Card from "../Card.svelte";
    import Button from "../Button.svelte";
    import Dropdown from "../form/Dropdown.svelte";

    export let guildId;

    const dispatch = createEventDispatcher();

    let timePeriod = "1";

    function dispatchClose() {
        dispatch('close', {});
    }

    // Dispatch with data
    function dispatchConfirm() {
        dispatch('confirm', {timePeriod: parseInt(timePeriod)});
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
        width: 40%;
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

    .body-wrapper {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .footer-wrapper {
        display: flex;
        flex-direction: row;
        gap: 12px;
    }
</style>