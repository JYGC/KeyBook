<script lang="ts">
  import { AspectRatio, Button, Link, PasswordInput, TextInput } from "carbon-components-svelte";

	import { BackendClient } from "$lib/api/backend-client.svelte";
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
<TextInput placeholder="Enter email..." bind:value={loginApi.email} />
<br />
<PasswordInput placeholder="Enter password..." bind:value={loginApi.password} />
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