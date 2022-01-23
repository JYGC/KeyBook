<template>
    <div>
        <div>
            <router-link to="/device-add">Add device</router-link>
        </div>
        <div>
            <table>
                <tr>
                    <th>Device Name</th>
                    <th>Device Identifier</th>
                    <th>Status</th>
                    <th>Type</th>
                    <th></th>
                </tr>
                <tr v-for="device in devices" v-bind:key="device.id">
                    <td>{{ device.name }}</td>
                    <td>{{ device.identifier }}</td>
                    <td>{{ device.status }}</td>
                    <td>{{ device.type }}</td>
                    <td>
                        <router-link :to="{ name: 'device-edit', params: { deviceId: device.id } }">Details</router-link>
                    </td>
                </tr>
            </table>
        </div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    
    export default Vue.extend({
        name: 'device-list-all',
        data() {
            return {
                devices: null
            };
        },
        created() {
            this.fetchAllDevices();
        },
        watch: {
            // call again the method if the route changes
            '$route': 'fetchDevices'
        },
        methods: {
            fetchAllDevices(): void {
                this.devices = null;
                fetch('device').then(r => r.json()).then(data => {
                    this.devices = data;
                });
            }
        }
    });
</script>