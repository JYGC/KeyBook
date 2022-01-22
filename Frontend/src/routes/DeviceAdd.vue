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

    export default Vue.extend({
        name: 'device-add',
        data() {
            return {
                device: {},
            };
        },
        methods: {
            addDevice(): void {
                fetch('device/add', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.device),
                }).then(r => r.json()).then(data => {
                    this.device = data;
                }).catch(e => {
                    console.log(e);
                });
            },
        },
    });
</script>