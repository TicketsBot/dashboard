<ManageLayout>
  {#if !$loadingScreen}
    <div class="content" transition:fade>
      <div class="main-col">
        <Card footer="{false}" fill="{false}">
          <span slot="title">
            Settings
          </span>

          <div slot="body">
            real
          </div>
        </Card>
      </div>
    </div>
  {/if}
</ManageLayout>

<script>
    import {fade} from 'svelte/transition';
    import ManageLayout from "../layouts/ManageLayout.svelte";
    import Card from '../components/Card.svelte'
    import {withLoadingScreen} from '../js/util'
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {loadingScreen} from "../js/stores.js"

    export let currentRoute;
    let guildId = currentRoute.namedParams.id

    setDefaultHeaders();

    withLoadingScreen(async () => {

    })
</script>

<style>
    .content {
        display: flex;
        height: 100%;
        width: 100%;
        padding: 0 45px;
        justify-content: space-between;
    }

    .main-col {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 66.6%;
        margin-top: 30px;
    }
</style>
