import { HttpError } from "react-admin";

/**
 * Fetches the federated catalog from the specified API endpoint.
 *
 * @returns {Promise<Array<Object>>} The array of datasets retrieved from the federated catalog.
 * @throws {HttpError} If the response is not successful, an HttpError is thrown with the error message and status code.
 */
export const fetchFederatedCatalog = async () => {
  const response = await fetch(window.config.federatedCatalogUrl, {
    headers: {
      "Content-Type": "application/json",
    },
    method: "POST",
    body: JSON.stringify({
      "@context": {
        edc: "https://w3id.org/edc/v0.0.1/ns/",
      },
      "@type": "QuerySpec",
    }),
  });

  const catalogs = await response.json();
  if (response.ok === false) {
    throw new HttpError(catalogs.message, response.status);
  }

  const datasets = [];
  for (const catalog of catalogs) {
    const { "dcat:dataset": catalogAssets, ...rest } = catalog;
    if (!Array.isArray(catalogAssets)) {
      catalog["dcat:dataset"] = [catalogAssets];
    }

    for (const dataset of catalog["dcat:dataset"]) {
      datasets.push({
        ...dataset,
        ...rest,
        // image: 'https://picsum.photos/536/354'
      });
    }
  }

  return datasets;
};
