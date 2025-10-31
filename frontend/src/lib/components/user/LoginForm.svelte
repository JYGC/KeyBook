<script lang="ts">
	import type { ILoginModule } from "$lib/modules/interfaces";
	import { Button, FluidForm, PasswordInput, TextInput } from "carbon-components-svelte";

  let  {
    login,
    loginModule = $bindable(),
  } = $props<{
    login: () => Promise<void>,
    loginModule: ILoginModule
  }>();
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