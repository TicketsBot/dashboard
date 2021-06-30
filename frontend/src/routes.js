import IndexLayout from './layouts/IndexLayout.svelte'
import ManageLayout from './layouts/ManageLayout.svelte'
import ErrorLayout from './layouts/ErrorPage.svelte'
import TranscriptViewLayout from './layouts/TranscriptViewLayout.svelte'

import Index from './views/Index.svelte'
import LoginCallback from './views/LoginCallback.svelte'
import Login from './views/Login.svelte'
import Logout from './views/Logout.svelte'
import Whitelabel from './views/Whitelabel.svelte'
import Settings from './views/Settings.svelte'
import Error from './views/Error.svelte'
import Error404 from './views/Error404.svelte'
import Transcripts from './views/Transcripts.svelte'
import TranscriptView from './views/TranscriptView.svelte'
import Blacklist from './views/Blacklist.svelte'
import Panels from './views/Panels.svelte'
import Tags from './views/Tags.svelte'
import Teams from './views/Teams.svelte'
import Tickets from './views/Tickets.svelte'
import TicketView from './views/TicketView.svelte'

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
            {
                name: 'transcripts',
                nestedRoutes: [
                    {
                        name: 'index',
                        component: Transcripts,
                        layout: ManageLayout,
                    },
                    {
                        name: 'view/:ticketid',
                        component: TranscriptView, // just to test
                        layout: TranscriptViewLayout,
                    }
                ]
            },
            // Backwards compatibility
            {
                name: 'logs',
                nestedRoutes: [
                    {
                        name: 'view/:ticketid',
                        component: TranscriptView,
                        layout: TranscriptViewLayout,
                    }
                ]
            },
            {name: 'panels', component: Panels, layout: ManageLayout},
            {name: 'blacklist', component: Blacklist, layout: ManageLayout},
            {name: 'tags', component: Tags, layout: ManageLayout},
            {name: 'teams', component: Teams, layout: ManageLayout},
            {
                name: 'tickets',
                nestedRoutes: [
                    {
                        name: 'index',
                        component: Tickets,
                        layout: ManageLayout,
                    },
                    {
                        name: 'view/:ticketid',
                        component: TicketView,
                        layout: ManageLayout,
                    }
                ]
            },
        ],
    }
]