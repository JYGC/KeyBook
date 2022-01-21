<template>
    <div>
        <div id="nav">
            <router-link to="/login">Login</router-link> |
            <router-link to="/register">Register</router-link>
        </div>
        <slot />
    </div>
</template>