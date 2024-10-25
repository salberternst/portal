import { fetchFederatedCatalog } from "../api/federated_catalog";
import { castArray } from "lodash";

function normalize(federatedCatalog) {
  const datasets = [];
  for (const catalog of federatedCatalog) {
    const { "dcat:dataset": catalogAssets, ...rest } = catalog;
    if (!Array.isArray(catalogAssets)) {
      catalog["dcat:dataset"] = [catalogAssets];
    }

    for (const dataset of catalog["dcat:dataset"]) {
      dataset["odrl:hasPolicy"] = castArray(dataset["odrl:hasPolicy"]);
      for (const policy of dataset["odrl:hasPolicy"]) {
        policy["odrl:permission"] = castArray(policy["odrl:permission"]);
        for (const permission of policy["odrl:permission"]) {
          permission["odrl:constraint"] = castArray(
            permission["odrl:constraint"]
          );
        }
        policy["odrl:prohibition"] = castArray(policy["odrl:prohibition"]);
        for (const prohibition of policy["odrl:prohibition"]) {
          prohibition["odrl:constraint"] = castArray(
            prohibition["odrl:constraint"]
          );
        }
        policy["odrl:obligation"] = castArray(policy["odrl:obligation"]);
        for (const obligation of policy["odrl:obligation"]) {
          obligation["odrl:constraint"] = castArray(
            obligation["odrl:constraint"]
          );
        }
      }
      datasets.push({
        ...dataset,
        ...rest,
      });
    }
  }

  return datasets;
}

export async function getList(params) {
  const federatedCatalog = await fetchFederatedCatalog();
  const normalizedData = normalize(federatedCatalog);
  if (params.filter?.name) {
    const filteredData = normalizedData.filter((dataset) =>
      dataset["name"].toLowerCase().includes(params.filter.name.toLowerCase())
    );
    return {
      data: filteredData,
      total: filteredData.length,
    };
  }

  return {
    data: normalizedData,
    total: normalizedData.length,
  };
}

export async function getMany() {
  const federatedCatalog = await fetchFederatedCatalog();
  const normalizedData = normalize(federatedCatalog);
  return {
    data: normalizedData,
  };
}
