<div class:col-1={col1} class:col-2={col2} class:col-3={col3} class:col-4={col4}>
  <label for="input" class="form-label">{label}</label>
  <input
          id="input" class="form-input" type="number"
          min={min} max={max} bind:value={value}
          on:input={validateMax} on:change={validateMin}>
</div>

<script>
    export let value;
    export let label;
    export let min;
    export let max;

    export let col1 = false;
    export let col2 = false;
    export let col3 = false;
    export let col4 = false;

    function validateMax() {
        if (value > max) {
            value = max;
        }
    }

    // If we validateMin on input, the user can never backspace to enter a number
    function validateMin() {
        if (value < min) {
            value = min;
        }
    }
</script>
