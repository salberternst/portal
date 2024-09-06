import {
  fetchContractAgreement,
  fetchContractAgreements,
  fetchContractAgreementNegotiation,
} from "../api/contract_agreements";
import { fetchCatalogDataset } from "../api/catalog";
import { fetchAsset } from "../api/assets";

export async function getList(params) {
  const contracts = await fetchContractAgreements(params.pagination);
  return {
    data: contracts.map((contract) => ({
      id: contract["@id"],
      contractAgreement: {
        ...contract,
        id: contract["@id"],
      },
      negotiation: {},
      dataset: {},
    })),
    pageInfo: {
      hasNextPage: contracts.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const contractAgreement = await fetchContractAgreement(params.id);
  const negotiation = await fetchContractAgreementNegotiation(
    contractAgreement["@id"]
  );
  let dataset;

  if (negotiation.type === "CONSUMER") {
    dataset = await fetchCatalogDataset(
      negotiation["counterPartyAddress"],
      contractAgreement.assetId
    );
  } else {
    const asset = await fetchAsset(contractAgreement.assetId);
    dataset = {
      "@id": asset["@id"],
      name: asset.properties.name,
      contenttype: asset.properties.contenttype,
    };
  }

  return {
    data: {
      contractAgreement,
      negotiation,
      dataset,
      id: contractAgreement["@id"],
    },
  };
}

export async function getMany(params) {
  const contracts = await Promise.all(
    params.ids.map((id) => fetchContractAgreement(id))
  );
  return {
    data: contracts.map((contract) => ({
      contractAgreement: {
        ...contract,
        id: contract["@id"],
      },
      negotiation: {},
      dataset: {},
      id: contract["@id"],
    })),
  };
}
