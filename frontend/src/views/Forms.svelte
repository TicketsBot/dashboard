<div class="parent">
  <div class="content">
    <Card footer={false}>
      <span slot="title">Forms</span>
      <div slot="body" class="body-wrapper">
        <div class="section">
          <h2 class="section-title">Create New Form</h2>

          <form on:submit|preventDefault={createForm}>
            <div class="row" id="creation-row">
              <Input placeholder="Form Title" col3={true} bind:value={newTitle}/>
              <div id="create-button-wrapper">
                <Button icon="fas fa-paper-plane" fullWidth={windowWidth <= 950}>Create</Button>
              </div>
            </div>
          </form>
        </div>
        <div class="section">
          <h2 class="section-title">Manage Forms</h2>

          <div class="col-1" style="flex-direction: row">
            <div class="col-4" style="margin-right: 12px">
              <div class="multiselect-super">
                <Dropdown col1={true} bind:value={activeFormId}>
                  <option value={null}>Select a form...</option>
                  {#each forms as form}
                    <option value="{form.form_id}">{form.title}</option>
                  {/each}
                </Dropdown>
              </div>
            </div>

            {#if activeFormId !== null}
              <div class="col-4">
                <Button danger={true} type="button"
                        on:click={() => deleteForm(activeFormId)}>Delete {getActiveFormTitle()}</Button>
              </div>
            {/if}
          </div>

          <div class="manage">
            {#if activeFormId !== null}
              {#each forms.find(form => form.form_id === activeFormId).inputs as input}
                <FormInputRow data={input} formId={activeFormId} withSaveButton={true} withDeleteButton={true}
                              on:save={(e) => editInput(activeFormId, input.id, e.detail)}
                              on:delete={() => deleteInput(activeFormId, input.id)}/>
              {/each}
            {/if}

            {#if activeFormId !== null}
              <FormInputRow bind:data={inputCreationData} withCreateButton={true}
                            on:create={(e) => createInput(e.detail)}/>
            {/if}
          </div>
        </div>
      </div>
    </Card>
  </div>
</div>

<svelte:window bind:innerWidth={windowWidth} />

<script>
  import Card from "../components/Card.svelte";
  import {notifyError, notifySuccess, withLoadingScreen} from '../js/util'
  import Button from "../components/Button.svelte";
  import axios from "axios";
  import {API_URL} from "../js/constants";
  import {setDefaultHeaders} from '../includes/Auth.svelte'
  import Input from "../components/form/Input.svelte";
  import Dropdown from "../components/form/Dropdown.svelte";
  import FormInputRow from "../components/manage/FormInputRow.svelte";

  export let currentRoute;
  let guildId = currentRoute.namedParams.id;

  let defaultTeam = {id: 'default', name: 'Default'};

  let newTitle;
  let forms = [];
  let activeFormId = null;
  let inputCreationData = {};

  $: windowWidth = 0;

  function getForm(formId) {
    return forms.find(form => form.form_id === formId);
  }

  async function createForm() {
    let data = {
      title: newTitle,
    };

    const res = await axios.post(`${API_URL}/api/${guildId}/forms`, data);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    notifySuccess(`Form ${newTitle} has been created`);
    newTitle = '';
    forms = [...forms, res.data];
    activeFormId = res.data.form_id;
  }

  async function deleteForm(id) {
    const res = await axios.delete(`${API_URL}/api/${guildId}/forms/${id}`);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    notifySuccess(`Form deleted successfully`);

    forms = forms.filter(form => form.form_id !== id);
    if (forms.length > 0) {
      activeFormId = forms[0].form_id;
    } else {
      activeFormId = null;
    }
  }

  async function createInput(data) {
    let mapped = {...data, style: parseInt(data.style)};

    const res = await axios.post(`${API_URL}/api/${guildId}/forms/${activeFormId}`, mapped);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    let form = getForm(res.data.form_id);
    form.inputs = [...form.inputs, res.data];
    forms = forms;
    inputCreationData = {};

    notifySuccess('Form input created successfully');
  }

  async function editInput(formId, inputId, data) {
    let mapped = {...data, style: parseInt(data.style)};

    const res = await axios.patch(`${API_URL}/api/${guildId}/forms/${formId}/${inputId}`, mapped);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    let form = getForm(formId);
    form.inputs = form.inputs.filter(input => input.id !== inputId);
    form.inputs = [...form.inputs, res.data];

    notifySuccess('Form input updated successfully');
  }

  async function deleteInput(formId, inputId) {
    const res = await axios.delete(`${API_URL}/api/${guildId}/forms/${formId}/${inputId}`);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    // TODO: delete keyword?
    let form = getForm(formId);
    form.inputs = form.inputs.filter(input => input.id !== inputId);
    forms = forms;

    notifySuccess('Form input deleted successfully');
  }

  async function loadForms() {
    const res = await axios.get(`${API_URL}/api/${guildId}/forms`);
    if (res.status !== 200) {
      notifyError(res.data.error);
      return;
    }

    forms = res.data || [];
  }

  function getActiveFormTitle() {
    return activeFormId !== null ? forms.find(f => f.form_id === activeFormId).title : 'Unknown';
  }

  withLoadingScreen(async () => {
    setDefaultHeaders();
    await loadForms();
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
        margin-bottom: 50px;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        padding: 1%;
    }

    .section {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .section:not(:first-child) {
        margin-top: 2%;
    }

    .section-title {
        font-size: 36px;
        font-weight: bolder !important;
    }

    h3 {
        font-size: 28px;
        margin-bottom: 4px;
    }

    .row {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
    }

    .manage {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-top: 12px;
    }

    #creation-row {
        justify-content: flex-start !important;
    }

    #create-button-wrapper {
        margin-left: 15px;
        height: 40px;
    }

    @media only screen and (max-width: 950px) {
        .manage {
            flex-direction: column;
        }

        .row {
            flex-direction: column;
        }

        #create-button-wrapper {
            margin-left: unset;
        }
    }
</style>
