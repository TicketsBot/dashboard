<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  {#if label !== undefined}
    <div class="label-wrapper">
        <label for="picker" class="form-label">
            {label}
        </label>
    </div>
  {/if}
  <input type="datetime-local" id="picker" class="form-input" disabled="{disabled}"
         bind:this={picker} on:input on:change bind:value={value} on:click={() => picker.showPicker()}>
</div>

<script>
    let picker;

    export let value;
    export let label;
    export let placeholder;
    export let disabled = false;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>

<style>
    input {
        width: 100%;
    }

    .label-wrapper {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 5px;
    }

    .tooltip-icon {
        cursor: pointer;
    }
</style>