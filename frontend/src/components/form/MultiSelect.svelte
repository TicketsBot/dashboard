<!--
  Based upon https://svelte.dev/repl/c7094fb1004b440482d2a88f4d1d7ef5?version=3.14.0
  Heavily ammended
-->

<script>
    import {fly} from 'svelte/transition';

    export let values = {};
    export let selected = [];
    let filtered = [];

    let input,
        inputValue = '',
        options = [],
        activeOption,
        showOptions = false,
        first = true,
        slot
    const iconClearPath = 'M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z';

    function updateFiltered() {
        filtered = Object.entries(values).filter(([_, name]) => name.includes(inputValue));
    }

    updateFiltered();

    /*afterUpdate(() => {
        let newOptions = [];
        slot.querySelectorAll('option').forEach(o => {
            o.selected && !value.includes(o.value) && (value = [...value, o.value]);
            newOptions = [...newOptions, {value: o.value, name: o.textContent}]
        });
        value && (selected = newOptions.reduce((obj, op) => value.includes(op.value) ? {
            ...obj,
            [op.value]: op
        } : obj, {}));
        first = false;
        options = newOptions;
    });

    $: if (!first) value = Object.values(selected).map(o => o.value);
    $: filtered = options.filter(o => inputValue ? o.name.toLowerCase().includes(inputValue.toLowerCase()) : o);
    $: if (activeOption && !filtered.includes(activeOption) || !activeOption && inputValue) activeOption = filtered[0];*/


    function add(value) {
        selected = [...selected, value];
    }

    function remove(value) {
        selected = selected.filter((e) => e !== value);
    }

    function optionsVisibility(show) {
        if (typeof show === 'boolean') {
            showOptions = show;
            show && input.focus();
        } else {
            showOptions = !showOptions;
        }
        if (!showOptions) {
            activeOption = undefined;
        }
    }

    function handleKeyup(e) {
        updateFiltered();
        /*if (e.keyCode === 13) {
            Object.keys(selected).includes(activeOption.value) ? remove(activeOption.value) : add(activeOption);
            inputValue = '';
        }
        if ([38, 40].includes(e.keyCode)) { // up and down arrows
            const increment = e.keyCode === 38 ? -1 : 1;
            const calcIndex = filtered.indexOf(activeOption) + increment;
            activeOption = calcIndex < 0 ? filtered[filtered.length - 1]
                : calcIndex === filtered.length ? filtered[0]
                    : filtered[calcIndex];
        }*/
    }

    function handleBlur(e) {
        optionsVisibility(false);
    }

    function handleTokenClick(e) {
        if (e.target.closest('.token-remove')) {
            e.stopPropagation();
            remove(e.target.closest('.token').dataset.id);
        } else {
            updateFiltered();
            optionsVisibility(true);
        }
    }

    function toggle(id) {
        if (isSelected(id)) {
            selected = selected.filter((e) => e !== id);
        } else {
            $: selected.push(id);
        }

        console.log(selected)
    }

    function isSelected(value) {
        return selected.find((option) => option === value);
    }
</script>

<style>
    .multiselect {
        background-color: #2e3136;
        border-color: #2e3136;
        border-radius: 4px;
        position: relative;
        width: 100%
    }

    .tokens {
        align-items: center;
        display: flex;
        flex-wrap: wrap;
        position: relative;
        padding: 5px 10px;
    }

    .tokens::after {
        background: none repeat scroll 0 0 transparent;
        bottom: -1px;
        content: "";
        display: block;
        height: 2px;
        left: 50%;
        position: absolute;
        background: hsl(45, 100%, 51%);
        transition: width 0.3s ease 0s, left 0.3s ease 0s;
        width: 0;
    }

    .tokens.showOptions::after {
        width: 100%;
        left: 0;
    }

    .token {
        align-items: center;
        background-color: #272727;
        border-radius: 1.25rem;
        display: flex;
        margin: .25rem .5rem .25rem 0;
        max-height: 1.3rem;
        padding: .25rem .5rem .25rem .5rem;
        transition: background-color .3s;
        white-space: nowrap;
    }

    .token-remove, .remove-all {
        align-items: center;
        background-color: #3472f7;
        transition: background-color .3s;
        border-radius: 50%;
        display: flex;
        justify-content: center;
        height: 1.25rem;
        margin-left: .25rem;
        min-width: 1.25rem;
    }

    .token-remove:hover, .remove-all:hover {
        background-color: #0062cc;
        cursor: pointer;
    }

    .actions {
        align-items: center;
        display: flex;
        flex: 1;
        min-width: 15rem;
    }

    input {
        border: none;
        color: white;
        line-height: 1.5rem;
        margin: 0;
        outline: none;
        padding: 0;
        width: 100%;
    }

    .dropdown-arrow path {
        fill: hsl(0, 0%, 70%);
    }

    .multiselect:hover .dropdown-arrow path {
        fill: hsl(0, 0%, 50%);
    }

    .icon-clear path {
        fill: white;
    }

    .options {
        box-shadow: 0 2px 4px rgba(0, 0, 0, .1), 0 -2px 4px rgba(0, 0, 0, .1);
        left: 0;
        list-style: none;
        margin-block-end: 0;
        margin-block-start: 0;
        max-height: 300px;
        overflow-y: scroll;
        overflow-x: hidden;
        padding-inline-start: 0;
        position: absolute;
        top: calc(100% + 1px);
        width: 100%;
    }

    li {
        background-color: #2e3136;
        color: white;
        cursor: pointer;
        padding: .5rem;
    }

    li.selected {
        background-color: #121212;
    }

    li.selected:hover {
        background-color: #121212;
    }

    li:hover {
        background-color: #272727;
    }

    li:last-child {
        border-bottom-left-radius: .2rem;
        border-bottom-right-radius: .2rem;
    }

    .hidden {
        display: none;
    }

    .search {
        background-color: #2e3136;
        caret-color: white;
    }
</style>

<div class="multiselect">
  <div class="tokens" class:showOptions on:click={handleTokenClick}>
    {selected}
    {#each selected as value}
      value
      <div class="token" data-id="{value}">
        <span>{values[selected]}</span>
        <div class="token-remove" title="Remove {values[selected]}">
          <svg class="icon-clear" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
            <path d="{iconClearPath}"/>
          </svg>
        </div>
      </div>
    {/each}
    <div class="actions">
      <input class="search" autocomplete="off" bind:value={inputValue} bind:this={input}
             on:keyup={handleKeyup} on:blur={handleBlur}>
      <svg class="dropdown-arrow" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 18 18">
        <path d="M5 8l4 4 4-4z"></path>
      </svg>
    </div>
  </div>

  <select bind:this={slot} type="multiple" class="hidden">
    <slot/>
  </select>

  {#if showOptions}
    <ul class="options" transition:fly="{{duration: 200, y: 5}}">
      {#each filtered as option}
        <li class:selected={isSelected[option[0]]} data-value="{option[0]}" on:click={toggle(option[0])}>{option[1]}</li>
      {/each}
    </ul>
  {/if}
</div>