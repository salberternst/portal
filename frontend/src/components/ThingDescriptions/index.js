import React from "react";
import styles from './index.module.css';
import lzs from "lz-string"

export default function ThingDescriptions({ data, onThingDeleted }) {
  const [error, setError] = React.useState()

  const deleteThing = async (id) => {
    try {
      const response = await fetch(`/api/registry/things/${id}`, {
        method: 'DELETE'
      })

      if (!response.ok) {
        return setError({
          status: response.status,
          statusText: response.statusText
        })
      }

      if (onThingDeleted !== undefined) {
        onThingDeleted()
      }
    } catch (e) {
      setError(e)
    }
  }

  const showThing = async (id) => {
    try {
      const response = await fetch(`/api/registry/things/${id}`)

      if (!response.ok) {
        return setError({
          status: response.status,
          statusText: response.statusText
        })
      }

      const thing = await response.json()
      const data = "td" + "json" + JSON.stringify(thing, null, 4)
      const compressed = lzs.compressToEncodedURIComponent(data)

      window.open(`http://plugfest.thingweb.io/playground/#${compressed}`);
    } catch (e) {
      setError(e)
    }
  }

  return (
    <>
      <dialog open={error !== undefined}>
        <p>Error deleting thing description</p>
        <p>{JSON.stringify(error)}</p>
        <button onClick={() => setError(undefined)}>OK</button>
      </dialog>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Types</th>
            <th>Created At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data.things.map((row) => (
            <tr key={row.id}>
              <td>{row.id}</td>
              <td>{row.title}</td>
              <td>{row.types}</td>
              <td>{row.createdAt}</td>
              <td>
                <button onClick={() => showThing(row.id)}>Show</button>
                <button onClick={() => deleteThing(row.id)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
