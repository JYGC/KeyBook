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

    export default Vue.extend({
        name: 'person-add',
        data() {
            return {
                person: {}
            };
        },
        methods: {
            addPerson(): void {
                fetch('person/add', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.person),
                }).then(r => r.json()).then(data => {
                    this.person = data;
                    this.$router.push('/persons-list-all');
                }).catch(e => {
                    console.log(e);
                });
            },
        }
    });
</script>