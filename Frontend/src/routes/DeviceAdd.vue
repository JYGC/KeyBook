<template>
    <div>
        <div>
            <device-details v-bind:device="device" v-bind:hideStatus=true></device-details>
        </div>
        <div v-on:click="addDevice()">Add device</div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import { deviceAdd } from '../api/device-api';

    export default Vue.extend({
        name: 'device-add',
        data() {
            return {
                device: {},
            };
        },
        methods: {
            addDevice(): void {
                deviceAdd(this.device).then(d => {
                    this.device = d;
                    this.$router.push('/');
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            },
        },
    });
</script>