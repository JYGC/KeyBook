import Vue from 'vue'

Vue.config.productionTip = false

// Declare layouts
import Default from './_layouts/Default.vue';
import NotLoggedIn from './_layouts/NotLoggedIn.vue';
Vue.component('default-layout', Default);
Vue.component('not-logged-in-layout', NotLoggedIn);

// Declare components
import DeviceDetails from './components/DeviceDetails.vue';
import PersonDetails from './components/PersonDetails.vue';
Vue.component('device-details', DeviceDetails);
Vue.component('person-details', PersonDetails);

// Declare routes
import Router from 'vue-router';
Vue.use(Router);
import DevicesListAll from '@/routes/DevicesListAll.vue';
import DeviceEdit from '@/routes/DeviceEdit.vue';
import DeviceAdd from '@/routes/DeviceAdd.vue';
import PersonsListAll from '@/routes/PersonsListAll.vue';
import PersonEdit from '@/routes/PersonEdit.vue';
import PersonAdd from '@/routes/PersonAdd.vue';
const router = new Router({
    mode: 'history',
    routes: [
        {
            path: '/',
            name: 'devices-list-all',
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
            path: '/persons-list-all',
            name: 'persons-list-all',
            component: PersonsListAll,
        },
        {
            path: '/person-edit/:personid',
            name: 'person-edit',
            component: PersonEdit,
            props: true,
        },
        {
            path: '/person-add',
            name: 'person-add',
            component: PersonAdd,
        },
        {
            path: '*',
            redirect: '/',
        },
    ],
});

// Configure Veux
import Vuex from 'vuex';
Vue.use(Vuex);
const store = new Vuex.Store({
    state: {

    },
    mutations: {

    },
    actions: {

    },
});

// Declare app
import App from './App.vue'
new Vue({
    router,
    store,
    render: h => h(App),
}).$mount('#app')
