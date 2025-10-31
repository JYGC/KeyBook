<script lang="ts">
	import type { IRegisterModule } from '$lib/modules/interfaces';
	import { Button, FluidForm, PasswordInput, TextInput } from 'carbon-components-svelte';
  
  let {
    registerAndLogin,
    registerModule = $bindable(),
  } = $props<{
    registerAndLogin: () => Promise<void>,
    registerModule: IRegisterModule
  }>();
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