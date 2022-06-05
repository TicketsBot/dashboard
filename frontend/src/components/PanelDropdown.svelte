<label class="form-label">{label}</label>
<div class="multiselect-super">
  <Select placeholder="Select..." items={panels} optionIdentifier="panel_id" getOptionLabel={labelMapper}
          getSelectionLabel={labelMapper} bind:selectedValue={selectedRaw}
          on:select={update} on:clear={handleClear} {isMulti} />
</div>

<script>
    import Select from 'svelte-select';
    import {onMount} from "svelte";

    export let label;
    export let panels;
    export let selected;
    export let isMulti = true;

    let selectedRaw = isMulti ? panels.filter((p) => selected.includes(p.panel_id)) : selected;

    function labelMapper(panel) {
        return panel.title || "";
    }

    function update() {
        if (selectedRaw === undefined) {
            selectedRaw = [];
        }

        if (isMulti) {
            selected = selectedRaw.map((panel) => panel.panel_id);
        } else {
            if (selectedRaw) {
                selected = selectedRaw.panel_id;
            } else {
                selected = undefined;
            }
        }
    }

    function handleClear() {
        if (isMulti) {
            selected = [];
        } else {
            selected = undefined;
        }
    }

    function applyOverrides() {
        if (isMulti) {
            //selected = [];
            selectedRaw = panels.filter((p) => selected.includes(p.panel_id));
        } else {
            if (selectedRaw) {
                selectedRaw = selectedRaw.panel_id;
            }
        }
    }

    onMount(() => {
        applyOverrides();
    });
</script>