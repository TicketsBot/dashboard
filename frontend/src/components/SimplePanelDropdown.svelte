<Dropdown {col1} {col2} {col3} {col4} bind:value label={label}>
  {#if allowNone}
    <option value=null>
      None
    </option>
  {/if}

  {#each panels as panel}
    <option value={panel.panel_id}>
      {panel.title}
    </option>
  {/each}
</Dropdown>

<script>
    import Dropdown from "./form/Dropdown.svelte"

    export let value;
    export let label;
    export let allowNone = false;
    export let panels = [];

    $: value, ensureStringified();

    function ensureStringified() {
      if (value === null) {
        value = "null";
      }
    }

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>