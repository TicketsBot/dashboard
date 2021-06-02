function clear(...elements) {
    for (const elementId of elements) {
        document.getElementById(elementId).value = '';
    }
}

function registerHideListener(elementId) {
    $(`#${elementId}`).on('hidden.bs.modal', hideBackdrop);
}
