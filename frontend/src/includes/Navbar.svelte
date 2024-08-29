<div class="navbar" class:dropdown={$dropdown}>
  <div class="wrapper" class:dropdown={$dropdown}>
    <div>
      <div class="burger-menu">
        <NavElement icon="fas fa-bars" on:click={dropdownNav}>Menu</NavElement>
      </div>
      <div class="nav-section" class:dropdown={$dropdown}>
        <!-- on:click required to close dropdown again -->

        {#if isAdmin}
          <NavElement icon="fas fa-cogs" link="/manage/{guildId}/settings" on:click={closeDropdown}>Settings
          </NavElement>
        {/if}

        <NavElement icon="fas fa-copy" link="/manage/{guildId}/transcripts" on:click={closeDropdown}>Transcripts
        </NavElement>

        {#if isAdmin}
          <NavElement icon="fas fa-mouse-pointer" link="/manage/{guildId}/panels" on:click={closeDropdown}>Ticket Panels</NavElement>
          <NavElement icon="fas fa-poll-h" link="/manage/{guildId}/forms" on:click={closeDropdown}>Forms</NavElement>
          <NavElement icon="fas fa-users" link="/manage/{guildId}/teams" on:click={closeDropdown}>Staff Teams</NavElement>
          <NavElement icon="fas fa-robot" link="/manage/{guildId}/integrations" on:click={closeDropdown}>Integrations</NavElement>
        {/if}

        <NavElement icon="fas fa-ticket-alt" link="/manage/{guildId}/tickets" on:click={closeDropdown}>Tickets</NavElement>
        <NavElement icon="fas fa-ban" link="/manage/{guildId}/blacklist" on:click={closeDropdown}>Blacklist</NavElement>
        <NavElement icon="fas fa-tags" link="/manage/{guildId}/tags" on:click={closeDropdown}>Tags</NavElement>
      </div>
    </div>
    <div>
      <div class="nav-section" class:dropdown={$dropdown}>
        <NavElement icon="fas fa-book" link="https://docs.ticketsbot.net">Documentation</NavElement>
        <NavElement icon="fas fa-server" link="/#">Servers</NavElement>
        <NavElement icon="fas fa-sign-out-alt" link="/logout">Logout</NavElement>
      </div>
    </div>
  </div>
</div>

<script>
    import NavElement from "../components/NavElement.svelte";
    import Badge from "../components/Badge.svelte";

    export let guildId;
    export let dropdown;
    export let permissionLevel;

    $: isAdmin = permissionLevel >= 2;

    function dropdownNav() {
        dropdown.update(v => !v);
    }

    function closeDropdown() {
        dropdown.set(false);
    }
</script>

<style>
    .navbar {
        display: none;
        justify-content: center;
        width: 100%;
        background-color: #272727;
    }

    .wrapper {
        display: flex;
        width: 98%;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
    }

    .nav-section {
        display: flex;
        flex-direction: row;
        gap: 15px;
        margin: 20px 0;
    }

    .burger-menu {
        display: none;
    }

    @media only screen and (max-width: 1154px) {
        .navbar {
            display: flex;
        }

        .nav-section {
            display: none;
        }

        .burger-menu {
            display: flex;
        }

        .nav-section.dropdown {
            display: flex;
            flex-direction: column;
        }

        .wrapper {
            flex-direction: column;
            align-items: flex-start;
            overflow: hidden;
        }

        .dropdown {
            transition-property: height;
        }

        .navbar {
            position: relative;
            height: 49px;
            transition: all .3s ease-in-out;
            overflow: hidden;
        }

        .navbar.dropdown, .wrapper.dropdown {
            height: 100%;
            overflow: hidden;
        }

        .wrapper.dropdown {
            position: absolute;
        }

        :global(.super-container.dropdown) {
            display: none;
        }
    }
</style>
