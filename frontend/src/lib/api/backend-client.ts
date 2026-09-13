import PocketBase from "pocketbase";
import { PUBLIC_POCKETBASE_URL, PUBLIC_UNSECURE_COOKIE } from '$env/static/public';

export const getBackendClient = () => {
  const backendClient = new PocketBase(PUBLIC_POCKETBASE_URL);
  backendClient.authStore.loadFromCookie(document.cookie);
  backendClient.authStore.onChange(() => {
    document.cookie = backendClient.authStore.exportToCookie({
      httpOnly: false,
      secure: PUBLIC_UNSECURE_COOKIE === undefined ||
        PUBLIC_UNSECURE_COOKIE === null ||
        PUBLIC_UNSECURE_COOKIE !== 'true'
    });
  });
  return backendClient;
}