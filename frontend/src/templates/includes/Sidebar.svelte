<script>
    export let referralShow;
    export let referralLink;

    export let name;
    export let avatar;
</script>

<div class="csidebar" id="sidebar-gradient">
  <div class="sidebar-container" id="sidebar-nav">
    <div class="sidebar-element">
      <a href="/" class="sidebar-link">
        <i class="fas fa-server sidebar-icon"></i>
        <span class="sidebar-text">Servers</span>
      </a>
    </div>
    <div class="sidebar-element">
      <a href="/whitelabel" class="sidebar-link">
        <i class="fas fa-edit sidebar-icon"></i>
        <span class="sidebar-text">Whitelabel</span>
      </a>
    </div>

    {#if referralShow}
      <div class="sidebar-element ref">
        <a class="ref-wrapper" href="{referralLink}">
          <img src="/assets/img/custom.png" id="custom-image" alt="Ref Logo"/>
        </a>
      </div>
    {/if}
  </div>
  <div class="sidebar-container" style="margin-bottom: 2%">
    <div class="sidebar-element">
      <a href="/logout" onclick="clearLocalStorage();" class="sidebar-link">
        <i class="sidebar-icon fas fa-sign-out-alt sidebar-icon"></i>
        <span class="sidebar-text">Logout</span>
      </a>
    </div>
    <div class="sidebar-element">
      <a class="sidebar-link">
        <i id="avatar-sidebar" style="background: url('{avatar}') center center;"></i>
        {#if name !== undefined}
          <span class="sidebar-text">{name}</span>
        {/if}
      </a>
    </div>
  </div>
</div>

<style>
    .csidebar {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 16.6%;
        background-color: #272727;
        float: left;
        background-size: cover;
        overflow-x: hidden !important;
    }

    .sidebar-element {
        display: flex;
        align-items: center;
        width: 100%;
        cursor: pointer;
        padding: 5px 0 5px 0;
    }

    .sidebar-element:hover {
        background-color: #121212;
        transition: background-color 0.5s ease;
    }

    #custom-image {
        max-height: 70px;
        max-width: 90%;
    }

    :global(.sidebar-link) {
        display: flex;
        align-items: center;
        width: 100%;
        color: white !important;
        font-size: 18px;
        margin-left: 4%;
    }

    .sidebar-text {
        margin-left: 4%;
        display: flex;
        align-items: center;
    }

    #sidebar-nav {
        flex: 1;
    }

    .ref {
        display: flex;
        justify-content: center;
    }

    .ref-wrapper {
        display: flex;
        justify-content: center;
        padding: 10px 0;
        margin: 0 !important
    }

    .sidebar-bottom {
        position: absolute !important;
        width: 100%;
    }

    #avatar-sidebar {
        width: 32px;
        height: 32px;
        display: block;
        background-size: cover !important;
        border-radius: 50%;
    }
</style>