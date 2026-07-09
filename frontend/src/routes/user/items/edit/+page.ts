import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => ({
  itemId: url.searchParams.get('id') ?? '',
});
