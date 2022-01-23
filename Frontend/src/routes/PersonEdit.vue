<template>
    <div v-if="person">
        <person-details v-bind:person="person" v-bind:disableType="true"></person-details>
    </div>
    <div v-else>
        Can't find person details.
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        name: 'person-edit',
        props: ['personId'],
        data() {
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
                fetch('person/view/id/' + this.personId).then(r => r.json()).then(data => {
                    this.person = data;
                });
            }
        },
    });
</script>