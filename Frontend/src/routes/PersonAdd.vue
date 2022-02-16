<template>
    <div>
        <div>
            <person-details v-bind:person="person" v-bind:hideIsGone=true></person-details>
        </div>
        <div v-on:click="addPerson()">Add person</div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';
    import { personAdd } from '../api/person-api';

    export default Vue.extend({
        name: 'person-add',
        data() {
            return {
                person: {}
            };
        },
        methods: {
            addPerson(): void {
                personAdd(this.person).then(data => {
                    this.person = data;
                    this.$router.push('/persons-list-all');
                })
            },
        }
    });
</script>