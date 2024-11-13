<div class="wrapper">
  <div class="content">
    <Card footer="{false}" fill="{false}">
      <h4 slot="title">Bot Staff</h4>
      <div slot="body" class="full-width body-wrapper">
        <form class="form-wrapper" on:submit|preventDefault={addStaff}>
          <Input label="User ID" placeholder="585576154958921739" bind:value={tempUserId} />
          <Button type="submit">Add</Button>
        </form>

        <table class="nice">
          <thead>
          <tr>
            <th>Username</th>
            <th>Remove</th>
          </tr>
          </thead>
          <tbody>
          {#each staff as user}
            <tr>
              <td>{user.username} ({user.id})</td>
              <td>
                <Button type="button" danger on:click={() => removeStaff(user.id)}>
                  Delete
                </Button>
              </td>
            </tr>
          {/each}
          </tbody>
        </table>
      </div>
    </Card>
  </div>
</div>

<script>
    import {notifyError, withLoadingScreen} from '../../js/util'
    import axios from "axios";
    import Card from '../../components/Card.svelte'
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Button from "../../components/Button.svelte";
    import Input from "../../components/form/Input.svelte";

    setDefaultHeaders()

    let staff = [];
    let tempUserId = "";

    async function addStaff() {
        if (tempUserId.length === 0) {
            return;
        }

        const res = await axios.post(`${API_URL}/api/admin/bot-staff/${tempUserId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        tempUserId = "";
        await loadData(); // TODO: Return user data with response
    }

    async function removeStaff(userId) {
        const res = await axios.delete(`${API_URL}/api/admin/bot-staff/${userId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        staff = staff.filter(user => user.id !== userId);
    }

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/admin/bot-staff`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        staff = res.data;
    }

    withLoadingScreen(async () => {
        await loadData();
    });
</script>

<style>
    .wrapper {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
        align-items: center;
    }

    .content {
        display: flex;
        justify-content: space-around;
        flex-direction: row;
        width: 95%;

        margin-top: 2%;
    }

    @media only screen and (max-width: 900px) {
        .content {
            flex-direction: column;
        }
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        gap: 30px;
    }

    .form-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
    }
</style>
