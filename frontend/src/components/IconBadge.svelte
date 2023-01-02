<Badge>
  <div class="inner">
    {#if icon !== undefined}
      <i class="{icon}"></i>
    {/if}
    <span class="text">
      <slot />
    </span>
  </div>
</Badge>

<script>
    import Badge from "./Badge.svelte";

    export let text;
    export let icon;
</script>

<style>
    .inner {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 4px;
        color: white;
    }
</style>