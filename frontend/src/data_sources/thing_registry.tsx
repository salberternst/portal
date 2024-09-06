import {
  createThing,
  deleteThing,
  fetchThing,
  fetchThings,
  updateThing,
} from "../api/thing_registry";

export async function getList(params) {
  const things = await fetchThings(params.pagination);
  return {
    data: things.map((thing) => ({
      ...thing,
      id: thing.id,
    })),
    pageInfo: {
      hasNextPage: things.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const thing = await fetchThing(params.id);
  return {
    data: {
      ...thing,
      id: thing.id,
    },
  };
}

export async function remove(params) {
  await deleteThing(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const thing = await createThing(params.data);
  return {
    data: {
      ...thing,
      id: thing.id,
    },
  };
}

export async function update(params) {
  const updatedThing = await updateThing(params.id, params.data);
  return {
    data: {
      ...updatedThing,
      id: updatedThing.id,
    },
  };
}
