<template>
    <div v-if="person">
        <person-details v-bind:person="person" v-bind:disableType="true"></person-details>
        <div v-if="person.personDevices && person.personDevices.length > 0">
            <div class="persns-devices">
                <table>
                    <tr>
                        <th>Device Name</th>
                        <th>Device Identifier</th>
                        <th>Status</th>
                        <th>Type</th>
                        <th></th>
                    </tr>
                    <tr v-for="personDevice in person.personDevices" v-bind:key="personDevice.id">
                        <td>{{ personDevice.device.name }}</td>
                        <td>{{ personDevice.device.identifier }}</td>
                        <td>{{ personDevice.device.status }}</td>
                        <td>{{ personDevice.device.type }}</td>
                        <td></td>
                    </tr>
                </table>
            </div>
        </div>
        <div v-else>
            Person has no devices
        </div>
    </div>
    <div v-else>
        Can't find person details.
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import { personView } from '../api/person-api';

    export default Vue.extend({
        name: 'person-edit',
        props: { personId: String },
        data(): {
            person: any
        } {
            return {
                person: null
            };
        },
        created() {
            this.fetchPersonDetails();
        },
        watch: {
            '$route': 'fetchPersonDetails'
        },
        methods: {
            fetchPersonDetails(): void {
                this.person = null;
                personView(this.personId).then(data => {
                    this.person = data;
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                })
            }
        },
    });
</script>