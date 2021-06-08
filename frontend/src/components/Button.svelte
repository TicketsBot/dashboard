<button on:click isTrigger="1" class:fullWidth disabled="{disabled}">
  {#if icon !== undefined}
    <i class="{icon}"></i>
  {/if}
  <span class="content">
    <slot />
  </span>
</button>

<script>
  export let icon;
  export let fullWidth = false;
  export let disabled = false;
</script>

<style>
  button {
      display: flex;

      justify-content: center;
      align-items: center;

      text-align: center;

      color: white;
      background-color: #3472f7;
      border-color: #3472f7;
      border-width: 2px;
      border-radius: .25rem;
      margin: 0;

      cursor: pointer;
      transition: background-color 150ms ease-in-out, border-color 150ms ease-in-out;
      box-shadow: 0 4px 4px rgb(0 0 0 / 25%);
  }

  button:active {
      background-color: #0062cc;
      border-color: #0062cc;
  }

  button:disabled {
      background-color: #6c757d;
      border-color: #6c757d;
      cursor: default;
  }

  button:hover:enabled {
      background-color: #0062cc;
      border-color: #0062cc;
  }

  .content {
      margin-left: 5px;
      margin-right: 5px;
  }

  .fullWidth {
      width: 100%;
  }
</style>