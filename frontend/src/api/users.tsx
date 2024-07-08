import { HttpError } from "react-admin";

/**
 * Fetches a user from the server.
 * @param id - The ID of the user to fetch.
 * @returns A Promise that resolves to the user object.
 * @throws {HttpError} If the server returns an error response.
 */
export const fetchUser = async (id: string) => {
  const response = await fetch(`/api/portal/users/${id}`, {
    headers: {
      "Content-Type": "application/json",
    },
  });

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};

/**
 * Creates a new user.
 * @param {any} data - The user data to be sent in the request body.
 * @returns {Promise<any>} - A promise that resolves to the JSON response from the server.
 * @throws {HttpError} - If the server returns an error response.
 */
export const createUser = async (data: any) => {
  const response = await fetch(`/api/portal/users`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  const json = await response.json();
  if (response.ok === false) {
    throw new HttpError(json.message, response.status);
  }

  return json;
};

/**
 * Deletes a user with the specified ID.
 * @param {string} id - The ID of the user to delete.
 * @throws {HttpError} If the server returns an error response.
 */
export const deleteUser = async (id: string) => {
  const response = await fetch(`/api/portal/users/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (response.ok === false) {
    const json = await response.json();
    throw new HttpError(json.message, response.status);
  }
};
