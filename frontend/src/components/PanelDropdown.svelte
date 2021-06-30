<label class="form-label">{label}</label>
<div class="multiselect-super">
  <Select placeholder="Select..." items={panels} optionIdentifier="panel_id" getOptionLabel={labelMapper}
          getSelectionLabel={labelMapper} bind:selectedValue={panelsRaw}
          on:select={update} isMulti={true}/>
</div>

<script>
    import Select from 'svelte-select';

    export let label;
    export let panels = [];
    export let selected = [];

    let panelsRaw = [];

    function labelMapper(panel) {
        return panel.title;
    }

    function update() {
        if (panelsRaw === undefined) {
            panelsRaw = [];
        }

        selected = panelsRaw.map((panel) => panel.panel_id);
    }
</script>