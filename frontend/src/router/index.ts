import {createRouter, createWebHashHistory} from 'vue-router'
import Remotify from "../views/Remotify.vue";

const routes = [
    {
        path: '/remotify',
        component: () => Remotify
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router