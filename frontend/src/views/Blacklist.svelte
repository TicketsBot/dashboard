<div class="parent">
  <div class="content">
    <div class="main-col">
      <Card footer={false}>
        <span slot="title">Blacklisted Users</span>
        <div slot="body" class="body-wrapper">
          <table class="nice">
            <thead>
            <tr>
              <th>Username</th>
              <th>User ID</th>
              <th>Remove</th>
            </tr>
            </thead>
            <tbody>
            {#each blacklistedUsers as user}
              <tr>
                {#if user.username !== '' && user.discriminator !== ''}
                  <td>{user.username}#{user.discriminator}</td>
                {:else}
                  <td>Unknown</td>
                {/if}

                <td>{user.id}</td>
                <td>
                  <Button type="button" on:click={() => removeBlacklist(user)}>Remove</Button>
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
        <span slot="title">Blacklist A User</span>
        <div slot="body" class="body-wrapper">
          <form class="body-wrapper" on:submit|preventDefault={addBlacklist}>
            <div class="row" style="flex-direction: column">
              <UserSelect {guildId} label="User" bind:value={addUser}/>
            </div>
            <div class="row" style="justify-content: center">
              <div class="col-2">
                <Button fullWidth={true} icon="fas fa-plus"
                        disabled={addUser === undefined || addUser === ''}>Blacklist</Button>
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
    import UserSelect from "../components/form/UserSelect.svelte";
    import {notifyError, notifySuccess, withLoadingScreen} from '../js/util'
    import Button from "../components/Button.svelte";
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let addUser;
    let blacklistedUsers = [];

    async function addBlacklist() {
        const res = await axios.post(`${API_URL}/api/${guildId}/blacklist/${addUser.id}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${addUser.username}#${addUser.discriminator} has been blacklisted`);
        blacklistedUsers = [...blacklistedUsers, {
            id: addUser.id,
            username: addUser.username,
            discriminator: addUser.discriminator,
        }];
    }

    async function removeBlacklist(user) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/blacklist/${user.id}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess(`${user.username}#${user.discriminator} has been removed from the blacklist`);
        blacklistedUsers = blacklistedUsers.filter((u) => u.id !== user.id);
    }

    async function loadUsers() {
        const res = await axios.get(`${API_URL}/api/${guildId}/blacklist`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        blacklistedUsers = res.data;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadUsers();
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
