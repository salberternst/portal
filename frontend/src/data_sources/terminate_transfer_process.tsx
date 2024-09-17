import { terminateTransferProcess } from "../api/transfer_processes";

export async function create(params) {
  await terminateTransferProcess(params.data.id, params.data.reason);
  return {
    data: {
      id: params.data.id,
    },
  };
}
