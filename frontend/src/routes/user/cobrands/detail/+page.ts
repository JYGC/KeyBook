import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => ({
  cobrandId: url.searchParams.get('id') ?? '',
});
