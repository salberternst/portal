import { HttpError } from "react-admin";

const fetchThings = async (pagination: any) => {
  const { page, perPage }: { page: number; perPage: number } = pagination;
  const response = await fetch(
    `/api/registry/things?page=${page}&page_size=${perPage}`
  );

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};

const fetchThing = async (id: string) => {
  const response = await fetch(`/api/registry/things/${id}`);

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};

const fetchThingCredentials = async (id: string, security: string) => {
  const response = await fetch(
    `/api/registry/things/${id}/${security}/credentials`
  );

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};

const updateThing = async (thing: any) => {
  const response = await fetch(`/api/registry/things/${thing.id}`, {
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(thing),
    method: "PUT",
  });

  if (response.ok === false) {
    throw new HttpError(response.statusText, response.status);
  }

  return thing;
};

const deleteThing = async (id: any) => {
  const response = await fetch(`/api/registry/things/${id}`, {
    headers: {
      "Content-Type": "application/json",
    },
    method: "DELETE",
  });

  if (response.ok === false) {
    throw new HttpError(response.statusText, response.status);
  }
};

const createThing = async (thing: any) => {
  const response = await fetch(`/api/registry/things`, {
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(thing),
    method: "POST",
  });

  if (response.ok === false) {
    throw new HttpError(response.statusText, response.status);
  }

  return response.json();
};

export default {
  getList: async (resource: any, params: any) => {
    if (resource === "thingDescriptions") {
      const result = await fetchThings(params.pagination);
      return {
        data: result.things,
        total: result.totalElements,
      };
    }
  },
  getMany: async (resource: any, params: any) => {
    if (resource === "thingCredentials") {
      const result = await Promise.all(
        params.ids.map((id) => fetchThingCredentials(params.meta.thingId, id))
      );
      console.log(result);
      return {
        data: result.map((_, index) => ({
          id: params.ids[index],
        })),
      };
    }
  },
  getOne: async (resource: any, params: any) => {
    if (resource === "thingDescriptions") {
      const description = await fetchThing(params.id);
      return {
        data: {
          id: description.id,
          description,
          securityDefinitions: Object.keys(description.securityDefinitions).map(
            (name) => ({
              ...description.securityDefinitions[name],
              name,
              thingId: description.id,
              id: description.id + name,
            })
          ),
          properties: Object.keys(description.properties || {}).map((name) => ({
            ...description.properties[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
          actions: Object.keys(description.actions || {}).map((name) => ({
            ...description.actions[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
          events: Object.keys(description.events || {}).map((name) => ({
            ...description.events[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
        },
      };
    }
  },
  update: async (resource: any, params: any) => {
    if (resource === "thingDescriptions") {
      const { description, ...rest } = params.data;
      const updatedThing = await updateThing({
        ...description,
        properties: rest.properties?.reduce(
          (properties, { name, ...property }) => ({
            ...properties,
            [name]: property,
          }),
          {}
        ),
        actions: rest.actions?.reduce(
          (actions, { name, ...action }) => ({ ...actions, [name]: action }),
          {}
        ),
        events: rest.events?.reduce(
          (events, { name, ...event }) => ({ ...events, [name]: event }),
          {}
        ),
      });
      return {
        data: {
          id: description.id,
          description: updatedThing,
          properties: Object.keys(description.properties || {}).map((name) => ({
            ...description.properties[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
          actions: Object.keys(description.actions || {}).map((name) => ({
            ...description.actions[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
          events: Object.keys(description.events || {}).map((name) => ({
            ...description.events[name],
            name,
            thingId: description.id,
            id: description.id + name,
          })),
        },
      };
    }
  },
  create: (resource: any, params: any) => {
    if (resource === "thingDescriptions") {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = async function (e) {
          createThing(JSON.parse(e.target.result))
            .then((createdThing) => resolve({ data: createdThing }))
            .catch((e) => reject(e));
        };
        reader.readAsText(params.data.attachments.rawFile);
      });
    }
  },
  delete: async (resource: any, params: any) => {
    if (resource === "thingDescriptions") {
      await deleteThing(params.id);
      return {
        data: {
          id: params.id,
        },
      };
    }
  },
};
