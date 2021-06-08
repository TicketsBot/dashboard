<Dropdown col1={col1} col2={col2} col3={col3} col4={col4} bind:value={value} label={label}>
  {#each channels as channel}
    {#if channel.type === 4}
      <option value={channel.id}>
        {channel.name}
      </option>
    {/if}
  {/each}
</Dropdown>

<script>
    import Dropdown from "./form/Dropdown.svelte"

    export let value;
    export let label;
    export let channels;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;
</script>