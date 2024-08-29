<div class="badge" style="--badge-background-color: {colour}">
  <slot></slot>
</div>

<script>
  export let colour = '';
</script>

<style>
    .badge {
        display: flex;
        align-items: center;

        background: var(--badge-background-color, var(--primary));
        border-radius: 2px;
        font-size: 14px;
        padding: 0 4px;
    }
</style>