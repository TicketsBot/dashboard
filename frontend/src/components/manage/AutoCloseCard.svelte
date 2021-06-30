<Card footer="{false}">
  <span slot="title">
    Auto Close
  </span>

  <div slot="body" class="body-wrapper">
    <div class="alert danger">
      <span class="alert-text">
        <span>
          This feature is currently disabled. Discord will soon be releasing message threads,
        which will incorporate their own auto-close behaviour. Thank you for your patience.
        </span>
      </span>
    </div>
  </div>
</Card>

<style>
    .body-wrapper {
        width: 100%;
    }

    .alert {
        display: flex;
        justify-content: center;
        width: 100%;
        border-radius: 4px;
        padding: 10px 0;
    }

    .danger {
        background-color: #fc727a;
    }

    .alert-text {
        width: 95%
    }
</style>

<script>
    import Card from "../Card.svelte";
</script>
