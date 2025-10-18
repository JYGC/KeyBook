import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { getBackendClient } from "$lib/api/backend-client";

export const load: LayoutLoad = async () => {
  const backendClient = getBackendClient();
  if (backendClient.authStore.isValid) {
    return redirect(303, "/user");
  }
};
