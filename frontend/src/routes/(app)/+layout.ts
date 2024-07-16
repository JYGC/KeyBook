import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { BackendClient } from "$lib/api/backendclient.svelte";

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

export const load: LayoutLoad = async () => {
  const authManager = new BackendClient();
  if (!authManager.isTokenValid) {
    return redirect(303, "/auth");
  }
  await authManager.authRefresh();
  authManager.pb.authStore.onChange((token, model) => {
    console.log('New store data:', token, model);
  });
};
