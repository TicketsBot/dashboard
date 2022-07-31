<div style="margin-bottom: 8px">
  <div style="cursor: pointer;" class="inline" on:click={toggle}>
    {#if expanded}
      <i class="{retractIcon}"></i>
    {:else}
      <i class="{expandIcon}"></i>
    {/if}

    <slot name="header"></slot>

    {#if tooltip !== undefined}
      <div style="">
        <Tooltip tip={tooltip} top color="#121212">
          <i class="fas fa-circle-info form-label tooltip-icon"></i>
        </Tooltip>
      </div>
    {/if}

    <hr/>
  </div>

  <div bind:this={content} class="content">
    <slot name="content"></slot>
  </div>
</div>

<svelte:window bind:innerWidth />

<script>
    import {onMount} from "svelte";
    import Tooltip from "svelte-tooltip";

    export let retractIcon = "fas fa-minus";
    export let expandIcon = "fas fa-plus";

    export let defaultOpen = false;
    export let tooltip;

    let expanded = false;
    let showOverflow = true;

    let content;

    let innerWidth;
    $: innerWidth, updateIfExpanded();

    export function toggle() {
        if (expanded) {
            content.style.maxHeight = 0;
        } else {
            updateSize();
        }

        expanded = !expanded;
    }

    export function updateSize() {
        content.style.maxHeight = `${content.scrollHeight}px`;
    }

    function updateIfExpanded() {
        if (expanded) {
            updateSize();
        }
    }

    onMount(() => {
        content.addEventListener('DOMNodeInserted', updateIfExpanded);
        content.addEventListener('DOMNodeRemoved', updateIfExpanded);

        if (defaultOpen) toggle();
    });
</script>

<style>
    .content {
        display: flex;
        transition: max-height .3s ease-in-out, margin-top .3s ease-in-out, margin-bottom .3s ease-in-out;
        position: relative;
        overflow: hidden;
        max-height: 0;
    }

    .inline {
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 100%;
        gap: 10px;
    }

    hr {
        border-top: 1px solid #777;
        border-bottom: 0;
        border-left: 0;
        border-right: 0;
        width: 100%;
        flex: 1;
    }
</style>