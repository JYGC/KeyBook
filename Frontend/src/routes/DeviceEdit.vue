<template>
    <div v-if="device">
        <device-details v-bind:device="device" v-bind:disableType=true></device-details>
        <div>
            <label for="person">Person holding device:</label>
        </div>
        <div>
            <b-dropdown id="dropdown-1" :text="(selectedPerson === null) ? 'Not Used' : selectedPerson.name" v-model="selectedPerson" class="m-md-2">
                <b-dropdown-item @click="selectedPerson = null">
                    Not Used
                </b-dropdown-item>
                <b-dropdown-item v-for="person in personUsers" :key="person.id" @click="selectPerson(person.id)">
                    {{ person.name }}
                </b-dropdown-item>
            </b-dropdown>
        </div>
        <div v-on:click="saveDevice()">Add device</div>
    </div>
    <div v-else>
        Can't find device details.
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';

    import { deviceView, deviceSave } from '../api/device-api';
    import { personAllForUser } from '../api/person-api';
    
    export default Vue.extend({
        name: 'device-edit',
        props: { deviceId: String },
        data(): {
            device: any,
            personUsers: Array<any>,
            selectedPerson: any
        } {
            return {
                device: null,
                personUsers: [],
                selectedPerson: null
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
                this.personUsers = [];
                personAllForUser().then(data => {
                    this.personUsers = data;
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            },
            selectPerson(persionId: string): void {
                this.selectedPerson = this.personUsers.find(p => p.id == persionId);
            },
            saveDevice(): void {
                deviceSave(this.device, this.selectedPerson).then(data => {
                    console.log(data)
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            }
        },
    });
</script>