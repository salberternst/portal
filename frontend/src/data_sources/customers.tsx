import {
  createCustomer,
  deleteCustomer,
  fetchCustomer,
  fetchCustomers,
  updateCustomer,
} from "../api/customers";

export async function getList(params) {
  const customers = await fetchCustomers(params.pagination);
  return {
    data: customers,
    pageInfo: {
      hasNextPage: customers.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const customer = await fetchCustomer(params.id);
  return {
    data: customer,
  };
}

export async function remove(params) {
  await deleteCustomer(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const customer = await createCustomer(params.data);
  return {
    data: customer,
  };
}

export async function update(params) {
  const updatedCustomer = await updateCustomer(params.id, params.data);
  return {
    data: updatedCustomer,
  };
}

export async function getMany(params) {
  const customers = await Promise.all(
    params.ids.map((id) => fetchCustomer(id))
  );
  return {
    data: customers,
  };
}
