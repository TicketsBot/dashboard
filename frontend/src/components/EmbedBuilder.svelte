<div class="modal" transition:fade>
  <div class="modal-wrapper">
    <Card footer="{true}" footerRight="{true}" fill="{false}">
      <span slot="title">Embed Builder</span>

      <div slot="body" class="body-wrapper">
        <form class="form-wrapper" on:submit|preventDefault>
          <div class="row">
            <Colour col3 label="Embed Colour" bind:value={data.colour}/>
            <Input col3 label="Title" placeholder="Embed Title" bind:value={data.title}/>
            <Input col3 label="Title URL (Optional)" placeholder="https://example.com" bind:value={data.url}/>
          </div>

          <div class="row">
            <Textarea col1 label="Description" placeholder="Large text area, up to 4096 characters"
                      bind:value={data.description}/>
          </div>

          <Collapsible>
            <span slot="header">Author</span>

            <div slot="content" class="row">
              <Input col3 label="Author Name" placeholder="Author Name" bind:value={data.author.name}/>
              <Input col3 label="Author Icon URL (Optional)" placeholder="https://example.com/image.png"
                     tooltipText="Small icon displayed in the top left" bind:value={data.author.icon_url}/>
              <Input col3 label="Author URL (Optional)" placeholder="https://example.com"
                     tooltipText="Hyperlink on the author's name" bind:value={data.author.url}/>
            </div>
          </Collapsible>

          <Collapsible>
            <span slot="header">Images</span>
            <div slot="content" class="row">
              <Input col2 label="Large Image URL" placeholder="https://example.com/image.png"
                     bind:value={data.image_url}/>
              <Input col2 label="Small Image URL" placeholder="https://example.com/image.png"
                     bind:value={data.thumbnail_url}/>
            </div>
          </Collapsible>

          <Collapsible>
            <span slot="header">Footer</span>
            <div slot="content" class="row">
              <Input col3 label="Footer Text" placeholder="Footer Text" badge="Premium" bind:value={data.footer.text}/>
              <Input col3 label="Footer Icon URL (Optional)" badge="Premium" placeholder="https://example.com/image.png"
                     bind:value={data.footer.icon_url}/>
              <DateTimePicker col3 label="Footer Timestamp (Optional)" bind:value={data.timestamp}/>
            </div>
          </Collapsible>

          <Collapsible>
            <span slot="header">Fields</span>
            <div slot="content" class="col-1">
              {#each data.fields as field, i}
                <div class="row" style="justify-content: flex-start; gap: 10px">
                  <Input col2 label="Field Name" placeholder="Field Name" bind:value={field.name}/>
                  <Checkbox label="Inline" bind:value={field.inline}/>

                  <div style="margin-top: 18px; display: flex; align-self: center">
                    <Button danger icon="fas fa-trash-can" on:click={() => deleteField(i)}>Delete</Button>
                  </div>
                </div>
                <div class="row">
                      <Textarea col1 label="Field Value" placeholder="Large text area, up to 1024 characters"
                                bind:value={field.value}/>
                </div>
              {/each}
              <Button type="button" icon="fas fa-plus" fullWidth on:click={addField}>Add Field</Button>
            </div>
          </Collapsible>
        </form>
      </div>

      <div slot="footer">
        <Button danger={true} on:click={dispatchClose}>Cancel</Button>
        <div style="margin-left: 12px">
          <Button icon="fas fa-paper-plane" on:click={dispatchConfirm}>Submit</Button>
        </div>
      </div>
    </Card>
  </div>
</div>

<div class="modal-backdrop" transition:fade>
</div>

<svelte:window on:keydown={handleKeydown}/>

<script>
    import {createEventDispatcher} from 'svelte';
    import {fade} from 'svelte/transition'
    import Card from "./Card.svelte";
    import Button from "./Button.svelte";
    import Input from "./form/Input.svelte";
    import Colour from "./form/Colour.svelte";
    import Textarea from "./form/Textarea.svelte";
    import DateTimePicker from "./form/DateTimePicker.svelte";
    import Collapsible from "./Collapsible.svelte";
    import Checkbox from "./form/Checkbox.svelte";

    export let guildId;

    export let data;

    if (data === undefined || data === null) {
        if (!data) {
            data = {};
        }

        data.fields = [];
        data.colour = '#2ECC71';
        data.author = {};
        data.footer = {};
    }

    function addField() {
        data.fields.push({name: '', value: '', inline: false});
        data = data;
    }

    function deleteField(i) {
        data.fields.splice(i, 1);
        data = data;
    }

    const dispatch = createEventDispatcher();

    function dispatchClose() {
        dispatch('close', {});
    }

    // Dispatch with data
    function dispatchConfirm() {
        // Map blank strings to null
        const mapper = (obj) => {
            Object.keys(obj).forEach(key => {
                if (typeof obj[key] === 'string' && obj[key] === '') {
                    obj[key] = null;
                } else if (typeof obj[key] === 'object' && obj[key] !== null && obj[key] !== undefined) {
                    mapper(obj[key]);
                }
            });
        }

        mapper(data);

        dispatch('confirm', data);
    }

    function handleKeydown(e) {
        if (e.key === "Escape") {
            dispatchClose();
        }
    }
</script>

<style>
    .modal {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 999;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-wrapper {
        display: flex;
        width: 60%;
        margin: 10% auto auto auto;
    }

    @media only screen and (max-width: 1280px) {
        .modal-wrapper {
            width: 96%;
        }
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 500;
        background-color: #000;
        opacity: .5;
    }

    .form-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
    }
</style>