import {
  fetchTransferProcess,
  fetchTransferProcesses,
  createTransferProcess,
} from "../api/transfer_processes";

export async function create(params) {
  const transferProcess = await createTransferProcess({
    ...params.data,
    "@context": {
      "@vocab": "https://w3id.org/edc/v0.0.1/ns/",
    },
    "@type": "TransferRequestDto",
  });
  return {
    data: {
      ...transferProcess,
      id: transferProcess["@id"],
    },
  };
}

export async function getList(params) {
  const transfers = await fetchTransferProcesses(params.pagination);
  return {
    data: transfers.map((transfer) => ({
      ...transfer,
      id: transfer["@id"],
    })),
    pageInfo: {
      hasNextPage: transfers.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const transfer = await fetchTransferProcess(params.id);
  return {
    data: {
      ...transfer,
      id: transfer["@id"],
    },
  };
}

export async function getMany(params) {
  const transfers = await Promise.all(
    params.ids.map((id) => fetchTransferProcess(id))
  );
  return {
    data: transfers.map((transfer) => ({
      ...transfer,
      id: transfer["@id"],
    })),
  };
}
