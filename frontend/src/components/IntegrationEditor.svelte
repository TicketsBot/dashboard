{#if deleteConfirmationOpen}
  <ConfirmationModal icon="fas fa-globe" isDangerous on:cancel={() => publicConfirmationOpen = false}
                     on:confirm={dispatchDelete}>
    <span slot="body">Are you sure you want to delete your integration {data.name}?</span>
    <span slot="confirm">Delete</span>
  </ConfirmationModal>
{/if}

{#if publicConfirmationOpen}
  <ConfirmationModal icon="fas fa-trash-icon" on:cancel={() => deleteConfirmationOpen = false}
                     on:confirm={dispatchMakePublic}>
    <p slot="body">Are you sure you want to make your integration <i>{data.name}</i> public? Everyone will be able to
      add it to their servers.</p>
    <span slot="confirm">Confirm</span>
  </ConfirmationModal>
{/if}

<div class="parent">
  <div class="content">
    {#if editingMetadata || editMode}
      <div class="row outer-row" bind:this={metadataRow} transition:fade>
        <div class="col-metadata">
          <Card footer footerRight>
            <span slot="title">Integration Metadata</span>
            <div slot="body" class="body-wrapper">
              <p>Let people know what your integration does. A preview will be generated as you type.</p>

              <div class="row">
                <Input col2 label="Name" placeholder="My Integration" bind:value={data.name}/>
              </div>
              <div class="row">
                <Input col2 label="Image URL" placeholder="https://example.com/logo.png" bind:value={data.image_url}
                       on:change={ensureNullIfBlank}/>
                <Input col2 label="Privacy Policy URL" placeholder="https://example.com/privacy"
                       bind:value={data.privacy_policy_url} on:change={ensureNullIfBlank}/>
              </div>
              <div class="row">
                <Textarea col1 label="Description" placeholder="Let people know what your integration does"
                          bind:value={data.description}/>
              </div>
            </div>
            <div slot="footer" style="gap: 12px">
              {#if editMode}
                <Button icon="fas fa-globe" disabled={data.public} on:click={confirmMakePublic}>Make Public</Button>

                <Button danger icon="fas fa-trash-can" on:click={requestDeleteConfirmation}>
                  Delete Integration
                </Button>
              {:else}
                <Button icon="fas fa-arrow-right" disabled={data.name.length === 0 || data.description.length === 0}
                        on:click={nextStage}>
                  Continue
                </Button>
              {/if}
            </div>
          </Card>
        </div>
        <div class="col-preview">
          <div class="preview">
            <Integration hideLinks name={data.name} imageUrl={data.image_url}>
              <span slot="description">{data.description}</span>
            </Integration>
          </div>
        </div>
      </div>
    {/if}

    {#if !editingMetadata || editMode}
      <div class="row outer-row" transition:fade>
        <div class="col-left">
          <Card footer footerRight>
            <span slot="title">HTTP Request</span>
            <div slot="body" class="body-wrapper">
              <div>
                <h3>API Endpoint</h3>
                <div class="section">
                  <p>When a user opens a ticket, a HTTP <code>{data.http_method}</code> request will be sent to the
                    provided
                    request
                    URL. The URL must respond with a valid JSON payload.</p>
                  <div class="row">
                    <Dropdown col4 bind:value={data.http_method} label="Request Method">
                      <option value="GET">GET</option>
                      <option value="POST">POST</option>
                    </Dropdown>
                    <div class="col-3-4">
                      <Input col1 label="Request URL" bind:value={data.webhook_url}
                             placeholder="https://api.example.com/users/find?discord=%user_id%"/>
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <h3>Secrets</h3>
                <div class="section">
                  <p>
                    If creating a public integration, you may wish to let users provide secret values, e.g. API keys,
                    instead of sending all requests through your own.
                  </p>

                  <p>Note: Do not include the <code>%</code> symbols in secret names, they will be automatically
                    included
                  </p>

                  <div class="col">
                    {#each data.secrets as secret, i}
                      <div class="col">
                        <div class="row">
                          {#if i === 0}
                            <Input col1 label="Secret Name" placeholder="api_key" bind:value={secret.name}/>
                          {:else}
                            <Input col1 placeholder="api_key" bind:value={secret.name}/>
                          {/if}

                          <div class="button-anchor">
                            <Button danger iconOnly icon="fas fa-trash-can" on:click={() => deleteSecret(i)}/>
                          </div>
                        </div>

                        <div class="row">
                          <Textarea col1 minHeight="60px" label="Description" bind:value={secret.description}
                                    placeholder="Tell users what value to enter for this secret, in up to 255 characters"/>
                        </div>

                        {#if i !== data.secrets.length - 1}
                          <hr/>
                        {/if}
                      </div>
                    {/each}
                  </div>

                  <Button fullWidth icon="fas fa-plus" on:click={addSecret} disabled={data.secrets.length >= 5}>
                    Add Additional Secret
                  </Button>
                </div>
              </div>

              <div>
                <h3>Request Headers</h3>
                <div class="section">
                  <p>You can specify up to 5 HTTP headers that will be sent with the request, for example, containing
                    authentication
                    keys. You may specify the user's ID in a header, via <code>%user_id%</code>.
                  </p>

                  <p>
                    You may also include the values of secrets you have created, via <code>%secret_name%</code>.
                    {#if data.secrets.length > 0}
                      For example, <code>%{data.secrets[0].name}%</code>.
                    {/if}
                  </p>

                  <div class="col">
                    {#each data.headers as header, i}
                      <div class="row">
                        {#if i === 0}
                          <Input col2 label="Header Name" placeholder="x-auth-key" bind:value={header.name}/>
                          <Input col2 label="Header Value" placeholder="super secret key" bind:value={header.value}/>
                          <div class="button-anchor">
                            <Button danger iconOnly icon="fas fa-trash-can" on:click={() => deleteHeader(i)}/>
                          </div>
                        {:else}
                          <Input col2 placeholder="x-auth-key" bind:value={header.name}/>
                          <Input col2 placeholder="super secret key" bind:value={header.value}/>
                          <div class="button-anchor">
                            <Button danger iconOnly icon="fas fa-trash-can" on:click={() => deleteHeader(i)}/>
                          </div>
                        {/if}
                      </div>
                    {/each}
                  </div>

                  <Button fullWidth icon="fas fa-plus" on:click={addHeader} disabled={data.headers.length >= 5}>
                    Add Additional Header
                  </Button>
                </div>
              </div>
            </div>
            <div slot="footer">
              {#if editMode}
                <Button icon="fas fa-floppy-disk" on:click={dispatchSubmit}>Save</Button>
              {:else }
                <Button icon="fas fa-floppy-disk" on:click={dispatchSubmit}>Create</Button>
              {/if}
            </div>
          </Card>
        </div>
        <div class="col-right">
          <Card footer={false} fill={false}>
            <span slot="title">Placeholders</span>
            <div slot="body" class="body-wrapper">
              <div class="section">
                <p>
                  The response must contain a valid JSON payload. This payload will be parsed, and values can be
                  extracted
                  to use as placeholders in your welcome message.
                </p>

                <p>
                  Do <b>not</b> include the % symbols in the placeholder names. They will be included automatically.
                </p>

                <p>
                  The JSON path is the key path to access a field in the response JSON. You can use a period
                  (e.g. <code>user.username</code>) to access nested objects.
                  You will be presented with an example JSON payload as you type.
                </p>

                <div class="col">
                  {#each data.placeholders as placeholder, i}
                    <div class="row">
                      {#if i === 0}
                        <Input col2 label="Placeholder" placeholder="ingame_username"
                               bind:value={placeholder.name}/>
                        <Input col2 label="JSON Path" placeholder="user.username" bind:value={placeholder.json_path}/>
                        <div class="button-anchor">
                          <Button danger iconOnly icon="fas fa-trash-can" on:click={() => deletePlaceholder(i)}/>
                        </div>
                      {:else}
                        <Input col2 placeholder="ingame_username" bind:value={placeholder.name}/>
                        <Input col2 placeholder="user.username" bind:value={placeholder.json_path}/>
                        <div class="button-anchor">
                          <Button danger iconOnly icon="fas fa-trash-can" on:click={() => deletePlaceholder(i)}/>
                        </div>
                      {/if}
                    </div>
                  {/each}
                </div>

                <Button fullWidth icon="fas fa-plus" on:click={addPlaceholder}
                        disabled={data.placeholders.length >= 15}>
                  Add Additional Placeholder
                </Button>
              </div>

              <div>
                <h3>Example Response</h3>
                <div class="section">
                  <p>The request must be responded to with a JSON payload in the following form:</p>
                  <code class="codeblock">
                    {exampleJson}
                  </code>
                </div>
              </div>
            </div>
          </Card>
        </div>
      </div>
    {/if}
  </div>
</div>

<script>
    import {fade} from 'svelte/transition';
    import Card from "./Card.svelte";
    import Button from "./Button.svelte";
    import Dropdown from "./form/Dropdown.svelte";
    import Input from "./form/Input.svelte";
    import Integration from "./manage/Integration.svelte";
    import Textarea from "./form/Textarea.svelte";
    import {createEventDispatcher, onMount} from "svelte";
    import ConfirmationModal from "./ConfirmationModal.svelte";

    const dispatch = createEventDispatcher();

    export let guildId;

    let metadataRow;
    let editingMetadata = true;

    let exampleJson = "{}";

    export let data = {
        name: "",
        description: "",
        image_url: "",
        privacy_policy_url: "",
        http_method: "GET",
        placeholders: [],
        headers: [],
        secrets: [],
    };

    export let editMode = false;
    let deleteConfirmationOpen = false;
    let publicConfirmationOpen = false;

    function requestDeleteConfirmation() {
        deleteConfirmationOpen = true;
    }

    function confirmMakePublic() {
        publicConfirmationOpen = true;
    }

    // on:input uses the old value!
    $: data.placeholders, normalisePlaceholders();
    $: data.placeholders, updateExampleJson();
    $: data.secrets, normaliseSecrets();
    $: data.headers, normaliseHeaders();
    $: data.name, data.name = data.name.substring(0, 32);
    $: data.description, data.description = data.description.substring(0, 255);

    function addPlaceholder() {
        data.placeholders.push({name: "", json_path: ""});
        data = data;
    }

    function deletePlaceholder(i) {
        data.placeholders.splice(i, 1);
        data = data;
    }

    function normalisePlaceholders() {
        data.placeholders = data.placeholders.map((placeholder) => {
            placeholder.name = placeholder.name.replaceAll(' ', '_').replaceAll('%', '');
            return placeholder;
        });
    }

    function addHeader() {
        data.headers.push({name: "", value: ""});
        data = data;
    }

    function deleteHeader(i) {
        data.headers.splice(i, 1);
        data = data;
    }

    function normaliseHeaders() {
        data.headers = data.headers.map((header) => {
            header.name = header.name.replaceAll(' ', '-');
            return header;
        });
    }

    function addSecret() {
        data.secrets.push({name: ""});
        data = data;
    }

    function deleteSecret(i) {
        data.secrets.splice(i, 1);
        data = data;
    }

    function normaliseSecrets() {
        data.secrets = data.secrets.map((secret) => {
            secret.name = secret.name.replaceAll(' ', '_').replaceAll('%', '');
            return secret;
        });
    }

    function nextStage() {
        editingMetadata = false;
        metadataRow.style.display = 'none';
    }

    function ensureNullIfBlank() {
        if (data.image_url !== undefined && data.image_url !== null && data.image_url.length === 0) {
            data.image_url = null;
        }

        if (data.privacy_policy_url !== undefined && data.privacy_policy_url !== null && data.privacy_policy_url.length === 0) {
            data.privacy_policy_url = null;
        }
    }

    function updateExampleJson() {
        try {
            let obj = {};
            for (const placeholder of data.placeholders) {
                let split = placeholder.json_path.split(".");
                let current = obj;
                for (const [index, part] of split.entries()) {
                    if (index === split.length - 1) {
                        current[part] = "...";
                    } else {
                        if (current[part]) {
                            current = current[part];
                        } else {
                            current[part] = {};
                            current = current[part];
                        }
                    }
                }
            }

            exampleJson = JSON.stringify(obj, null, 2);
        } catch (e) {
            exampleJson = "Invalid JSON";
        }
    }

    function dispatchSubmit() {
        ensureNullIfBlank();
        dispatch("submit", data);
    }

    function dispatchMakePublic() {
        publicConfirmationOpen = false;
        dispatch("makePublic", data);
    }

    function dispatchDelete() {
        dispatch("delete", data);
    }

    onMount(() => {
        updateExampleJson();

        if (!editMode) {
            addPlaceholder();
        }
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
        flex-direction: column;
        align-items: flex-start;
        row-gap: 5vh;
        width: 96%;
        height: 100%;
        margin-top: 30px;
        padding-bottom: 5vh;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        row-gap: 1vh;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: flex-start;
        gap: 10px;
        width: 100%;
    }

    .col-left, .col-right {
        width: 50%;
    }

    .col-metadata {
        display: flex;
        flex: 1;
    }

    .col-preview {
        display: flex;
        width: 33%;
        max-width: 33%;
    }

    .preview {
        width: 100%;
    }

    .button-anchor {
        align-self: flex-end;
        margin-bottom: 8px;
    }

    .codeblock {
        border-color: #2e3136;
        background-color: #2e3136;
        color: white;
        outline: none;
        border-radius: 4px;
        padding: 8px 12px;
        box-sizing: border-box;

        width: 100%;
        height: 100%;

        display: flex;
        white-space: pre-wrap;
    }

    .section {
        display: flex;
        flex-direction: column;
        row-gap: 1vh;
    }

    hr {
        border-top: 1px solid #777;
        border-bottom: 0;
        border-left: 0;
        border-right: 0;
        width: 100%;
    }

    @media only screen and (max-width: 950px) {
        .outer-row {
            flex-direction: column;
            align-items: center;
        }

        .col-left, .col-right {
            width: 100%;
        }
    }
</style>
