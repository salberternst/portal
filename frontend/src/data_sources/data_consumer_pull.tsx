import { fetchTransferProcessDataConsumerPull} from "../api/transfer_processes";

export async function getOne(params) {
  const dataRequest = await fetchTransferProcessDataConsumerPull(params.id);
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
        fetchTransferProcessDataConsumerPull(id).then((dataRequest: any) => ({
          ...dataRequest,
          id,
        }))
      )
    );
  
    return {
      data: dataRequests,
    };
  }
  
  