<script lang="ts">
  import { Link, TextInput, PasswordInput, Button, FluidForm } from "carbon-components-svelte";

	import { getBackendClient } from "$lib/api/backend-client";
  import { LoginModule } from "$lib/modules/user/login-module.svelte";
  import { RegisterModule } from "$lib/modules/user/register-module.svelte";
	import { goto } from "$app/navigation";
	import RegisterForm from "$lib/components/user/RegisterForm.svelte";

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

<RegisterForm {registerAndLogin} {registerModule} />
<br />
<br />
<br />
{#if registerModule.error}
  <p class="error">{registerModule.error}</p>
{/if}
<p>Already have an account?</p>
<Link href="/auth/login">Login</Link>