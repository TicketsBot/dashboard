{#if label !== undefined}
  <label class="form-label">{label}</label>
{/if}

<WrappedSelect placeholder="Search..." loadOptionsInterval={500} {loadOptions} optionIdentifier="id"
               bind:selectedValue={value} nameMapper={labelMapper}/>

<script>
    import axios from "axios";
    import {onMount} from 'svelte'
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {API_URL} from "../../js/constants";
    import {notifyError, notifyRatelimit} from "../../js/util";
    import WrappedSelect from "../WrappedSelect.svelte";

    export let label;
    export let guildId;

    export let value;

    async function loadOptions(filterText) {
        const res = await axios.get(`${API_URL}/api/${guildId}/members/search?query=${filterText}`)
        if (res.status !== 200) {
            if (res.status === 429) {
                notifyRatelimit();
            } else {
                notifyError(res.data.error);
            }

            return {cancelled: true}
        }

        return res.data.map((m) => m.user);
    }

    function labelMapper(user) {
        if (!user.discriminator || user.discriminator === "0" || user.discriminator === "0000") {
            return user.username;
        } else {
            return `${user.username}#${user.discriminator}`
        }
    }

    onMount(() => {
        setDefaultHeaders();
    })
</script>
