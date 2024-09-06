import { HttpError } from "react-admin";
import { castArray } from "lodash";

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

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};
