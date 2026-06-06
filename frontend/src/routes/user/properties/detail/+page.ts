import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => ({
  propertyId: url.searchParams.get('id') ?? '',
});
