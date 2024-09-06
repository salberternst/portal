import { fetchTransferProcessDataRequest } from "../api/transfer_processes";

export async function getOne(params) {
  const dataRequest = await fetchTransferProcessDataRequest(params.id);
  return {
    data: {
      ...dataRequest,
      id: params.id,
    },
  };
}

export async function getMany(params) {
  const dataRequests = await Promise.all(
    params.ids.map((id: any) =>
      fetchTransferProcessDataRequest(id).then((dataRequest: any) => ({
        ...dataRequest,
        id,
      }))
    )
  );

  return {
    data: dataRequests,
  };
}
