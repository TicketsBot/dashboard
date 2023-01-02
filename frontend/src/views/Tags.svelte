{#if tagCreateModal}
  <TagEditor on:cancel={() => tagCreateModal = false} on:confirm={createTag}/>
{:else if tagEditModal}
  <TagEditor bind:data={editData} on:cancel={cancelEdit} on:confirm={editTag}/>
{/if}

<div class="parent">
  <div class="content">
    <div class="main-col">
      <Card footer footerRight>
        <span slot="title">Tags</span>
        <div slot="body" class="body-wrapper">
          <table class="nice">
            <thead>
            <tr>
              <th>Tag</th>
              <th style="text-align: right">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each Object.entries(tags) as [id, tag]}
              <tr>
                <td>{id}</td>
                <td class="actions">
                  <Button type="button" on:click={() => openEditModal(id)}>Edit</Button>
                  <Button type="button" danger={true} on:click={() => deleteTag(id)}>Delete</Button>
                </td>
              </tr>
            {/each}
            </tbody>
          </table>
        </div>
        <div slot="footer">
          <Button icon="fas fa-plus" on:click={openCreateModal}>Create Tag</Button>
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
    import {fade} from "svelte/transition";
    import TagEditor from "../components/manage/TagEditor.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let tags = {};
    let editData;
    let editId;

    let tagCreateModal = false;
    let tagEditModal = false;

    function openCreateModal(id) {
        tagCreateModal = true;
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    function openEditModal(id) {
        editId = id;
        editData = tags[id];
        tagEditModal = true;

        window.scrollTo({ top: 0, behavior: 'smooth' });
    }

    function cancelEdit() {
        editId = undefined;
        editData = undefined;
        tagEditModal = false;
    }

    async function createTag(e) {
        const data = e.detail;
        if (!data.id || data.id.length === 0) {
            notifyError("Tag ID is required");
            return;
        }

        if (data.content !== null && data.content !== undefined && data.content.length === 0) {
            data.content = null;
        }

        const res = await axios.put(`${API_URL}/api/${guildId}/tags`, data);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Tag ${data.id} has been created`);
        tagCreateModal = false;
        tags[data.id] = data;
    }

    async function editTag(e) {
        const data = e.detail;

        if (editId !== data.id) {
            // Delete old tag
            const res = await axios.delete(`${API_URL}/api/${guildId}/tags`, {data: {tag_id: editId}});
            if (res.status !== 204) {
                notifyError(res.data.error);
                return;
            }

            delete tags[editId];
        }

        if (!data.id || data.id.length === 0) {
            notifyError("Tag ID is required");
            return;
        }

        if (data.content !== null && data.content !== undefined && data.content.length === 0) {
            data.content = null;
        }

        const res = await axios.put(`${API_URL}/api/${guildId}/tags`, data);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        tags[data.id] = data;

        tagEditModal = false;
        editData = undefined;
        editId = undefined;

        notifySuccess("Tag edited successfully");
    }

    async function deleteTag(id) {
        const data = {
            tag_id: id
        };

        const res = await axios.delete(`${API_URL}/api/${guildId}/tags`, {data: data});
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Tag deleted successfully`);
        delete tags[id];
        tags = tags;
    }

    async function loadTags() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tags`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        tags = res.data;
        for (const id in tags) {
            tags[id].use_embed = tags[id].embed !== null;
        }
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
        padding-bottom: 4%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: 2vh;
        width: 100%;
        height: 100%;
    }

    table {
        width: 100%;
    }

    .actions {
        display: flex;
        flex-direction: row;
        gap: 10px;
        justify-content: flex-end;
    }

    @media only screen and (max-width: 950px) {
        .content {
            flex-direction: column-reverse;
        }

        .main-col {
            width: 100%;
            margin-top: 4%;
        }
    }
</style>
