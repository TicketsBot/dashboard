export const API_URL = env.API_URL || "http://localhost:3000"
export const PLACEHOLDER_DOCS_URL = "https://docs.ticketsbot.net/setup/placeholders.html"

export const OAUTH = {
    clientId: env.CLIENT_ID || "700742994386747404",
    redirectUri: env.REDIRECT_URI || "http://localhost:5000/callback"
}

export const ENABLE_EXPORT = env.ENABLE_EXPORT || false;
