<div class="nav-element">
  {#if link}
    <Navigate to="{link}" styles="link row" on:click>
      <div class="icon">
        <i class="{icon}"></i>
      </div>
      <slot/>
    </Navigate>
  {:else}
    <a class="link row" on:click>
      <div class="icon">
        <i class="{icon}"></i>
      </div>
      <slot/>
    </a>
  {/if}
</div>

<script>
    export let icon;
    export let link;

    import {Navigate} from 'svelte-router-spa'
</script>

<style>
    .nav-element {
        padding: 20px 0 20px 15px;
    }

    :global(.link) {
        display: flex;
        color: inherit;
        text-decoration: none;
        cursor: pointer;
    }

    .icon {
        width: 24px;
    }
</style>
