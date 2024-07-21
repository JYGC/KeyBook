import type { PageLoad } from './$types';

export const load = (async ({ params }) => {
  return {
    propertyId: params.property_id
  };
}) satisfies PageLoad;