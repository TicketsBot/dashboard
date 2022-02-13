<Card footer="{false}" fill="{false}">
  <span slot="title">
    Settings
  </span>

  <div slot="body" class="body-wrapper">
    <form class="settings-form" on:submit|preventDefault={updateSettings}>
      <div class="row">
        <Input label="prefix (max len. 8)" placeholder="t!" col4=true bind:value={data.prefix}/>
        <Number label="per user simultaneous ticket limit" col4=true min=1 max=10 bind:value={data.ticket_limit}/>
        <Checkbox label="allow users to close tickets" col4=true bind:value={data.users_can_close}/>
        <Checkbox label="ticket close confirmation" col4=true bind:value={data.close_confirmation}/>
      </div>
      <div class="row">
        <Textarea label="welcome message" placeholder="Thanks for opening a ticket!" col1=true
                  bind:value={data.welcome_message}/>
      </div>
      <div class="row">
        <ChannelDropdown label="Archive Channel" col3=true channels={channels} withNull={true} bind:value={data.archive_channel}/>
        <CategoryDropdown label="Channel Category" col3=true channels={channels} bind:value={data.category}/>
        <Dropdown label="Overflow Category" col3=true bind:value={data.overflow_category_id}>
          <option value=-1>Disabled</option>
          <option value=-2>Uncategorised (Appears at top of channel list)</option>
          {#each channels as channel}
            {#if channel.type === 4}
              <option value={channel.id}>
                {channel.name}
              </option>
            {/if}
          {/each}
        </Dropdown>
      </div>
      <div class="row">
        <NamingScheme col4=true bind:value={data.naming_scheme}/>
        <Checkbox label="Enable User Feedback" col4=true bind:value={data.feedback_enabled}/>
        <Checkbox label="Hide Claim Button" col4=true bind:value={data.hide_claim_button}/>
        <Checkbox label="Disable /open Command" col4=true bind:value={data.disable_open_command}/>
      </div>
      <div class="row">
        <Checkbox label="Store Ticket Transcripts" col4=true bind:value={data.store_transcripts}/>
      </div>
      <div class="from-message-settings">
        <h3>Start Ticket From Message Settings</h3>
        <div class="row">
          <Dropdown col3={true} label="Required Permission Level" bind:value={data.context_menu_permission_level}>
            <option value="0">Everyone</option>
            <option value="1">Support Representative</option>
            <option value="2">Administrator</option>
          </Dropdown>

          <Checkbox label="Add Message Sender To Ticket" col3={true} bind:value={data.context_menu_add_sender}/>
          <SimplePanelDropdown label="Use Settings From Panel" col3={true} allowNone={true} bind:panels
                               bind:value={data.context_menu_panel}/>
        </div>
      </div>
      <div class="row">
        <div class="col-1">
          <Button icon="fas fa-paper-plane" fullWidth=true>Submit</Button>
        </div>
      </div>
    </form>
  </div>
</Card>

<script>
    import ChannelDropdown from "../ChannelDropdown.svelte";
    import Card from "../Card.svelte";
    import Input from "../form/Input.svelte";
    import Number from "../form/Number.svelte";
    import Checkbox from "../form/Checkbox.svelte";
    import Textarea from "../form/Textarea.svelte";

    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import axios from "axios";
    import {notify, notifyError, notifySuccess, withLoadingScreen} from "../../js/util";
    import {API_URL} from "../../js/constants";
    import CategoryDropdown from "../CategoryDropdown.svelte";
    import Button from "../Button.svelte";
    import NamingScheme from "../NamingScheme.svelte";
    import Dropdown from "../form/Dropdown.svelte";
    import SimplePanelDropdown from "../SimplePanelDropdown.svelte";

    export let guildId;

    setDefaultHeaders();

    let channels = [];
    let panels = [];

    async function loadPanels() {
        const res = await axios.get(`${API_URL}/api/${guildId}/panels`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        panels = res.data;
    }

    async function loadChannels() {
        const res = await axios.get(`${API_URL}/api/${guildId}/channels`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        channels = res.data;
    }

    let data = {
        ticket_limit: 5,
        users_can_close: true,
        close_confirmation: true,
    };

    async function updateSettings() {
        // Svelte hack - I can't even remember what this does
        let mapped = Object.fromEntries(Object.entries(data).map(([k, v]) => {
            if (v === "null") {
                return [k, null];
            } else {
                return [k, v];
            }
        }));

        // "Normalise" data.overflow_category_id
        // Svelte doesn't always keep its promise of using integers, so == instead of ===
        if (mapped.overflow_category_id == -1) {
            mapped.overflow_enabled = false;
            mapped.overflow_category_id = null;
        } else if (mapped.overflow_category_id == -2) {
            mapped.overflow_enabled = true
            mapped.overflow_category_id = null;
        } else {
            mapped.overflow_enabled = true;
        }

        const res = await axios.post(`${API_URL}/api/${guildId}/settings`, mapped);
        if (res.status === 200) {
            if (showValidations(res.data)) {
                notifySuccess('Your settings have been saved.');
            } else {
                // Load valid data
                await loadData();
            }
        } else {
            notifyError(res.data.error);
        }
    }

    function showValidations(data) {
        let success = true;

        if (!data.prefix) {
            success = false;
            notify("Warning", "Your prefix has not been saved.\nPrefixes must be between 1 - 8 characters in length.")
        }

        if (!data.welcome_message) {
            success = false;
            notify("Warning", "Your welcome message has not been saved.\nWelcome messages must be between 1 - 1000 characters in length.")
        }

        if (!data.ticket_limit) {
            success = false;
            notify("Warning", "Your ticket limit has not been saved.\nTicket limits must be in the range 1 - 10.")
        }

        if (!data.archive_channel) {
            success = false;
            notify("Warning", "Your archive channel has not been saved.")
        }

        if (!data.category) {
            success = false;
            notify("Warning", "Your channel category has not been saved.")
        }

        if (!data.naming_scheme) {
            success = false;
            notify("Warning", "Your archive channel has not been saved.")
        }

        return success;
    }

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/${guildId}/settings`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data = res.data;

        // Overrides
        if (data.archive_channel === "0") {
            let first = channels.find((c) => c.type === 0);
            if (first !== undefined) {
                data.archive_channel = first.id;
            }
        }

        if (data.category === "0") {
            let first = channels.find((c) => c.type === 4);
            if (first !== undefined) {
                data.category = first.id;
            }
        }

        if (data.overflow_enabled === false) {
          data.overflow_category_id = "-1";
        } else if (data.overflow_enabled === true) {
          if (data.overflow_category_id === null) {
            data.overflow_category_id = "-2";
          }
        }
    }

    withLoadingScreen(async () => {
        await loadPanels();
        await loadChannels();
        await loadData();
    });
</script>

<style>
    :global(.body-wrapper) {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-top: 10px;
    }

    .settings-form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .from-message-settings {
        border-top: 1px solid rgba(0, 0, 0, .25);
        margin-top: 25px;
        padding-top: 10px;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
            justify-content: center;
        }

        :global(.col-4, .col-3, .col-2, .col-3-4) {
            width: 100% !important;
        }
    }

    :global(.col-1) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 100%;
        height: 100%;
    }

    :global(.col-2) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 49%;
        height: 100%;
    }

    :global(.col-3) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 31%;
        height: 100%;
    }

    :global(.col-4) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 23%;
        height: 100%;
    }
</style>
