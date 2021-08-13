<div class="parent">
  <div class="content">
    <div class="main-col">
      <Card footer={false}>
        <span slot="title">Tags</span>
        <div slot="body" class="body-wrapper">
          <table class="nice">
            <thead>
            <tr>
              <th>Tag</th>
              <th>Edit</th>
              <th>Delete</th>
            </tr>
            </thead>
            <tbody>
            {#each Object.entries(tags) as [id, content]}
              <tr>
                <td>{id}</td>
                <td>
                  <Button type="button" on:click={() => editTag(id)}>Edit</Button>
                </td>
                <td>
                  <Button type="button" danger={true} on:click={() => deleteTag(id)}>Delete</Button>
                </td>
              </tr>
            {/each}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
    <div class="right-col">
      <Card footer={false}>
        <span slot="title">Create A Tag</span>
        <div slot="body" class="body-wrapper">
          <form class="body-wrapper" on:submit|preventDefault={createTag}>
            <div class="col" style="flex-direction: column">
              <Input label="Tag ID" placeholder="mytag" bind:value={createData.id}/>
            </div>
            <div class="col" style="flex-direction: column">
              <Textarea label="Tag Content" placeholder="Enter the text that the bot should respond with"
                        bind:value={createData.content}/>
            </div>

            <div class="row" style="justify-content: center">
              <div class="col-2">
                <Button fullWidth={true} icon="fas fa-plus">Submit</Button>
              </div>
            </div>
          </form>
        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, notifySuccess, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Input from "../components/form/Input.svelte";
    import Textarea from "../components/form/Textarea.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let createData = {};
    let tags = {};

    function editTag(id) {
        createData.id = id;
        createData.content = tags[id];
    }

    async function createTag() {
        const res = await axios.put(`${API_URL}/api/${guildId}/tags`, createData);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Tag ${createData.id} has been created`);
        tags[createData.id] = createData.content;
        createData = {};
    }

    async function deleteTag(id) {
        const data = {
            tag_id: id
        };

        const res = await axios.delete(`${API_URL}/api/${guildId}/tags`, {data: data});
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Tag deleted successfully`);
        delete tags[id];
        tags = tags; // svelte terrible
    }

    async function loadTags() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tags`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        tags = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadTags();
    });
</script>

<style>
    .parent {
        display: flex;
        justify-content: center;
        width: 100%;
        height: 100%;
    }

    .content {
        display: flex;
        justify-content: space-between;
        width: 96%;
        height: 100%;
        margin-top: 30px;
    }

    .main-col {
        display: flex;
        flex-direction: column;
        width: 64%;
        height: 100%;
    }

    .right-col {
        display: flex;
        flex-direction: column;
        width: 34%;
        height: 100%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
        margin-bottom: 2%;
    }

    .col {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        margin-bottom: 2%;
    }

    @media only screen and (max-width: 950px) {
        .content {
            flex-direction: column-reverse;
        }

        .main-col {
            width: 100%;
            margin-top: 4%;
        }

        .right-col {
            width: 100%;
        }
    }
</style>
