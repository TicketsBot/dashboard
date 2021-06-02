<script>
    export let footer = true;
    export let fill = true;
    export let footerRight = false;
</script>

<div class="card" class:fill>
  <div class="card-header">
    <h4 class="card-title">
      <slot name="title">
        No Title :(
      </slot>
    </h4>
  </div>
  <div class="card-body">
    <slot name="body">
      No Content :(
    </slot>
  </div>

  {#if footer}
    <div class="card-footer">
      <div class="footer-content" class:footerRight>
        <slot name="footer" />
      </div>
    </div>
  {/if}
</div>

<style>
    .card {
        display: flex;
        flex-direction: column;

        background-color: #272727 !important;

        width: 100%;;
        border-radius: 5px;
        box-shadow: 0 4px 4px rgba(0, 0, 0, 0.25);
    }

    .fill {
        height: 100%;
    }

    .card-title {
        color: white;
        font-size: 22px;
        font-weight: bolder;

        padding: 10px 20px;
        margin: 0;
    }

    .card-header {
        display: flex;
        border-bottom: 1px solid rgba(0, 0, 0, .125);
    }

    .card-body {
        display: flex;
        flex: 1;
        min-height: 75px;

        color: white;
        margin: 10px 20px;
    }

    .card-footer {
        display: flex;
        color: white;
        border-top: 1px solid rgba(0, 0, 0, .125);
        padding: 10px 20px;
    }

    .footer-content {
        display: flex;
        align-items: center;
        height: 100%;
        width: 100%;
    }

    .footerRight {
        flex-direction: row-reverse;
    }
</style>