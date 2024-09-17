import {
  createDevice,
  deleteDevice,
  fetchDevice,
  fetchDevices,
  updateDevice,
} from "../api/devices";

export async function getList(params) {
  const devices = await fetchDevices(params.pagination);
  return {
    data: devices.data,
    total: devices.totalElements,
  };
}

export async function getOne(params) {
  const device = await fetchDevice(params.id);
  return {
    data: device,
  };
}

export async function remove(params) {
  await deleteDevice(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const device = await createDevice(params.data);
  return {
    data: {
      ...device,
      id: device.id.id,
    },
  };
}

export async function update(params) {
  const updatedDevice = await updateDevice(params.id, params.data);
  return {
    data: updatedDevice,
  };
}
