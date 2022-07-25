<div class="content">
  <div class="main-col">
    <div class="card">
      <SettingsCard {guildId}/>
    </div>
  </div>
</div>

<script>
    import SettingsCard from "../components/manage/SettingsCard.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id
</script>

<style>
    .content {
        display: flex;
        flex-direction: row;
        height: 100%;
        width: 100%;
        padding: 0 45px;
        justify-content: space-between;
    }

    .main-col {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 100%;
        margin-top: 30px;
    }

    .right-col {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 32%;
        margin-top: 30px;
    }

    .card {
        margin-bottom: 2%;
    }

    @media only screen and (max-width: 1100px) {
        .content {
            flex-direction: column;
        }

        .main-col, .right-col {
            width: 100%;
        }

        .right-col {
            margin-top: 0;
        }

        .card {
            margin-bottom: 4%;
        }
    }

    @media only screen and (max-width: 576px) {
        .content {
            padding: 0 5px;
        }
    }
</style>
