<div class="super">
  <RoleSelect {guildId} placeholder="Add another role..."
              roles={roles.filter((r) => !acl.find((s) => s.role_id === r.id))} disabled={acl.length >= maxAclSize}
              on:change={(e) => addToACL(e.detail)} bind:value={roleSelectorValue}/>

  <div class="container">
    {#each acl as subject, i}
      {@const role = roles.find(r => r.id === subject.role_id)}
      <div class="subject">
        <div class="inner-left">
          <div class="row" style="gap: 10px">
            <div class="arrow-container">
              <i class="fa-solid fa-arrow-up position-arrow" class:disabled={i<=0} on:click={() => moveUp(i)}></i>
              <i class="fa-solid fa-arrow-down position-arrow" class:disabled={i>=acl.length-1}
                 on:click={() => moveDown(i)}></i>
            </div>
            <span>{role ? role.name : 'Deleted Role'}</span>
          </div>
          {#key rerender}
            <Toggle on="Allow" off="Deny"
                    hideLabel
                    toggledColor="#66bb6a"
                    untoggledColor="#e84141"
                    toggled={subject.action === "allow"}
                    on:toggle={(e) => handleToggle(subject.role_id, e.detail)}/>
          {/key}
        </div>
        <div class="inner-right">
          {#if subject.role_id !== guildId}
            <div class="delete-button">
              <i class="fas fa-x" on:click={() => removeFromACL(subject)}></i>
            </div>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<script>
    import Toggle from "svelte-toggle";
    import RoleSelect from "../form/RoleSelect.svelte";

    export let guildId;
    export let roles;
    export let maxAclSize = 10;
    export let acl = [
        {
            role_id: guildId,
            action: "allow"
        }
    ];

    let rerender = 0;
    let roleSelectorValue;

    function handleToggle(roleId, enabled) {
        const subject = acl.find(s => s.role_id === roleId);
        subject.action = enabled ? "allow" : "deny";
    }

    function addToACL(role) {
        acl = [
            {
                role_id: role.id,
                action: "allow"
            },
            ...acl
        ];

        roleSelectorValue = null;

        rerender++;
    }

    function removeFromACL(subject) {
        if (subject.role_id === guildId) return;
        acl = acl.filter(s => s.role_id !== subject.role_id);

        rerender++;
    }

    function moveUp(index) {
        if (index <= 0) {
            return;
        }

        const tmp = acl[index];
        acl[index] = acl[index - 1];
        acl[index - 1] = tmp;

        rerender++;
    }

    function moveDown(index) {
        if (index >= acl.length - 1) {
            return;
        }

        const tmp = acl[index];
        acl[index] = acl[index + 1];
        acl[index + 1] = tmp

        rerender++;
    }
</script>

<style>
    .super {
        display: flex;
        flex-direction: column;
        width: 100%;
        gap: 10px;
    }

    .container {
        display: flex;
        flex-direction: column;

        gap: 5px;
        padding: 10px;
        border-radius: 5px;
        background-color: #121212;
    }

    .subject {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;

        padding: 5px 20px;
        border-radius: 5px;
        background-color: #2e3136;
    }

    .subject > .inner-left {
        display: flex;
        flex-direction: row;
        flex: 1;
        gap: 20%;
    }

    .subject > .inner-right {
        display: flex;
        flex-direction: row;
    }

    .delete-button {
        cursor: pointer;
        transform: scale(1.33333333333, 1);
    }

    .row {
        display: flex;
        flex-direction: row;
    }

    .arrow-container {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 5px;
    }

    .position-arrow {
        cursor: pointer;
    }

    .position-arrow.disabled {
        cursor: unset !important;
        color: #777;
    }
</style>