<template>
    <div>
        <div id="nav">
            <router-link to="/">List All Devices</router-link> |
            <router-link to="/TenantListAll">List All Tenants</router-link> |
            <div class="btn" v-on:click="logout()">Logout</div>
        </div>
        <slot />
    </div>
</template>

<script>
    //import userAPI from '../_api/user';
    export default {
      methods: {
        logout() {
          //const vm = this;
          //userAPI.logoutAPI(() => {
          //  vm.$router.push('/login');
          //}, () => {
          //  console.log('Logout failure');
          //});
              alert("logout");
        },
      },
    };
</script>