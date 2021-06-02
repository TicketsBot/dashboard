import Index from './templates/views/Index.svelte'
import LoginCallback from './templates/views/LoginCallback.svelte'
import Login from './templates/views/Login.svelte'
import Logout from './templates/views/Logout.svelte'
import Whitelabel from './templates/views/Whitelabel.svelte'
import Settings from './templates/views/Settings.svelte'
import Error from './templates/views/Error.svelte'
import Error404 from './templates/views/Error404.svelte'

export const routes = [
    {name: '/', component: Index,},
    {name: '404', path: '404', component: Error404},
    {name: '/callback', component: LoginCallback},
    {name: '/login', component: Login},
    {name: '/logout', component: Logout},
    {name: '/error', component: Error},
    {name: '/whitelabel', component: Whitelabel},
    {
        name: 'manage/:id',
        nestedRoutes: [
            {name: 'index', component: Error404},
            {name: 'settings', component: Settings},
        ],
    }
]