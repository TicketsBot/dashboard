{#if label !== undefined}
  <label class="form-label">{label}</label>
{/if}

<WrappedSelect {placeholder} optionIdentifier="id" items={roles} {disabled}
               bind:selectedValue={value} nameMapper={labelMapper} on:change />

<script>
    import {onMount} from 'svelte'
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import WrappedSelect from "../WrappedSelect.svelte";

    export let label;
    export let placeholder = "Search...";
    export let roles = [];
    export let guildId;
    export let disabled = false;

    export let value;

    function labelMapper(role) {
        return role.name;
    }

    onMount(() => {
        setDefaultHeaders();
    })
</script>
