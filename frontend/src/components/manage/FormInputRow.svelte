<form on:submit|preventDefault={forwardCreate} class="input-form">
  <div class="row">
    <div class="sub-row" style="flex: 1">
      <Input col3={true} label="Label" bind:value={data.label} placeholder="Name of the field" />
    </div>
    <div class="sub-row buttons-row">
      {#if windowWidth > 950}
        {#if withDirectionButtons}
          <form on:submit|preventDefault={() => forwardMove("down")} class="button-form">
            <Button disabled={index >= formLength - 1}>
              <i class="fas fa-chevron-down"></i>
            </Button>
          </form>
          <form on:submit|preventDefault={() => forwardMove("up")} class="button-form">
            <Button disabled={index === 0}>
              <i class="fas fa-chevron-up"></i>
            </Button>
          </form>
        {/if}
        {#if withDeleteButton}
          <form on:submit|preventDefault={forwardDelete} class="button-form">
            <Button icon="fas fa-trash" danger={true}>Delete</Button>
          </form>
        {/if}
      {/if}
    </div>
  </div>
  <div class="row settings-row">
    <Textarea col2={true} label="Placeholder" bind:value={data.placeholder} minHeight="120px"
           placeholder="Placeholder text for the field, just like this text" />
    <div class="col-2 properties-group">
      <div class="row">
        <Dropdown col2={true} label="Style" bind:value={data.style}>
          <option value=1 selected>Short</option>
          <option value=2>Multi-line</option>
        </Dropdown>
      </div>
      <div class="row" style="gap: 10px">
        <Checkbox label="Required" bind:value={data.required}/>
        <DoubleRangeSlider label="Answer Length Range" bind:start={data.min_length} bind:end={data.max_length} min={0} max={1024} />
      </div>
    </div>
  </div>

  {#if windowWidth <= 950}
    <div class="col-1">
      {#if withDirectionButtons}
        <div class="row">
          <div class="col-2-force">
            <form on:submit|preventDefault={() => forwardMove("down")} class="button-form">
              <Button fullWidth={true} disabled={index >= formLength - 1}>
                <i class="fas fa-chevron-down"></i>
              </Button>
            </form>
          </div>
          <div class="col-2-force">
            <form on:submit|preventDefault={() => forwardMove("up")} class="button-form">
              <Button fullWidth={true} disabled={index === 0}>
                <i class="fas fa-chevron-up"></i>
              </Button>
            </form>
          </div>
        </div>
      {/if}
      <div class="row">
        <div class="col-2-force">
          {#if withDeleteButton}
            <form on:submit|preventDefault={forwardDelete} class="button-form">
              <Button icon="fas fa-trash" danger={true}>Delete</Button>
            </form>
          {/if}
        </div>
      </div>
    </div>
  {/if}

  {#if withCreateButton && false}
    <div class="row" style="justify-content: center; margin-top: 10px">
      <Button type="submit" icon="fas fa-plus" {disabled}>Add Input</Button>
    </div>
  {/if}
</form>

<svelte:window bind:innerWidth={windowWidth} />

<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  import Input from "../form/Input.svelte";
  import Dropdown from "../form/Dropdown.svelte";
  import Button from "../Button.svelte";
  import Textarea from "../form/Textarea.svelte";
  import Checkbox from "../form/Checkbox.svelte";
  import DoubleRangeSlider from "../form/DoubleRangeSlider.svelte";

  export let withCreateButton = false;
  export let withDeleteButton = false;
  export let withDirectionButtons = false;
  export let disabled = false;

  export let index;
  export let formLength;

  export let data = {};

  $: windowWidth = 0;

  function forwardCreate() {
    dispatch('create', data);
  }

  function forwardDelete() {
    dispatch('delete', {});
  }

  function forwardMove(direction) {
    dispatch('move', {direction: direction});
  }
</script>

<style>
    .input-form {
        display: flex;
        flex-direction: column;
        width: 100%;
        border-top: 1px solid rgba(0, 0, 0, .25);
        padding-top: 10px;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
    }

    .sub-row {
        display: flex;
        flex-direction: row;
    }

    .button-form {
        display: flex;
        flex-direction: column;
        justify-content: flex-end;
        padding-bottom: 0.5em;
    }

    .buttons-row > :not(:last-child) {
        margin-right: 10px;
    }

    @media only screen and (max-width: 950px) {
        .settings-row {
            flex-direction: column-reverse !important;
        }

        .button-form {
            width: 100%;
        }
    }

    @media only screen and (max-width: 576px) {
        .properties-group > div:nth-child(2) {
            flex-direction: column;
        }
    }
</style>