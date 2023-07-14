{#if label !== undefined}
  <label class="form-label">{label}</label>
{/if}

<WrappedSelect placeholder="Search..." optionIdentifier="id" items={roles}
               bind:selectedValue={value} nameMapper={labelMapper}/>

<script>
    import {onMount} from 'svelte'
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import WrappedSelect from "../WrappedSelect.svelte";

    export let label;
    export let roles = [];
    export let guildId;

    export let value;

    function labelMapper(role) {
        return role.name;
    }

    onMount(() => {
        setDefaultHeaders();
    })
</script>
