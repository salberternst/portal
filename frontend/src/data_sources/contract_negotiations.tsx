import {
  createContractNegotiation,
  fetchContractNegotiation,
  terminateContractNegotiation,
} from "../api/contract_negotiations";

export async function getOne(params) {
  const contractNegotiation = await fetchContractNegotiation(params.id);
  return {
    data: {
      ...contractNegotiation,
      id: contractNegotiation["@id"],
    },
  };
}

export async function create(params) {
  const contractNegotiation = await createContractNegotiation({
    ...params.data,
    "@type": "ContractRequest",
    policy: {
      ...params.data.policy,
      "@context": "http://www.w3.org/ns/odrl.jsonld",
    },
    "@context": {
      "@vocab": "https://w3id.org/edc/v0.0.1/ns/",
      edc: "https://w3id.org/edc/v0.0.1/ns/",
      odrl: "http://www.w3.org/ns/odrl/2/",
    },
  });
  return {
    data: {
      ...contractNegotiation,
      id: contractNegotiation["@id"],
    },
  };
}

export async function getMany(params) {
  const contractNegotiations = await Promise.all(
    params.ids.map((id) => fetchContractNegotiation(id))
  );
  return {
    data: contractNegotiations.map((contractNegotiation) => ({
      ...contractNegotiation,
      id: contractNegotiation["@id"],
    })),
  };
}

export async function terminate(params) {
  await terminateContractNegotiation(params.data.id, params.data.reason);
  return {
    data: {
      id: params.data.id,
    },
  };
}
