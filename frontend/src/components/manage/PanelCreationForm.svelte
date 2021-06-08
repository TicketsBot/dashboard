<Card footer="{false}">
  <span slot="title">Create Panel</span>

  <div slot="body" class="body-wrapper">
    <form class="settings-form" on:submit|preventDefault={updateSettings}>
      <div class="row">
        <Input label="Panel Title" placeholder="Open a ticket!" col3=true bind:value={data.prefix} />
        <Textarea label="Panel Content" placeholder="By clicking the button, a ticket will be opened for you."></Textarea>
      </div>
      <div class="row">
        <Textarea label="welcome message" placeholder="Thanks for opening a ticket!" col1=true bind:value={data.welcome_message} />
      </div>
      <div class="row">
        <ChannelDropdown label="Archive Channel" col2=true channels={channels} bind:value={data.archive_channel} />
        <CategoryDropdown label="Channel Category" col2=true channels={channels} bind:value={data.category} />
      </div>
      <div class="row">
        <NamingScheme col4=true bind:value={data.naming_scheme} />
      </div>
      <div class="row">
        <div class="col-1">
          <Button icon="fas fa-paper-plane" fullWidth=true>Submit</Button>
        </div>
      </div>
    </form>
  </div>
</Card>

<script>
    import Card from "../Card.svelte";
    import Input from "../form/Input.svelte";
    import Textarea from "../form/Textarea.svelte";

    export let data = {
        title: 'Open a ticket!',
        welcome_message: 'By clicking the button, a ticket will be opened for you.'
    };
</script>

<style>
  .body-wrapper {
      width: 100%;
  }
</style>