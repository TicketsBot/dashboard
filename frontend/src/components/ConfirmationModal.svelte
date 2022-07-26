<div class="modal" transition:fade>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Embed Builder</span>

      <div slot="body" class="body-wrapper">
        <slot name="body"></slot>
      </div>

      <div slot="footer" style="gap: 12px">
        <Button danger={!isDangerous} on:click={() => dispatch("cancel", {})}>Cancel</Button>
        <Button danger={isDangerous} {icon} on:click={() => dispatch("confirm", {})}>
          <slot name="confirm"></slot>
        </Button>
      </div>
    </Card>
  </div>
</div>

<div class="modal-backdrop" transition:fade>
</div>

<script>
    import Card from "./Card.svelte";
    import {fade} from "svelte/transition";

    import {createEventDispatcher} from "svelte";
    import Button from "./Button.svelte";

    const dispatch = createEventDispatcher();

    export let icon;
    export let isDangerous = false;
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