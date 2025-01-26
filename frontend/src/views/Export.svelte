<div class="content">
    <Card fill={false} footerRight>
        <span slot="title">Data Export</span>
        <div slot="body" class="card-body">
            <p>
                A subset of server data can be exported via the dashboard. In the coming weeks, we will release
                a dedicated data export site, from which all data can be downloaded. Data available via the dashboard
                includes:
            </p>
            <ul>
                <li>Panels & multi-panels</li>
                <li>Staff teams</li>
                <li>Forms</li>
                <li>Tags</li>
                <li>Blacklist data</li>
            </ul>
        </div>
        <div slot="footer">
            <Button icon="fas fa-download" on:click={download}>Download</Button>
        </div>
    </Card>
</div>

<script>
    import Card from "../components/Card.svelte";
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {notifyError} from "../js/util";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    async function download() {
        const res = await axios.get(`${API_URL}/api/${guildId}/export`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        console.log(res.data)
    }
</script>

<style>
    .content {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
    }
</style>
