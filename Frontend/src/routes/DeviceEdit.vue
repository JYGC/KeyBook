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
                fetch('device/view/id/' + this.deviceId).then(r => r.json()).then(data => {
                    this.device = data;
                });
            },
            fetchAllPersonForUser(): void {
                this.personUsers = null;
                fetch('device/allforuser').then(r => r.json()).then(data => {
                    this.personUsers = data;
                });
            }
        },
    });
</script>