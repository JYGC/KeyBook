<template>
    <div v-if="device">
        <div>
            <label for="name">Device Name:</label>
        </div>
        <div>
            <input type="text" name="name" id="name" v-model="device.name" />
        </div>
        <div>
            <label for="identifier">Device Identifier:</label>
        </div>
        <div>
            <input type="text" name="identifier" id="identifier" v-model="device.identifier" />
        </div>
        <div>
            <label for="status">Device Status:</label>
        </div>
        <div>
            <input type="text" name="status" id="status" v-model="device.status" />
        </div>
        <div>
            <label for="type">Device Type:</label>
        </div>
        <div>
            <input type="text" name="type" id="type" v-model="device.type" />
        </div>
    </div>
    <div v-else>
        fail
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    
    export default Vue.extend({
        name: 'DeviceEdit',
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
        }
    });
</script>