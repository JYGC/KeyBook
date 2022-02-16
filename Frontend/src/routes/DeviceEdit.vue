<template>
    <div v-if="device">
        <device-details v-bind:device="device" v-bind:disableType=true></device-details>
        <div>
            <label for="person">Person holding device:</label>
        </div>
        <div>
            <!--<input type="text" name="person" id="person" v-model="device.personDevice.person" />-->
            <v-select name="person" id="person" :options="personUsers" label="name">
            </v-select>
        </div>
    </div>
    <div v-else>
        Can't find device details.
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';

    import { deviceView, personAllForUser } from '../api/device-api';
    
    export default Vue.extend({
        name: 'device-edit',
        props: ['deviceId'],
        data() {
            return {
                device: null,
                personUsers: null,
            };
        },
        created() {
            this.fetchDeviceDetails();
            this.fetchAllPersonForUser();
        },
        watch: {
            // call again the method if the route changes
            '$route': 'fetchDevice'
        },
        methods: {
            fetchDeviceDetails(): void {
                this.device = null;
                deviceView(this.deviceId).then(data => {
                    this.device = data;
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            },
            fetchAllPersonForUser(): void {
                this.personUsers = null;
                personAllForUser().then(data => {
                    this.personUsers = data;
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            }
        },
    });
</script>