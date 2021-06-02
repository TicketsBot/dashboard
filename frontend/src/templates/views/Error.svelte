<ErrorPage>
  <span slot="title">
    <b>Error {status}</b>
  </span>

  <span slot="body">
    {message}
  </span>
</ErrorPage>

<script>
    import ErrorPage from '../layouts/ErrorPage.svelte'

    const urlParams = new URLSearchParams(window.location.search);
    const status = urlParams.get('status');
    const message = urlParams.get('message');
</script>
