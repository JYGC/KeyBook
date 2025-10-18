<script lang="ts">
  import { Link, TextInput, PasswordInput, Button, FluidForm } from "carbon-components-svelte";

	import { getBackendClient } from "$lib/api/backend-client";
  import { LoginModule } from "$lib/modules/user/login-module.svelte";
  import { RegisterModule } from "$lib/modules/user/register-module.svelte";
	import { goto } from "$app/navigation";

  const backendClient = getBackendClient();
  const registerModule = new RegisterModule(backendClient);
  const loginModule = new LoginModule(backendClient);

  const registerAndLogin = async () => {
    const isRegistrationSuccessful = await registerModule.callApi();
    if (!isRegistrationSuccessful) {
      return;
    }
    loginModule.email = registerModule.email;
    loginModule.password = registerModule.password;
    document.cookie = await loginModule.callApi();
    goto("/user");

  };
</script>

<h1>Create an Account</h1>
<br />
<FluidForm>
  <TextInput labelText="Name" placeholder="Enter your name..." bind:value={registerModule.name} />
  <TextInput labelText="Email" placeholder="Enter email..." bind:value={registerModule.email} />
  <PasswordInput labelText="Password" placeholder="Enter password..." bind:value={registerModule.password} />
  <PasswordInput labelText="Confirm password" placeholder="Confirm password..." bind:value={registerModule.passwordConfirm} />
</FluidForm>
<br />
<Button
  disabled={
    registerModule.name.length === 0 ||
    registerModule.email.length === 0 ||
    registerModule.password.length < 8 ||
    registerModule.passwordConfirm !== registerModule.password
  }
  onclick={registerAndLogin}
>Create account!</Button>
<br />
<br />
<br />

{#if registerModule.error}
  <p class="error">{registerModule.error}</p>
{/if}

<div>
  <p>Already have an account?</p>
  <Link href="/auth/login">Login</Link>
</div>