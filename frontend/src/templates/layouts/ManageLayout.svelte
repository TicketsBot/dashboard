<Head/>

<div class="wrapper">
  <Navbar />
  <div class="super-container">
    <LoadingScreen/>
    <NotifyModal/>
    <div class="content-container">
      <slot/>
    </div>
  </div>
</div>

<style>
    body {
        padding: 0 !important;
    }

    .wrapper {
        margin: 0 !important;
        padding: 0 !important;
        width: 100%;
        height: 100%;
    }

    .super-container {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .content-container {
        display: flex;
        width: 100%;
        height: 100%;
    }
</style>

<script>
    import Head from '../includes/Head.svelte'
    import LoadingScreen from '../includes/LoadingScreen.svelte'
    import NotifyModal from '../includes/NotifyModal.svelte'
    import Navbar from '../includes/Navbar.svelte'
</script>