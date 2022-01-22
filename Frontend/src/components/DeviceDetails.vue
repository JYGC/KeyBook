<template>
    <div>
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
        <div div v-if="!hideStatus">
            <div>
                <label for="status">Device Status:</label>
            </div>
            <div>
                <input type="text" name="status" id="status" v-model="device.status" />
            </div>
        </div>
        <div>
            <label for="type">Device Type:</label>
        </div>
        <div>
            <input type="text" name="type" id="type" v-model="device.type" />
        </div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        props: {
            'device': Object,
            'hideStatus': {
                type: Boolean,
                default: false,
            },
        },
    });
</script>