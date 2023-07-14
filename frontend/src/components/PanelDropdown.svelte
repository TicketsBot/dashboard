<label class="form-label">{label}</label>

<WrappedSelect placeholder="Select panel..." items={panels} optionIdentifier="panel_id" nameMapper={labelMapper}
               bind:selectedValue={selectedRaw} on:select={update} on:clear={handleClear} {isMulti} {isSearchable} />

<script>
    import {onMount} from "svelte";
    import WrappedSelect from "./WrappedSelect.svelte";

    export let label;
    export let panels;
    export let selected;
    export let isMulti = true;
    export let isSearchable = false;

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