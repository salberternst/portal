import {
  createThing,
  deleteThing,
  fetchThing,
  fetchThings,
  updateThing,
} from "../api/thing_registry";

export async function getList(params) {
  const result = await fetchThings(params.pagination, params.sort);
  return {
    data: result.things.map((thing: any) => ({
      ...thing,
      description: {},
    })),
    total: result.totalPages * result.pageSize,
  };
}

export async function getOne(params) {
  const description = await fetchThing(params.id);
  return {
    data: {
      id: description.id,
      description,
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
  const createdThing = await createThing(params.data.description);
  return {
    data: {
      id: createdThing.id,
      description: {},
    },
  };
}

export async function update(params) {
  const updatedThing = await updateThing(params.id, params.data.description);
  return {
    data: {
      id: updatedThing.id,
      description: updatedThing,
    },
  };
}
