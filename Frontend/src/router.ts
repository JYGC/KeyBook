import Vue from 'vue';
import Router from 'vue-router';

import DevicesListAll from '@/routes/DevicesListAll.vue';
import DeviceEdit from '@/routes/DeviceEdit.vue';
import DeviceAdd from '@/routes/DeviceAdd.vue';

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
            path: '/device-edit/:deviceid',
            name: 'device-edit',
            component: DeviceEdit,
            props: true,
        },
        {
            path: '/device-add',
            name: 'device-add',
            component: DeviceAdd,
        },
        {
            path: '*',
            redirect: '/',
        },
    ],
});

export default router;
