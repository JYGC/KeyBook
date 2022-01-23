<template>
    <div v-if="device">
        <device-details v-bind:device="device" v-bind:disableType=true></device-details>
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
                device: null
            };
        },
        created() {
            this.fetchDeviceDetails();
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
            }
        },
    });
</script>