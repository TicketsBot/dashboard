{#if blacklistUserModal}
    <div class="modal" transition:fade>
        <div class="modal-wrapper">
            <Card footer footerRight fill={false}>
                <span slot="title">Blacklist User</span>

                <div slot="body" class="modal-inner">
                    <div>
                        <label class="form-label" style="margin-bottom: 0 !important;">Use User ID</label>
                        <Toggle hideLabel
                                toggledColor="#66bb6a"
                                untoggledColor="#ccc"
                                bind:toggled={blacklistById}/>
                    </div>

                    {#if blacklistById}
                        <Input label="User ID" placeholder="592348585904198711" bind:value={blacklistUserId}/>
                    {:else}
                        <div class="user-select-wrapper">
                            <UserSelect {guildId} label="User" bind:value={blacklistUser}/>
                        </div>
                    {/if}
                </div>

                <div slot="footer" style="gap: 12px">
                    <Button danger on:click={() => blacklistUserModal = false}>Cancel</Button>
                    <Button on:click={addUser}>Confirm</Button>
                </div>
            </Card>
        </div>
    </div>

    <div class="modal-backdrop" transition:fade>
    </div>
{:else if blacklistRoleModal}
    <div class="modal" transition:fade>
        <div class="modal-wrapper">
            <Card footer footerRight fill={false}>
                <span slot="title">Blacklist Role</span>

                <div slot="body" class="modal-inner user-select-wrapper">
                    <RoleSelect {guildId} {roles} label="Role" bind:value={blacklistRole}/>
                </div>

                <div slot="footer" style="gap: 12px">
                    <Button danger on:click={() => blacklistRoleModal = false}>Cancel</Button>
                    <Button on:click={addRole}>Confirm</Button>
                </div>
            </Card>
        </div>
    </div>

    <div class="modal-backdrop" transition:fade>
    </div>
{/if}

{#if data}
    <div class="content">
        <div class="main-col">
            <Card footer={false}>
                <span slot="title">Blacklist</span>
                <div slot="body" class="body-wrapper">
                    <div class="row" style="gap: 10px">
                        <Button icon="fas fa-ban" on:click={() => blacklistUserModal = true}>Blacklist New User</Button>
                        <Button icon="fas fa-ban" on:click={() => blacklistRoleModal = true}>Blacklist New Role</Button>
                    </div>

                    <hr/>

                    <div class="tables">
                        <table class="nice">
                            <thead>
                            <tr>
                                <th class="full-width">Role</th>
                                <th>Remove</th>
                            </tr>
                            </thead>
                            <tbody>
                            {#each data.roles as roleId}
                                {@const role = roles.find(role => role.id === roleId)}
                                <tr>
                                    {#if role === undefined}
                                        <td class="full-width">Unknown ({roleId})</td>
                                    {:else}
                                        <td class="full-width">{role.name}</td>
                                    {/if}

                                    <td>
                                        <Button type="button" danger icon="fas fa-trash-can"
                                                on:click={() => removeRoleBlacklist(roleId, role)}>
                                            Remove
                                        </Button>
                                    </td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>

                        <table class="nice">
                            <thead>
                            <tr>
                                <th class="full-width">User</th>
                                <th>Remove</th>
                            </tr>
                            </thead>
                            <tbody>
                            {#each data.users as user}
                                <tr>
                                    {#if user.username !== ''}
                                        <td class="full-width">{user.username} ({user.id})</td>
                                    {:else}
                                        <td class="full-width">Unknown ({user.id})</td>
                                    {/if}

                                    <td>
                                        <Button type="button" danger icon="fas fa-trash-can"
                                                on:click={() => removeUserBlacklist(user)}>
                                            Remove
                                        </Button>
                                    </td>
                                </tr>
                            {/each}
                            </tbody>
                        </table>
                    </div>

                    <div class="row nav">
                        <i class="fas fa-chevron-left pagination-chevron" class:disabled-chevron={page <= 1}
                           on:click={loadPrevious}></i>
                        <span>Page {page}</span>
                        <i class="fas fa-chevron-right pagination-chevron"
                           class:disabled-chevron={data.users.length < data.page_limit && data.roles.length < data.page_limit}
                           on:click={loadNext}></i>
                    </div>
                </div>
            </Card>
        </div>
    </div>
{/if}

<script>
    import Card from "../components/Card.svelte";
    import UserSelect from "../components/form/UserSelect.svelte";
    import {notifyError, notifySuccess, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {fade} from "svelte/transition";
    import Toggle from "svelte-toggle";
    import Input from "../components/form/Input.svelte";
    import RoleSelect from "../components/form/RoleSelect.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let page = 1;
    let data;
    let roles = [];

    let blacklistUserModal = false;
    let blacklistRoleModal = false;
    let blacklistById = false;
    let blacklistUserId;
    let blacklistUser;
    let blacklistRole;

    function loadPrevious() {
        if (page > 1) {
            page--;
            loadData();
        }
    }

    function loadNext() {
        if (data.users.length >= data.page_limit || data.roles.length >= data.page_limit) {
            page++;
            loadData();
        }
    }

    async function addUser() {
        let snowflake;
        if (blacklistById) {
            snowflake = blacklistUserId;
        } else {
            snowflake = blacklistUser.id;
        }

        const body = {
            entity_type: 0,
            snowflake: snowflake
        };

        const res = await axios.post(`${API_URL}/api/${guildId}/blacklist`, body);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        if (res.data.resolved) {
            notifySuccess(`${res.data.username} has been blacklisted`);

            data.users = [...data.users, {
                id: res.data.id,
                username: res.data.username,
            }];
        } else {
            notifySuccess(`User with ID ${res.data.id} has been blacklisted`);
            data.users = [...data.users, {
                id: res.data.id,
                username: "Unknown",
            }];
        }

        blacklistById = false;
        blacklistUser = undefined;
        blacklistUserId = undefined;
        blacklistUserModal = false;
    }

    async function addRole() {
        const body = {
            entity_type: 1,
            snowflake: blacklistRole.id,
        };

        const res = await axios.post(`${API_URL}/api/${guildId}/blacklist`, body);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data.roles = [...data.roles, blacklistRole.id];

        notifySuccess(`${blacklistRole.name} has been blacklisted`);
        blacklistRole = undefined;
        blacklistRoleModal = false;
    }

    async function removeUserBlacklist(user) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/blacklist/user/${user.id}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${user.username || `User with ID ${user.id}`} has been removed from the blacklist`);
        data.users = data.users.filter((u) => u.id !== user.id);
    }

    async function removeRoleBlacklist(roleId, role) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/blacklist/role/${roleId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        if (role) {
            notifySuccess(`${role.name} has been removed from the blacklist`);
        } else {
            notifySuccess(`Role with ID ${roleId} has been removed from the blacklist`);
        }

        data.roles = data.roles.filter((otherId) => otherId !== roleId);
    }

    async function loadRoles() {
        const res = await axios.get(`${API_URL}/api/${guildId}/roles`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        roles = res.data.roles;
    }

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/${guildId}/blacklist?page=${page}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();

        await Promise.all([
            loadData(),
            loadRoles()
        ]);
    });
</script>

<style>
    .content {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .main-col {
        display: flex;
        flex-direction: column;
        width: 100%;
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
    }

    hr {
        border-top: 1px solid #777;
        border-bottom: 0;
        border-left: 0;
        border-right: 0;
        width: 100%;
        flex: 1;
    }

    .tables {
        display: flex;
        flex-direction: column;
        row-gap: 4vh;
    }

    .full-width {
        width: 100%;
    }

    .nav {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        gap: 2px;
        margin-top: 20px;
    }

    .pagination-chevron {
        cursor: pointer;
        color: #3472f7;
    }

    .disabled-chevron {
        color: #777 !important;
        cursor: default !important;
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

    .modal {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 999;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-wrapper {
        display: flex;
        width: 60%;
        margin: 10% auto auto auto;
    }

    .modal-inner {
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
        gap: 2%;
        width: 100%;
    }

    .user-select-wrapper {
        display: flex;
        flex-direction: column;
        width: 33%;
    }

    @media only screen and (max-width: 1280px) {
        .modal-wrapper {
            width: 96%;
        }
    }

    @media only screen and (max-width: 950px) {
        .content {
            width: 96%;
        }
    }

    @media only screen and (max-width: 700px) {
        .user-select-wrapper {
            width: 100%;
        }
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 500;
        background-color: #000;
        opacity: .5;
    }
</style>
