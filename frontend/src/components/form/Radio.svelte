<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  <label class="form-label">{label}</label>
  <slot />
</div>

<script>
    export let label;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>

<style>
    select {
        width: 100%;
    }

    :global(.radio-row) {
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    :global(.radio-label) {
        display: flex;
        align-items: center;
        text-transform: uppercase;
        font-size: 12px;
        margin: 0 0 7px 5px;
    }

    :global(.radio-input) {
        display: flex;
        align-items: center;
    }
</style>
