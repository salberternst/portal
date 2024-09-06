import { createUser, deleteUser, fetchUser } from "../api/users";

export async function getOne(params) {
  const user = await fetchUser(params.id);
  return {
    data: {
      ...user,
      id: user.id,
    },
  };
}

export async function remove(params) {
  await deleteUser(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const user = await createUser(params.data);
  return {
    data: {
      ...user,
      id: user.id,
    },
  };
}

export async function getMany(params) {
  const users = await Promise.all(params.ids.map((id: any) => fetchUser(id)));
  return {
    data: users,
  };
}
