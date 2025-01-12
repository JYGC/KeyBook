<script lang="ts">
  import { Button, FluidForm, Link, PasswordInput, TextInput } from "carbon-components-svelte";

	import { BackendClient } from "$lib/api/backend-client";
	import { LoginApi } from "$lib/api/login-api.svelte";
	import { goto } from "$app/navigation";

  const backendClient = new BackendClient();
  const loginApi = new LoginApi(backendClient);
  const login = async () => {
    document.cookie = await loginApi.callApi();
    goto("/");
  };
</script>

<h1>Log into account</h1>
<br />
<FluidForm>
  <TextInput labelText="Email" placeholder="Enter email..." bind:value={loginApi.email} />
  <PasswordInput labelText="Password" placeholder="Enter password..." bind:value={loginApi.password} />
</FluidForm>
<br />
<Button
  disabled={loginApi.email.length === 0 || loginApi.password.length < 8}
  onclick={login}
>Log in!</Button>
<br />
<br />
<br />
<p>Don't have an account?</p>
<Link href="/auth/register">Create an Account</Link>