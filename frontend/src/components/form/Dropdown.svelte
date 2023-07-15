<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  {#if label !== undefined}
    <div class="label-wrapper">
      <label for="input" class="form-label">{label}</label>
      {#if premiumBadge}
        <div style="margin-bottom: 5px">
          <PremiumBadge />
        </div>
      {/if}
    </div>
  {/if}
  <select id="input" class="form-input" on:change bind:value={value} {disabled}>
    <slot />
  </select>
</div>

<style>
  select {
      width: 100%;
  }

  .label-wrapper {
      display: flex;
      flex-direction: row;
      align-items: center;
      gap: 5px;
  }
</style>

<script>
    import PremiumBadge from "../PremiumBadge.svelte";

    export let value;
    export let label;
    export let disabled = false;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;

    export let premiumBadge = false;
</script>
