<script lang="ts">
  import { Button, FluidForm, Link, PasswordInput, TextInput } from "carbon-components-svelte";

	import { getBackendClient } from "$lib/api/backend-client";
	import { LoginModule } from "$lib/modules/user/login-module.svelte";
	import { goto } from "$app/navigation";

  const backendClient = getBackendClient();
  const loginModule = new LoginModule(backendClient);
  const login = async () => {
    document.cookie = await loginModule.callApi();
    goto("/user");
  };
</script>

<h1>Log into account</h1>
<br />
<FluidForm>
  <TextInput labelText="Email" placeholder="Enter email..." bind:value={loginModule.email} />
  <PasswordInput labelText="Password" placeholder="Enter password..." bind:value={loginModule.password} />
</FluidForm>
<br />
<Button
  disabled={loginModule.email.length === 0 || loginModule.password.length < 8}
  onclick={login}
>Log in!</Button>
<br />
<br />
<br />
<p>Don't have an account?</p>
<Link href="/auth/register">Create an Account</Link>
