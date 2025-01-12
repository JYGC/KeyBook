import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { BackendClient } from "$lib/api/backend-client";

export const load: LayoutLoad = async () => {
  const authManager = new BackendClient();
  if (authManager.isTokenValid) {
    return redirect(303, "/");
  }
};
