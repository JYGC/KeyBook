<script lang="ts">
	import { Button, Modal } from "carbon-components-svelte";

  let {
    submitAction,
    buttonText = "",
    modalHeading = "",
    bodyMessage = "Are you sure?",
    primaryButtonText = "Confirm",
    secondaryButtonText = "Cancel",
  } = $props<{
    submitAction: () => null,
    buttonText?: string,
    modalHeading?: string,
    bodyMessage?: string,
    primaryButtonText?: string,
    secondaryButtonText?: string,
  }>();
  
  let open = $state(false);

  let onSubmit = () => {
    submitAction();
    open = false;
  };
</script>
<Button onclick={() => (open = true)}>{buttonText}</Button>
<Modal
  bind:open
  {modalHeading}
  {primaryButtonText}
  {secondaryButtonText}
  on:click:button--secondary={() => (open = false)}
  on:open
  on:close
  on:submit={onSubmit}
>
  {bodyMessage}
</Modal>