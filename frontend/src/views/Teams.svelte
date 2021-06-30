<div class="parent">
  <div class="content">
    <Card footer={false}>
      <span slot="title">Support Teams</span>
      <div slot="body" class="body-wrapper">
        <div class="section">
          <h2 class="section-title">Create Team</h2>

          <form on:submit|preventDefault={createTeam}>
            <div class="row" style="max-height: 40px"> <!-- hacky -->
              <Input placeholder="Team Name" col4={true} bind:value={createName}/>
              <div style="margin-left: 30px">
                <Button icon="fas fa-paper-plane">Submit</Button>
              </div>
            </div>
          </form>
        </div>
        <div class="section">
          <h2 class="section-title">Manage Teams</h2>

          <div class="col-1" style="flex-direction: row">
            <div class="col-4" style="margin-right: 12px">
              <div class="multiselect-super">
                <Select isSearchable={false} isClearable={false} optionIdentifier="id" items={teams}
                        bind:selectedValue={activeTeam} getOptionLabel={labelMapper} getSelectionLabel={labelMapper}
                        on:select={updateActiveTeam}/>
              </div>
            </div>

            {#if activeTeam.id !== 'default'}
              <div class="col-1">
                <Button danger={true} type="button"
                        on:click={() => deleteTeam(activeTeam.id)}>Delete {activeTeam.name}</Button>
              </div>
            {/if}
          </div>

          <div class="manage">
            <div class="col">
              <h3>Manage Members</h3>

              <table class="nice">
                <tbody>
                {#each members as member}
                  <tr>
                    <td>{member.name}</td>
                    <td style="display: flex; flex-direction: row-reverse">
                      <Button type="button" danger={true} on:click={() => removeMember(activeTeam.id, member)}>Delete
                      </Button>
                    </td>
                  </tr>
                {/each}
                </tbody>
              </table>
            </div>

            <div class="col">
              <h3>Add Member</h3>
              <div class="user-select">
                <div style="display: flex; flex: 1">
                  <UserSelect {guildId} bind:value={selectedUser}/>
                </div>

                <div style="margin-left: 10px">
                  <Button type="button" icon="fas fa-plus" disabled={selectedUser === undefined}
                          on:click={addUser}>Add To Team
                  </Button>
                </div>
              </div>

              <h3>Add Role</h3>
              <div class="user-select">
                <div style="display: flex; flex: 1">
                  <RoleSelect {guildId} {roles} bind:value={selectedRole}/>
                </div>

                <div style="margin-left: 10px">
                  <Button type="button" icon="fas fa-plus" disabled={selectedRole === undefined}
                          on:click={addRole}>Add To Team
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Card>
  </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, notifyRatelimit, notifySuccess, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Input from "../components/form/Input.svelte";
    import Select from 'svelte-select';
    import UserSelect from "../components/form/UserSelect.svelte";
    import RoleSelect from "../components/form/RoleSelect.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let defaultTeam = {id: 'default', name: 'Default'};

    let createName;
    let teams = [];
    let roles = [];
    let activeTeam = defaultTeam;
    let members = [];

    let selectedUser;
    let selectedRole;

    function labelMapper(team) {
        return team.name;
    }

    async function updateActiveTeam() {
        const res = await axios.get(`${API_URL}/api/${guildId}/team/${activeTeam.id}`);
        if (res.status !== 200) {
            if (res.status === 429) {
                notifyRatelimit();
            } else {
                notifyError(res.data.error);
            }
            return;
        }

        members = res.data;
    }

    async function addUser() {
        const res = await axios.put(`${API_URL}/api/${guildId}/team/${activeTeam.id}/${selectedUser.id}?type=0`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${selectedUser.username}#${selectedUser.discriminator} has been added to the support team ${activeTeam.name}`);

        let entity = {
            id: selectedUser.id,
            type: 0,
            name: `${selectedUser.username}#${selectedUser.discriminator}`
        }
        members = [...members, entity];
        selectedUser = undefined;
    }

    async function addRole() {
        const res = await axios.put(`${API_URL}/api/${guildId}/team/${activeTeam.id}/${selectedRole.id}?type=1`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${selectedRole.name} has been added to the support team ${activeTeam.name}`);

        let entity = {
            id: selectedRole.id,
            type: 1,
            name: selectedRole.name,
        }
        members = [...members, entity];
        selectedRole = undefined;
    }

    async function removeMember(teamId, entity) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/team/${teamId}/${entity.id}?type=${entity.type}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${entity.name} has been removed from the team`);
        members = members.filter((member) => member.id !== entity.id);
    }

    async function createTeam() {
        let data = {
            name: createName,
        };

        const res = await axios.post(`${API_URL}/api/${guildId}/team`, data);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Team ${createName} has been created`);
        createName = '';
        teams = [...teams, res.data];
    }

    async function deleteTeam(id) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/team/${id}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`Team deleted successfully`);
        activeTeam = defaultTeam;
        teams = teams.filter((team) => team.id !== id);
    }

    async function loadTeams() {
        const res = await axios.get(`${API_URL}/api/${guildId}/team`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        teams = [defaultTeam, ...res.data];
    }

    async function loadRoles() {
        const res = await axios.get(`${API_URL}/api/${guildId}/roles`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        roles = res.data.roles;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadTeams();
        await loadRoles();
        await updateActiveTeam();
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

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        padding: 1%;
    }

    .section {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .section:not(:first-child) {
        margin-top: 2%;
    }

    .section-title {
        font-size: 36px;
        font-weight: bolder !important;
    }

    h3 {
        font-size: 28px;
        margin-bottom: 4px;
    }

    .row {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
    }

    .manage {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-top: 12px;
    }

    .col {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 49%;
        height: 100%;
    }

    table.nice > tbody > tr:first-child {
        border-top: 1px solid #dee2e6;
    }

    .user-select {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-bottom: 1%;
    }

    @media only screen and (max-width: 950px) {
        .manage {
            flex-direction: column;
        }

        .col {
            width: 100%;
        }
    }
</style>
