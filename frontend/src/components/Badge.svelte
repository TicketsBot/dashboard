<div class="badge" style="--badge-background-color: {colour}">
  <slot></slot>
</div>

<script>
  export let colour = '#3472f7';
</script>

<style>
    .badge {
        display: flex;
        align-items: center;

        background-color: var(--badge-background-color, #3472f7);
        border-radius: 2px;
        font-size: 14px;
        padding: 0 4px;
    }
</style>