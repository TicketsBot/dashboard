export function toDays(value) {
    return Math.floor(value / 86400);
}

export function toHours(value) {
    return Math.floor((value % 86400) / 3600);
}

export function toMinutes(value) {
    return Math.floor((value % 86400 % 3600) / 60);
}