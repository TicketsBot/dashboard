<button on:click isTrigger="1" class:fullWidth class:danger class:iconOnly {disabled} {type}>
  {#if icon !== undefined}
    <i class="{icon}"></i>
  {/if}
  {#if !iconOnly}
    <span class="content">
      <slot/>
    </span>
  {/if}
</button>

<script>
    export let icon = undefined;
    export let fullWidth = false;
    export let disabled = false;
    export let type = "submit";
    export let danger = false;
    export let iconOnly = false;
</script>

<style>
    button {
        display: flex;

        justify-content: center;
        align-items: center;

        text-align: center;

        color: white;
        background: var(--primary-gradient);
        border: none;
        /*border-color: var(--primary);*/
        /*border-width: 2px;*/
        border-radius: .25rem;
        margin: 0;

        cursor: pointer;
        transition: background-color 150ms ease-in-out, border-color 150ms ease-in-out;
        box-shadow: 0 4px 4px rgb(0 0 0 / 25%);
    }

    /*button:active, button:hover:enabled {*/
    /*    background-color: #0062cc;*/
    /*    border-color: #0062cc;*/
    /*}*/

    button:disabled {
        background: #6c757d;
        border: #6c757d;
        cursor: default;
    }

    .content {
        margin-left: 5px;
        margin-right: 5px;
    }

    .fullWidth {
        width: 100%;
    }

    .danger {
        background: #dc3545 !important;
        border-color: #dc3545 !important;
    }

    .danger:hover:enabled, .danger:active {
        background: #c32232 !important;
        border-color: #c32232 !important;
    }

    .iconOnly {
        width: 40px;
        height: 40px;
    }
</style>