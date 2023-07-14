<Select {items}
        bind:value={selectedValue}
        on:select
        searchable={searchable}
        clearable={isClearable}
        itemId={optionIdentifier}
        debounceWait={loadOptionsInterval}
        {loadOptions}
        {placeholder}
        {placeholderAlwaysShow}
        multiple={isMulti}
        --background="#2e3136"
        --border="#2e3136"
        --border-focused="0"
        --border-hover="0"
        --borderRadius="4px"
        --item-hover-bg="#121212"
        --list-background="#2e3136"
        --item-color="white"
        --multi-item-bg="#272727"
        --multi-item-outline="0"
        --multi-item-border-radius="16px"
        --multi-item-padding="0 10px 0 15px"
        --multi-item-clear-icon-color="#52616F"
        --multi-item-height="32px"
        --clear-icon-color="white"
        --clear-select-color="white"
        --input-color="white"
        --font-size="16px">
  <div slot="item" let:item>{nameMapper(item)}</div>
  <div slot="selection" let:selection>{nameMapper(selection)}</div>
</Select>

<script>
    import Select from "svelte-select";

    export let items;
    export let selectedValue;

    export let isSearchable = true;
    export let isClearable = true;
    export let isMulti = false;
    export let placeholder;
    export let placeholderAlwaysShow = false;
    export let loadOptions;
    export let loadOptionsInterval;

    export let optionIdentifier;
    export let nameMapper = (x) => x;

    let searchable = isSearchable && selectedValue == undefined;
    $: selectedValue, updateSearchable();

    function updateSearchable() {
        searchable = isSearchable && selectedValue == undefined;
    }
</script>
