<div class="wrapper">
  <div class="card-wrapper">
    <Card footer=true footerRight=true>
      <span slot="title">
        Error
      </span>

      <div slot="body">
        <Route {currentRoute} {params}/>
      </div>

      <span slot="footer">
        <span class="buttons">
          <span style="margin-right: 20px">
            <Button icon="fas fa-arrow-left" on:click={back}>Back</Button>
          </span>
          <Button icon="fas fa-home" on:click={home}>Home</Button>
        </span>
      </span>
    </Card>
  </div>
</div>

<style>
    :global(body) {
        padding: 0 !important;
    }

    .wrapper {
        background-color: #121212;

        margin: 0 !important;
        padding: 0 !important;
        display: flex;
        width: 100%;
        height: 100%;

        justify-content: center;
        align-items: center;
    }

    .card-wrapper {
        display: flex;
        width: 50%;
        height: 30%;
    }

    .buttons {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
    }
</style>

<script>
    import Head from '../includes/Head.svelte'
    import Card from '../components/Card.svelte'
    import Button from "../components/Button.svelte";

    import {Route} from 'svelte-router-spa'

    export let currentRoute;
    export let params = {};

    function back() {
        window.history.back()
    }

    function home() {
        window.location.href = '/'
    }
</script>