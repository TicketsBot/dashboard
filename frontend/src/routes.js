import IndexLayout from './templates/layouts/IndexLayout.svelte'
import ManageLayout from './templates/layouts/ManageLayout.svelte'
import ErrorLayout from './templates/layouts/ErrorPage.svelte'

import Index from './templates/views/Index.svelte'
import LoginCallback from './templates/views/LoginCallback.svelte'
import Login from './templates/views/Login.svelte'
import Logout from './templates/views/Logout.svelte'
import Whitelabel from './templates/views/Whitelabel.svelte'
import Settings from './templates/views/Settings.svelte'
import Error from './templates/views/Error.svelte'
import Error404 from './templates/views/Error404.svelte'
import Transcripts from './templates/views/Transcripts.svelte'

export const routes = [
    {name: '/', component: Index, layout: IndexLayout},
    {name: '404', path: '404', component: Error404, layout: ErrorLayout},
    {name: '/callback', component: LoginCallback},
    {name: '/login', component: Login},
    {name: '/logout', component: Logout},
    {name: '/error', component: Error, layout: ErrorLayout},
    {name: '/whitelabel', component: Whitelabel, layout: IndexLayout},
    {
        name: 'manage/:id',
        nestedRoutes: [
            {name: 'index', component: Error404, layout: ErrorLayout},
            {name: 'settings', component: Settings, layout: ManageLayout},
            {name: 'transcripts', component: Transcripts, layout: ManageLayout},
        ],
    }
]