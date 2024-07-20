<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  <div class="label-wrapper">
      <slot name="label">
          <label for="input" class="form-label">
              {label}
          </label>
      </slot>
      {#if premiumBadge}
          <div style="margin-bottom: 5px">
              <PremiumBadge />
          </div>
      {/if}
  </div>
  <input id="input" class="form-checkbox" type=checkbox bind:checked={value} on:change {disabled}>
</div>

<script>
    import PremiumBadge from "../PremiumBadge.svelte";

    export let value;
    export let label;
    export let disabled = false;
    export let premiumBadge = false;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>

<style>
    .form-checkbox {
        height: 40px;
        width: 40px;
        margin: 0 !important;
    }

    .label-wrapper {
        display: flex;
        flex-direction: row;
        gap: 4px;
        align-items: center;
    }
</style>
