import React from "react";

export default function ThingInbox({ data, onThingApproved }) {
  const [error, setError] = React.useState()

  const approveThing = async (id) => {
    try {
      const response = await fetch(`/api/registry/inbox/${id}/approve`, {
        method: 'POST'
      })

      if(!response.ok) {
        return setError({ 
          status: response.status,
          statusText: response.statusText
        })
      }

      if (onThingApproved !== undefined) {
        onThingApproved()
      } 

    } catch (e) {
      setError(e)
    }
  }

  return (
    <>
      <dialog open={error !== undefined}>
        <p>Error importing thing description</p>
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
              <td><button onClick={() => approveThing(row.id)}>Approve</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
}
