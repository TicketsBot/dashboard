<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  {#if label !== undefined}
    <label class="form-label">{label}</label>
  {/if}
  <div class="wrapper">
    <input id="input" class="form-input" placeholder="{placeholder}" disabled="{disabled}" bind:value={value}>
    {#if !disabled}
      <EmojiSelector on:emoji={onUpdate}/>
    {/if}
  </div>
</div>

<script>
    export let value;
    export let label;
    export let placeholder;
    export let disabled = false;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;

    import EmojiSelector from 'svelte-emoji-selector'

    function onUpdate(e) {
        value = e.detail;
    }
</script>

<style>
    input {
        width: 100%;
        border-top-right-radius: 0 !important;
        border-bottom-right-radius: 0 !important;
    }

    .wrapper {
        display: flex;
        flex-direction: row;
        width: 100%;
    }

    :global(.svelte-emoji-picker__trigger) {
        border-bottom-left-radius: 0;
        border-top-left-radius: 0;
        background-color: #272727;
        border-color: #2e3136 !important;
        border-left: none;
        color: white;
        z-index: 2;
    }

    :global(.svelte-emoji-picker__trigger:active) {
        background-color: #2e3136 !important;
    }
</style>