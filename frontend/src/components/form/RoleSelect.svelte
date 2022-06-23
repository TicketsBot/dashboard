{#if label !== undefined}
  <label class="form-label">{label}</label>
{/if}

<div class="multiselect-super">
  <Select placeholder="Search..." optionIdentifier="id" items={roles}
          bind:selectedValue={value} getOptionLabel={labelMapper} getSelectionLabel={labelMapper}/>
</div>

<script>
    import Select from 'svelte-select';
    import axios from "axios";
    import {onMount} from 'svelte'
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {API_URL} from "../../js/constants";
    import {notifyError, notifyRatelimit} from "../../js/util";

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
