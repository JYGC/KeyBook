<script lang="ts">
  import { Link } from "carbon-components-svelte";

	import { BackendClient } from "$lib/api/backend-client.svelte";
	import { LoginApi } from "$lib/api/login-api.svelte";

  const backendClient = new BackendClient();
  const loginApi = new LoginApi(backendClient);
  const login = async () => {
    document.cookie = await loginApi.callApi();
    window.location.replace("");
  };
</script>

<h1>Log into account</h1>

<div>
  <label for="email">Email</label>
  <input
    type="email"
    name="email"
    bind:value={loginApi.email}
  />
</div>
<div>
  <label for="password">Password</label>
  <input
    type="password"
    name="password"
    bind:value={loginApi.password}
  />
</div>
<div>
  <button
    disabled={loginApi.email.length === 0 || loginApi.password.length < 8}
    onclick={login}
  >Log in!</button>
</div>

<div>
  <p>Don't have an account?</p>
  <Link href="/auth/register">Create an Account</Link>
</div>