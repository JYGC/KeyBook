<script lang="ts">
  import { Link, TextInput, PasswordInput, Button } from "carbon-components-svelte";

	import { BackendClient } from "$lib/api/backend-client.svelte";
  import { LoginApi } from "$lib/api/login-api.svelte";
  import { RegisterApi } from "$lib/api/register-api.svelte";
	import { goto } from "$app/navigation";

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
    goto("/");

  };
</script>

<h1>Create an Account</h1>
<br />
<TextInput placeholder="Enter your name..." bind:value={registerApi.name} />
<br />
<TextInput placeholder="Enter email..." bind:value={registerApi.email} />
<br />
<PasswordInput placeholder="Enter password..." bind:value={registerApi.password} />
<br />
<PasswordInput placeholder="Confirm password..." bind:value={registerApi.passwordConfirm} />
<br />
<Button
  disabled={
    registerApi.name.length === 0 ||
    registerApi.email.length === 0 ||
    registerApi.password.length < 8 ||
    registerApi.passwordConfirm !== registerApi.password
  }
  onclick={registerAndLogin}
>Create account!</Button>
<br />
<br />
<br />

{#if registerApi.error}
  <p class="error">{registerApi.error}</p>
{/if}

<div>
  <p>Already have an account?</p>
  <Link href="/auth/login">Login</Link>
</div>