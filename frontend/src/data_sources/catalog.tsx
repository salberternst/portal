import { fetchCatalog } from "../api/catalog";

export async function getOne(params) {
  const catalog = await fetchCatalog(params.id);
  return {
    data: {
      id: params.id,
      ...catalog,
    },
    id: params.id,
  };
}
