<Dropdown {col1} {col2} {col3} {col4} bind:value label={label}>
  {#if withNull}
    <option value=null>
      {nullLabel}
    </option>
  {/if}
  {#each channels as channel}
    {#if channel.type === 0}
      <option value="{channel.id}">
        #{channel.name}
      </option>
    {/if}
  {/each}
</Dropdown>

<script>
    import Dropdown from "./form/Dropdown.svelte"

    export let value;
    export let label;
    export let channels = [];
    export let withNull = false;
    export let nullLabel = "Disabled";

    $: value, ensureStringified();

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;

    function ensureStringified() {
      if (value === null) {
        value = "null";
      }
    }
</script>