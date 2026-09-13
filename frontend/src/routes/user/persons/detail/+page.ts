import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => ({
  personId: url.searchParams.get('id') ?? '',
});
