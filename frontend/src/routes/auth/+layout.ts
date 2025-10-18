import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { getBackendClient } from "$lib/api/backend-client";
import { AuthService } from "$lib/services/auth-service";

export const load: LayoutLoad = async () => {
  const authService = new AuthService(getBackendClient());
  if (authService.isTokenValid) {
    return redirect(303, "/user");
  }
};
