<li>
    <div style="width: 100%">
        <a href="{href}" class:active target="{newWindow === true ? '_blank' : '_self'}">
            {#if icon}
                <i class="fas {icon}"/>
            {/if}

            <span>{title}</span>
        </a>

        {#if active}
            <slot/>
        {/if}
    </div>
</li>

<style>
    a.active {
        background: linear-gradient(71.3deg, #873ef5 0%, #995DF3 100%);
        box-shadow: 0 6px 6px rgba(10, 10, 10, .1), 0 0 0 1px rgba(10, 10, 10, .1);
    }

    a, a:link, a:hover, a:visited, a:active {
        display: block;

        color: inherit;
        text-decoration: none;

        font-size: 16px;
        padding: 5px 10px 5px 20px;
        border-radius: 4px;
    }

    i {
        width: 20px;
        text-align: center;
    }
</style>

<script>
    export let currentRoute;
    export let title;
    export let icon;
    export let href = "#";
    export let routePrefix;
    export let newWindow;

    let active = href !== "/" && ((routePrefix || href)?.toLowerCase() === currentRoute.name.toLowerCase() ||
        currentRoute.name.toLowerCase().startsWith((routePrefix || href).toLowerCase()));

    $: active;
</script>