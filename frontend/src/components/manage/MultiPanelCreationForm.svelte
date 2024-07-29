<form on:submit|preventDefault>
    <Collapsible defaultOpen>
        <span slot="header">Properties</span>
        <div slot="content" class="col-1">
            <div class="col-1">
                <ChannelDropdown col1 allowAnnouncementChannel {channels} label="Panel Channel"
                                 bind:value={data.channel_id}/>
            </div>
            <div class="col-1" style="padding-right: 10px">
                <PanelDropdown label="Panels (Minimum 2)" {panels} bind:selected={data.panels}/>
            </div>
            <div class="col-1">
                <div class="row dropdown-menu-settings">
                    <Checkbox label="Use Dropdown Menu" bind:value={data.select_menu}/>
                    <div class="placeholder-input">
                        <Input label="Dropdown Menu Placeholder" col1 placeholder="Select a topic..."
                               bind:value={data.select_menu_placeholder} disabled={!data.select_menu} />
                    </div>
                </div>
            </div>
        </div>
    </Collapsible>

    <Collapsible defaultOpen>
        <span slot="header">Message</span>
        <div slot="content" class="col-1">
            <EmbedForm footerPremiumOnly={true} bind:data={data.embed}/>
        </div>
    </Collapsible>
</form>

<script>
    import ChannelDropdown from "../ChannelDropdown.svelte";
    import PanelDropdown from "../PanelDropdown.svelte";
    import Checkbox from "../form/Checkbox.svelte";
    import Collapsible from "../Collapsible.svelte";
    import EmbedForm from "../EmbedForm.svelte";
    import Input from "../form/Input.svelte";

    export let data;

    export let guildId;
    export let channels = [];
    export let panels = [];

    export let seedDefault = true;
    if (seedDefault) {
        const firstChannel = channels[0];

        data = {
            channels: firstChannel ? firstChannel.id : undefined,
            panels: [],
            embed: {
                title: 'Open a ticket!',
                fields: [],
                colour: 0x2ECC71,
                author: {},
                footer: {},
            },
        }
    }
</script>

<style>
    form {
        display: flex;
        flex-direction: column;
        width: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
    }

    .dropdown-menu-settings {
        gap: 10px;
        margin-top: 10px;
    }

    .dropdown-menu-settings > .placeholder-input {
        flex: 1;
    }

    @media only screen and (max-width: 950px) {
        .row {
            flex-direction: column;
        }
    }

    :global(.col-1-4) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 25%;
        height: 100%;
    }

    :global(.col-3-4) {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 75%;
        height: 100%;
    }
</style>
