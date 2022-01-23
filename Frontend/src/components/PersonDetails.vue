<template>
    <div>
        <div>
            <label for="name">Person Name:</label>
        </div>
        <div>
            <input type="text" name="name" id="name" v-model="person.name" />
        </div>
        <div v-if="!hideIsGone">
            <div>
                <label for="isgone">Person Is Gone:</label>
            </div>
            <div>
                <input type="checkbox" name="isgone" id="name" v-model="person.isGone" />
            </div>
        </div>
        <div>
            <label for="type">Person Type:</label>
        </div>
        <div>
            <input type="text" name="type" id="type" v-model="person.type" :disabled="disableType" />
        </div>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        props: {
            'person': Object,
            'hideIsGone': {
                type: Boolean,
                default: false,
            },
            'disableType': {
                type: Boolean,
                default: false,
            },
        },
    });
</script>