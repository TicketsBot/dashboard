<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  {#if label !== undefined}
    <div class="label-wrapper" class:no-margin={tooltipText !== undefined}>
      <label for="input" class="form-label" style="margin-bottom: 0">{label}</label>
      {#if badge !== undefined}
        <Badge>{badge}</Badge>
      {/if}
      {#if tooltipText !== undefined}
        <div>
          <Tooltip tip={tooltipText} top color="#121212">
            {#if tooltipLink !== undefined}
              <a href={tooltipLink} target="_blank">
                <i class="fas fa-circle-info form-label tooltip-icon"></i>
              </a>
            {:else}
              <i class="fas fa-circle-info form-label tooltip-icon"></i>
            {/if}
          </Tooltip>
        </div>
      {/if}
    </div>
  {/if}
  <input id="input" class="form-input" placeholder="{placeholder}" disabled="{disabled}" on:input on:change
         bind:value={value}>
</div>

<script>
    import Tooltip from 'svelte-tooltip';
    import Badge from "../Badge.svelte";

    export let value;
    export let label;
    export let placeholder;
    export let badge = undefined;
    export let disabled = false;

    export let tooltipText = undefined;
    export let tooltipLink = undefined;

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
        margin-bottom: 5px;
    }

    .no-margin {
        margin-bottom: 0 !important;
    }

    .tooltip-icon {
        cursor: pointer;
    }
</style>