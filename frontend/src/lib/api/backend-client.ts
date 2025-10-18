import PocketBase from "pocketbase";
import { PUBLIC_POCKETBASE_URL, PUBLIC_UNSECURE_COOKIE } from '$env/static/public';

export const getBackendClient = () => {
  const pb = new PocketBase(PUBLIC_POCKETBASE_URL);
  pb.authStore.loadFromCookie(document.cookie);
  pb.authStore.onChange(() => {
    document.cookie = pb.authStore.exportToCookie({
      httpOnly: false,
      secure: PUBLIC_UNSECURE_COOKIE === undefined ||
        PUBLIC_UNSECURE_COOKIE === null ||
        PUBLIC_UNSECURE_COOKIE !== 'true'
    });
  });
  return pb;
}