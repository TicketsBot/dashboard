<div class:inline>
  {#if label !== undefined}
    <label class="form-label" style="margin-bottom: 0 !important;">{label}</label>
  {/if}
  <Toggle hideLabel
          toggledColor="{toggledColour}"
          untoggledColor="{untoggledColour}"
          bind:toggled={value}/>
</div>

<script>
    import Toggle from "svelte-toggle";

    export let value;
    export let label;

    export let toggledColour = "#66bb6a";
    export let untoggledColour = "#ccc";
    export let inline = false;
</script>

<style>
  div {
      display: flex;
      flex-direction: column;
  }

  .inline {
      flex-direction: row !important;
      gap: 4px;
  }
</style>