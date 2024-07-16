<script lang="ts">
	import { BackendClient } from "$lib/api/backendclient.svelte";
  import { LoginApi } from "$lib/api/loginapi.svelte";
  import { RegisterApi } from "$lib/api/registerapi.svelte";

  const backendClient = new BackendClient();
  const registerApi = new RegisterApi(backendClient);
  const loginApi = new LoginApi(backendClient);

  const registerAndLogin = async () => {
    const isRegistrationSuccessful = await registerApi.callApi();
    if (!isRegistrationSuccessful) {
      return;
    }
    loginApi.email = registerApi.email;
    loginApi.password = registerApi.password;
    document.cookie = await loginApi.callApi();
    window.location.replace("");

  };
</script>

<h1>Create an Account</h1>
<div>
  <label for="name">Name</label>
  <input
    type="text"
    name="name"
    bind:value={registerApi.name}
  />
</div>
<div>
  <label for="email">Email</label>
  <input
    type="email"
    name="email"
    bind:value={registerApi.email}
  />
</div>
<div>
  <label for="password">Password</label>
  <input
    type="password"
    name="password"
    bind:value={registerApi.password}
  />
</div>

{#if registerApi.error}
  <p class="error">{registerApi.error}</p>
{/if}

<div>
  <button
    disabled={
      registerApi.name.length === 0 ||
      registerApi.email.length === 0 ||
      registerApi.password.length < 8
    }
    onclick={registerAndLogin}
  >Create account!</button>
</div>


<div>
  <p>Already have an account?</p>
  <a href="/auth/login">Login</a>
</div>