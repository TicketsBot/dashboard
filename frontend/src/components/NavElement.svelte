<div class="nav-element">
  {#if link}
    <a href="{link}" class="link row" on:click>
      <div class="icon">
        <i class="{icon}"></i>
      </div>
      <slot/>
    </a>
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
</script>

<style>
    .nav-element {
        padding: 20px 0 20px 15px;
    }

    .row {
        display: flex;
        flex-direction: row;
    }

    :global(.link) {
        display: flex;
        color: inherit;
        text-decoration: none;
        cursor: pointer;
        text-align: center;
    }

    :global(.link-blue) {
        color: #3472f7;
        text-decoration: none;
        cursor: pointer;
        text-align: center;
    }

    .icon {
        width: 24px;
    }
</style>
