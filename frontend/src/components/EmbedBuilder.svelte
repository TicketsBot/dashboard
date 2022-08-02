<div class="modal" transition:fade>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Embed Builder</span>

      <div slot="body" class="body-wrapper">
        <EmbedForm bind:data />
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
    import Card from "./Card.svelte";
    import Button from "./Button.svelte";
    import Input from "./form/Input.svelte";
    import Colour from "./form/Colour.svelte";
    import Textarea from "./form/Textarea.svelte";
    import DateTimePicker from "./form/DateTimePicker.svelte";
    import Collapsible from "./Collapsible.svelte";
    import Checkbox from "./form/Checkbox.svelte";
    import EmbedForm from "./EmbedForm.svelte";

    export let guildId;

    export let data;

    const dispatch = createEventDispatcher();

    function dispatchClose() {
        dispatch('close', {});
    }

    // Dispatch with data
    function dispatchConfirm() {
        // Map blank strings to null
        const mapper = (obj) => {
            Object.keys(obj).forEach(key => {
                if (typeof obj[key] === 'string' && obj[key] === '') {
                    obj[key] = null;
                } else if (typeof obj[key] === 'object' && obj[key] !== null && obj[key] !== undefined) {
                    mapper(obj[key]);
                }
            });
        }

        mapper(data);

        dispatch('confirm', data);
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