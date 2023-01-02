{#if $notifyModal}
  <div class="modal" transition:fade="{{duration: 500}}">
    <div class="modal-wrapper" bind:this={wrapper}>
      <Card footer="{true}" footerRight="{true}" fill="{false}">
        <span slot="title">{$notifyTitle}</span>

        <span slot="body">{$notifyMessage}</span>

        <div slot="footer">
          <Button on:click={closeNotificationModal}>
            Close
          </Button>
        </div>
      </Card>
    </div>
  </div>

  <div class="modal-backdrop" transition:fade="{{duration: 500}}">
  </div>
{/if}

<script>
    import {notifyMessage, notifyModal, notifyTitle} from "../js/stores";
    import {closeNotificationModal} from "../js/util";
    import {fade} from 'svelte/transition'
    import Card from '../components/Card.svelte'
    import Button from '../components/Button.svelte'

    let wrapper;

    document.addEventListener('click', (e) => {
        if (!notifyModal) {
            return;
        }

        let current = e.target;
        let wrapperFound = false;

        while (current) {
            if (current.attributes) {
                if (current.hasAttribute('istrigger')) {
                    wrapperFound = true;
                    break;
                }
            }

            if (current === wrapper) {
                wrapperFound = true;
                break;
            } else {
                current = current.parentNode;
            }
        }

        if (!wrapperFound) {
            closeNotificationModal();
        }
    });
</script>

<style>
    .modal {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1001;

        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal-wrapper {
        display: flex;
        width: 50%;
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1000;
        background-color: #000;
        opacity: .5;
    }

    .footer {
        display: flex;
        width: 100%;
        height: 100%;
    }
</style>