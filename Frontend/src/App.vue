<template>
    <div id="app">
        <component :is="layout">
            <router-view />
        </component>
    </div>
</template>

<script>
    import Vue from 'vue';

    const defaultLayout = 'default';
    export default Vue.extend({
        computed: {
            // If layout not specified in router, assume defualt layout
            layout() {
                return `${defaultLayout}-layout`;
            },
        },
    });
</script>
