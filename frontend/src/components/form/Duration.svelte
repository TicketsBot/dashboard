<div class="col">
  <div class="row label">
    <slot name="header"></slot>
  </div>

  <div class="row fields">
    <div class="parent">
      <input class="form-input" type="number" min=0 {disabled} bind:value={days}/>
      <div class="period" class:disabled>D</div>
    </div>

    <div class="parent">
      <input class="form-input" type="number" min=0 {disabled} bind:value={hours}/>
      <div class="period" class:disabled>H</div>
    </div>

    <div class="parent">
      <input class="form-input" type="number" min=0 {disabled} bind:value={minutes}/>
      <div class="period" class:disabled>M</div>
    </div>
  </div>
</div>

<script>
    export let disabled = false; // note: bind:disabled isn't valid

    export let days = 0;
    export let hours = 0;
    export let minutes = 0;
</script>

<style>
    .col {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
    }

    .fields > .parent:not(:first-child) {
        margin-left: 10px;
    }

    input {
        border-top-right-radius: 0 !important;
        border-bottom-right-radius: 0 !important;
        width: 100%;
        -moz-appearance: textfield;
    }

    input:disabled {
        opacity: 0.6;
    }

    .period.disabled {
        opacity: 0.6;
    }

    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

    label {
        display: flex;
        align-items: center;
        margin: 0;
    }

    .label {
        margin-bottom: 4px;
    }

    .header-wrapper {
        display: flex;
        flex-direction: row;
        gap: 4px;
    }

    .parent {
        display: flex;
        flex-direction: row;
    }

    .period {
        display: flex;
        align-items: center;
        border-color: #2e3136 !important;
        background-color: #2e3136 !important;
        color: white !important;
        outline: none;
        border-top-right-radius: 4px;
        border-bottom-right-radius: 4px;
        padding: 0 10px;
        margin: 0 0 0.5em 0;
        height: 40px;
    }
</style>