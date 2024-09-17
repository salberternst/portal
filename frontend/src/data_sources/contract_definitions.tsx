import {
  createContractDefinition,
  deleteContractDefinition,
  fetchContractDefinition,
  fetchContractDefinitions,
} from "../api/contract_definitions";

export async function getList(params) {
  const contracts = await fetchContractDefinitions(params.pagination);
  return {
    data: contracts.map((contract) => ({
      ...contract,
      id: contract["@id"],
    })),
    pageInfo: {
      hasNextPage: contracts.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const contractDefinition = await fetchContractDefinition(params.id);
  return {
    data: {
      ...contractDefinition,
      id: contractDefinition["@id"],
    },
  };
}

export async function remove(params) {
  await deleteContractDefinition(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const contractDefinition = await createContractDefinition({
    ...params.data,
    "@context": {
      "@vocab": "https://w3id.org/edc/v0.0.1/ns/",
    },
  });
  return {
    data: {
      ...contractDefinition,
      id: contractDefinition["@id"],
    },
  };
}

export async function getMany(params) {
  const contracts = await Promise.all(
    params.ids.map((id) => fetchContractDefinition(id))
  );
  return {
    data: contracts.map((contract) => ({
      ...contract,
      id: contract["@id"],
    })),
  };
}
