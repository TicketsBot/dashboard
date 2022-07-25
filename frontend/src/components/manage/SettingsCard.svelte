{#if data}
  <Card footer="{false}" fill="{false}">
  <span slot="title">
    Settings
  </span>

    <div slot="body" class="body-wrapper">
      <form class="settings-form" on:submit|preventDefault={updateSettings}>
        <Collapsible defaultOpen>
          <span slot="header">General</span>
          <div slot="content" class="col-1">
            <div class="row">
              <Input label="prefix (max len. 8)" placeholder="t!" col4 bind:value={data.prefix}/>
              <Number label="per user simultaneous ticket limit" min=1 max=10 bind:value={data.ticket_limit}/>
              <Dropdown label="Language" bind:value={data.language}>
                <option value=null selected="selected">Server Default</option>
                {#if data.languages}
                  {#each data.languages as language}
                    <option value={language}>{data.language_names[language]}</option>
                  {/each}
                {/if}
              </Dropdown>
              <Checkbox label="allow users to close tickets" bind:value={data.users_can_close}/>
              <Checkbox label="ticket close confirmation" bind:value={data.close_confirmation}/>
              <Checkbox label="Enable User Feedback" bind:value={data.feedback_enabled}/>
            </div>
          </div>
        </Collapsible>

        <Collapsible defaultOpen>
          <span slot="header">Tickets</span>
          <div slot="content" class="col-1">
            <div class="row">
              <ChannelDropdown label="Archive Channel" col3=true channels={channels} withNull={true}
                               bind:value={data.archive_channel}/>
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
              <Checkbox label="Store Ticket Transcripts" bind:value={data.store_transcripts}/>
              <Checkbox label="Hide Claim Button" bind:value={data.hide_claim_button}/>
            </div>

            <div class="row">
            </div>
          </div>
        </Collapsible>

        <Collapsible>
          <span slot="header">/Open Command</span>
          <div slot="content" class="col-1">
            <div class="row">
              <Checkbox label="Disable /open Command" bind:value={data.disable_open_command}/>
              <CategoryDropdown label="Channel Category" col3 channels={channels} bind:value={data.category}/>
              <NamingScheme bind:value={data.naming_scheme}/>
            </div>
            <div class="row">
              <div class="col-1-flex">
              <Textarea label="welcome message" placeholder="Thanks for opening a ticket!" col1
                        bind:value={data.welcome_message}/>
              </div>
            </div>
          </div>
        </Collapsible>

        <Collapsible>
          <span slot="header">Context Menu (Start Ticket Dropdown)</span>
          <div slot="content" class="col-1">
            <div class="row">
              <Dropdown col3 label="Required Permission Level" bind:value={data.context_menu_permission_level}>
                <option value="0">Everyone</option>
                <option value="1">Support Representative</option>
                <option value="2">Administrator</option>
              </Dropdown>

              <Checkbox label="Add Message Sender To Ticket" bind:value={data.context_menu_add_sender}/>
              <SimplePanelDropdown label="Use Settings From Panel" col3 allowNone={true} bind:panels
                                   bind:value={data.context_menu_panel}/>
            </div>
          </div>
        </Collapsible>

        <Collapsible>
          <span slot="header">Claiming</span>
          <div slot="content" class="col-1">
            <div class="row">
              <Checkbox label="SUPPORT REPS CAN VIEW CLAIMED TICKETS" bind:value={data.claim_settings.support_can_view}
                        on:change={validateView}/>
              <Checkbox label="SUPPORT REPS CAN TYPE IN CLAIMED TICKETS" bind:value={data.claim_settings.support_can_type}
                        on:change={validateType}/>
            </div>
          </div>
        </Collapsible>

        <Collapsible>
          <span slot="header">Auto Close</span>
          <div slot="content" class="col-1">
            <div class="row">
              <Checkbox label="Enabled" bind:value={data.auto_close.enabled}/>
              <Checkbox label="Close On User Leave" disabled={!data.auto_close.enabled} bind:value={data.auto_close.on_user_leave}/>
            </div>

            <div class="row" style="justify-content: space-between">
              <div class="col-2" style="flex-direction: row">
                <Duration disabled={!isPremium || !data.auto_close.enabled} bind:days={sinceOpenDays} bind:hours={sinceOpenHours}
                          bind:minutes={sinceOpenMinutes}>
                  <div slot="header" class="header">
                    <label class="form-label" style="margin-bottom: unset">Since Open With No Response</label>
                    <PremiumBadge/>
                  </div>
                </Duration>
              </div>
              <div class="col-2" style="flex-direction: row">
                <Duration disabled={!isPremium || !data.auto_close.enabled} bind:days={sinceLastDays} bind:hours={sinceLastHours}
                          bind:minutes={sinceLastMinutes}>
                  <div slot="header" class="header">
                    <label class="form-label" style="margin-bottom: unset">Since Last Message</label>
                    <PremiumBadge/>
                  </div>
                </Duration>
              </div>
            </div>
          </div>
        </Collapsible>

        <Collapsible>
          <div slot="header" class="header">
            <span>Colour Scheme</span>
            <PremiumBadge/>
          </div>
          <div slot="content" class="col-1">
            <div class="row">
              <Colour col4 label="Success" bind:value={data.colours["0"]} disabled={!isPremium}/>
              <Colour col4 label="Failure" bind:value={data.colours["1"]} disabled={!isPremium}/>
            </div>
          </div>
        </Collapsible>

        <div class="row">
          <div class="col-1">
            <Button icon="fas fa-paper-plane" fullWidth=true>Submit</Button>
          </div>
        </div>
      </form>
    </div>
  </Card>
{/if}

<svelte:head>
  <style>
      body {
          overflow-y: scroll;
      }
  </style>
</svelte:head>

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
    import Collapsible from "../Collapsible.svelte";
    import Duration from "../form/Duration.svelte";
    import Colour from "../form/Colour.svelte";
    import PremiumBadge from "../PremiumBadge.svelte";
    import {toDays, toHours, toMinutes} from "../../js/timeutil";

    export let guildId;

    setDefaultHeaders();

    let channels = [];
    let panels = [];
    let isPremium = false;

    let data;

    let sinceOpenDays = 0, sinceOpenHours = 0, sinceOpenMinutes = 0;
    let sinceLastDays = 0, sinceLastHours = 0, sinceLastMinutes = 0;

    function validateView() {
        if (!data.support_can_view && data.support_can_type) {
            data.support_can_type = false;
        }
    }

    function validateType() {
        if (!data.support_can_view && data.support_can_type) {
            data.support_can_view = true;
        }
    }

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

    async function loadPremium() {
        const res = await axios.get(`${API_URL}/api/${guildId}/premium`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        isPremium = res.data.premium;
    }

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

        // Normalise autoclose
        data.auto_close.since_open_with_no_response = sinceOpenDays * 86400 + sinceOpenHours * 3600 + sinceOpenMinutes * 60;
        data.auto_close.since_last_message = sinceLastDays * 86400 + sinceLastHours * 3600 + sinceLastMinutes * 60;

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

        if (data.error !== null) {
            success = false;
            notify("Warning", data.error);
        }

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

            if (!channels.some((c) => c.id === data.overflow_category_id)) {
                data.overflow_category_id = "-2";
            }
        }

        if (data.language === null) {
            data.language = "null";
        }

        // Auto close overrides
        if (data.auto_close.since_open_with_no_response) {
            sinceOpenDays = toDays(data.auto_close.since_open_with_no_response);
            sinceOpenHours = toHours(data.auto_close.since_open_with_no_response);
            sinceOpenMinutes = toMinutes(data.auto_close.since_open_with_no_response);
        }

        if (data.auto_close.since_last_message) {
            sinceLastDays = toDays(data.auto_close.since_last_message);
            sinceLastHours = toHours(data.auto_close.since_last_message);
            sinceLastMinutes = toMinutes(data.auto_close.since_last_message);
        }
    }

    withLoadingScreen(async () => {
        await Promise.all([
            loadPanels(),
            loadChannels(),
            loadPremium()
        ]);

        await loadData(); // Depends on channels
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
        justify-content: flex-start;
        flex-wrap: wrap;
        gap: 2%;
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

    .col-1-flex {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        flex: 0 0 100%;
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

    :global(.col-1, .col-1-force) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 100%;
        height: 100%;
    }

    :global(.col-2, .col-2-force) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 49%;
        height: 100%;
    }

    :global(.col-3, .col-3-force) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 31%;
        height: 100%;
    }

    :global(.col-4, .col-4-force) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 23%;
        height: 100%;
    }

    .header {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 4px;
    }
</style>
