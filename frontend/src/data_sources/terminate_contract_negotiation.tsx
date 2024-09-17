import { terminateContractNegotiation } from "../api/contract_negotiations";

export async function create(params) {
  await terminateContractNegotiation(params.data.id, params.data.reason);
  return {
    data: {
      id: params.data.id,
    },
  };
}
