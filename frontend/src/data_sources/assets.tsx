import {
  createAsset,
  deleteAsset,
  fetchAsset,
  fetchAssets,
} from "../api/assets";

export async function getList(params) {
  const assets = await fetchAssets(params.pagination);
  return {
    data: assets.map((asset) => ({
      ...asset,
      id: asset["@id"],
    })),
    pageInfo: {
      hasNextPage: assets.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const asset = await fetchAsset(params.id);
  return {
    data: {
      ...asset,
      id: asset["@id"],
    },
  };
}

export async function remove(params) {
  await deleteAsset(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  await createAsset({
    ...params.data,
    "@context": {
      "@vocab": "https://w3id.org/edc/v0.0.1/ns/",
    },
  });

  return {
    data: {
      ...params.data,
      id: params.data["@id"],
    },
  };
}

export async function getMany(params) {
  const assets = await Promise.all(params.ids.map((id) => fetchAsset(id)));
  return {
    data: assets.map((asset) => ({
      ...asset,
      id: asset["@id"],
    })),
  };
}
