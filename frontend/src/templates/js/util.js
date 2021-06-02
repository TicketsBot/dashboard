import * as Stores from './stores'

export async function withLoadingScreen(func) {
    Stores.loadingScreen.set(true);
    await func();
    Stores.loadingScreen.set(false);
}

export function errorPage(status, message) {
    window.location.href = `/error?status=${encodeURIComponent(status)}&message=${encodeURIComponent(message)}`
}

export function notify(title, message) {
    Stores.notifyTitle.set(title);
    Stores.notifyMessage.set(message);
    Stores.notifyModal.set(true);
}

export function notifyError(message) {
    notify('Error', message);
}

export function notifySuccess(message) {
    notify('Success', message);
}

export function notifyRatelimit() {
    notifyError("You're doing that too fast: please wait a few seconds and try again");
}

export function closeNotificationModal() {
    Stores.notifyModal.set(false);
}