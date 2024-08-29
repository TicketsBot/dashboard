<div class="content">
    <Card footer={false}>
        <span slot="title">Open Tickets</span>
        <div slot="body" class="body-wrapper">
            <table class="nice">
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Panel</th>
                    <th>User</th>
                    <th>View</th>
                </tr>
                </thead>
                <tbody>
                {#each tickets as ticket}
                    <tr>
                        <td>{ticket.id}</td>
                        <td>{ticket.panel_title}</td>
                        {#if ticket.user !== undefined}
                            <td>{ticket.user.username}</td>
                        {:else}
                            <td>Unknown</td>
                        {/if}
                        <td>
                            <Navigate to="/manage/{guildId}/tickets/view/{ticket.id}" styles="link">
                                <Button type="button">View</Button>
                            </Navigate>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </Card>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, withLoadingScreen} from '../js/util'
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Button from "../components/Button.svelte";
    import {Navigate} from 'svelte-router-spa';

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let tickets = [];

    async function loadTickets() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tickets`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        tickets = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadTickets();
    });
</script>

<style>
    .content {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }
</style>
