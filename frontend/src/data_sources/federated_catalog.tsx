import { fetchFederatedCatalog } from "../api/federated_catalog";

export async function getList(params) {
  const federatedCatalog = await fetchFederatedCatalog();
  return {
    data: federatedCatalog,
    total: federatedCatalog.length,
  };
}
