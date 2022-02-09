<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  {#if label !== undefined}
    <label for="input" class="form-label">{label}</label>
  {/if}
  <select id="input" class="form-input" on:change bind:value={value}>
    <slot />
  </select>
</div>

<style>
  select {
      width: 100%;
  }
</style>

<script>
    export let value;
    export let label;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>
