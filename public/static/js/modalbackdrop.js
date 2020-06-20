function clear(...elements) {
    for (const elementId of elements) {
        document.getElementById(elementId).value = '';
    }
}

function hideBackdrop() {
    for (const backdrop of document.getElementsByClassName('modal-backdrop fade show')) {
        backdrop.remove();
    }
}

function registerHideListener(elementId) {
    $(`#${elementId}`).on('hidden.bs.modal', hideBackdrop);
}

function showBackdrop() {
    hideBackdrop();

    const backdrop = document.createElement('div');
    backdrop.classList.add('modal-backdrop', 'fade', 'show');
    document.getElementsByClassName('main-panel')[0].appendChild(backdrop);
}
