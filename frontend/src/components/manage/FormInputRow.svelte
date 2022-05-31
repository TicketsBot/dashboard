<form on:submit|preventDefault={forwardCreate} class="input-form">
  <div class="row">
    <div class="sub-row" style="flex: 1">
      <Input col4={true} label="Label" bind:value={data.label} placeholder="Name of the field" />
    </div>
    <div class="sub-row buttons-row">
      {#if windowWidth > 950}
        {#if withSaveButton}
          <form on:submit|preventDefault={forwardSave} class="button-form">
            <Button icon="fas fa-save">Save</Button>
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
    <Textarea col3_4={true} label="Placeholder" bind:value={data.placeholder} minHeight="120px"
           placeholder="Placeholder text for the field, just like this text" />
    <div class="col-4">
      <div class="row">
        <Dropdown col1={true} label="Style" bind:value={data.style}>
          <option value=1 selected>Short</option>
          <option value=2>Paragraph</option>
        </Dropdown>
      </div>
      <div class="row">
        <Checkbox label="Optional" bind:value={data.optional}/>
      </div>
    </div>
  </div>

  {#if windowWidth <= 950}
    <div class="row">
      {#if withSaveButton}
        <form on:submit|preventDefault={forwardSave} class="button-form">
          <Button icon="fas fa-save">Save</Button>
        </form>
      {/if}
      {#if withDeleteButton}
        <form on:submit|preventDefault={forwardDelete} class="button-form">
          <Button icon="fas fa-trash" danger={true}>Delete</Button>
        </form>
      {/if}
    </div>
  {/if}

  {#if withCreateButton}
    <div class="row" style="justify-content: center; margin-top: 10px">
      <Button type="submit" icon="fas fa-plus">Add Input</Button>
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

  export let withCreateButton = false;
  export let withSaveButton = false;
  export let withDeleteButton = false;

  export let data = {};

  $: windowWidth = 0;

  function forwardCreate() {
    dispatch('create', data);
  }

  function forwardSave() {
    dispatch('save', data);
  }

  function forwardDelete() {
    dispatch('delete', {});
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

    .buttons-row > form:first-of-type {
        margin-right: 10px;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }

        .settings-row {
            flex-direction: column-reverse !important;
        }
    }
</style>