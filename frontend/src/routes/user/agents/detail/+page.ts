import type { PageLoad } from './$types';

export const load: PageLoad = ({ url }) => ({
  agentId: url.searchParams.get('id') ?? '',
});
