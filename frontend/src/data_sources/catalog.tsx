import { fetchCatalog } from "../api/catalog";
import { castArray } from "lodash";

function normalize(catalog) {
  const datasets = [];
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

  return datasets;
}

export async function getOne(params) {
  const catalog = await fetchCatalog(params.id);
  return {
    data: {
      id: params.id,
      datasets: normalize(catalog),
    },
    id: params.id,
  };
}
