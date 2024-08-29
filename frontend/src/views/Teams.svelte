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

          <div class="col-1" style="flex-direction: row; gap: 12px">
            <Dropdown col3 label="Team" bind:value={activeTeam} on:change={(e) => updateActiveTeam(e.target.value)}>
              {#each teams as team}
                <option value={team.id}>{team.name}</option>
              {/each}
            </Dropdown>

            {#if activeTeam !== 'default'}
              <div style="margin-top: auto; margin-bottom: 0.5em">
                <Button danger={true} type="button"
                    on:click={() => deleteTeam(activeTeam)}>Delete {getTeam(activeTeam)?.name}</Button>
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
                    {#if member.type === USER_TYPE}
                      <td>{member.name}</td>
                    {:else if member.type === ROLE_TYPE}
                      {@const role = roles.find(role => role.id === member.id)}
                      <td>{role === undefined ? "Unknown Role" : role.name}</td>
                    {/if}
                    <td style="display: flex; flex-direction: row-reverse">
                      <Button type="button" danger={true} on:click={() => removeMember(activeTeam, member)}>Delete
                      </Button>
                    </td>
                  </tr>
                {/each}
                </tbody>
              </table>
            </div>

            <div class="col">
              <h3>Add Role</h3>
              <div class="user-select">
                <div class="col-1" style="display: flex; flex: 1">
                  <RoleSelect {guildId} {roles} bind:value={selectedRole}/>
                </div>

                <div style="margin-left: 10px">
                  <Button type="button" icon="fas fa-plus" disabled={selectedRole === null || selectedRole == undefined}
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
    import WrappedSelect from "../components/WrappedSelect.svelte";
    import Dropdown from "../components/form/Dropdown.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    const USER_TYPE = 0;
    const ROLE_TYPE = 1;

    let defaultTeam = {id: 'default', name: 'Default'};

    let createName;
    let teams = [defaultTeam];
    let roles = [];
    let activeTeam = defaultTeam.id;
    let members = [];

    let selectedUser;
    let selectedRole;

    function getTeam(id) {
        return teams.find((team) => team.id === id);
    }

    async function updateActiveTeam(teamId) {
        const res = await axios.get(`${API_URL}/api/${guildId}/team/${teamId}`);
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

    async function addRole() {
        const res = await axios.put(`${API_URL}/api/${guildId}/team/${activeTeam}/${selectedRole.id}?type=1`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${selectedRole.name} has been added to the support team ${getTeam(activeTeam).name}`);

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

        members = members.filter((member) => member.id !== entity.id);

        if (entity.type === USER_TYPE) {
            notifySuccess(`${entity.name} has been removed from the team`);
        } else {
            const role = roles.find((role) => role.id === entity.id);
            notifySuccess(`${role === undefined ? "Unknown role" : role.name} has been removed from the team`);
        }
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

        activeTeam = defaultTeam.id;
        teams = teams.filter((team) => team.id !== id);

        await updateActiveTeam(defaultTeam.id);
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

        await Promise.all([
            loadTeams(),
            loadRoles()
        ]);

        await updateActiveTeam(defaultTeam.id); // Depends on teams
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
