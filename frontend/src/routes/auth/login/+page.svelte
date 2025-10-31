<script lang="ts">
  import { Button, FluidForm, Link, PasswordInput, TextInput } from "carbon-components-svelte";

	import { getBackendClient } from "$lib/api/backend-client";
	import { LoginModule } from "$lib/modules/user/login-module.svelte";
	import { goto } from "$app/navigation";
	import LoginForm from "$lib/components/user/LoginForm.svelte";

  const backendClient = getBackendClient();
  const loginModule = new LoginModule(backendClient);
  const login = async () => {
    document.cookie = await loginModule.callApi();
    goto("/user");
  };
</script>

<LoginForm {login} {loginModule} />
<br />
<br />
<br />
<p>Don't have an account?</p>
<Link href="/auth/register">Create an Account</Link>
