import Vue from 'vue';
import Router from 'vue-router';

import DevicesListAll from '@/views/DevicesListAll.vue';
import DeviceEdit from '@/views/DeviceEdit.vue';

Vue.use(Router);

const router = new Router({
    mode: 'history',
    routes: [
        {
            path: '/',
            name: 'DevicesListAll',
            component: DevicesListAll,
        },
        {
            path: '/DeviceEdit/:id',
            name: 'DeviceEdit',
            component: DeviceEdit,
            props: true,
        },
        {
            path: '*',
            redirect: '/',
        },
    ],
});

export default router;
