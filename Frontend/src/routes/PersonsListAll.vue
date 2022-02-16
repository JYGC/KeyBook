<template>
    <div>
        <div>
            <router-link to="/person-add">Add Person</router-link>
        </div>
        <div>
            <table>
                <tr>
                    <th>Name</th>
                    <th>IsGone</th>
                    <th>Person Type</th>
                    <th>Current Devices</th>
                    <th></th>
                </tr>
                <tr v-for="person in persons" v-bind:key="person.id">
                    <td>{{ person.name }}</td>
                    <td>{{ person.isGone }}</td>
                    <td>{{ person.type }}</td>
                    <td>{{ person.devices }}</td>
                    <td>
                        <router-link :to="{ name: 'person-edit', params: { personId: person.id } }">Details</router-link>
                    </td>
                </tr>
            </table>
        </div>
    </div>
</template>

<script lang="ts">
    import { personAllForUser } from './../api/person-api';
    import Vue from 'vue';

    export default Vue.extend({
        name: 'persons-list-all',
        data(): {
            persons: Array<any>
        } {
            return {
                persons: [],
            };
        },
        created() {
            this.fetchAllPersonsForUser();
        },
        watch: {
            '$route': 'fetchAllPersonsForUser'
        },
        methods: {
            fetchAllPersonsForUser(): void {
                this.persons = [];
                personAllForUser().then(d => {
                    this.persons = d;
                }).catch(e => {
                    alert('error:' + e); // Add error handling later
                });
            }
        }
    });
</script>