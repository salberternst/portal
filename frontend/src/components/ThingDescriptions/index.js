import React from "react";
import styles from './index.module.css';
import lzs from "lz-string"

export default function ThingDescriptions({ data, onThingDeleted, onPrevious, onNext, hasNext, hasPrevious }) {
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
        <button className="button button--secondary" onClick={() => setError(undefined)}>OK</button>
      </dialog>
      <div className={styles.container}>
        <div className={styles.actionBar}>
          <button className="button button--primary">Hinzufügen</button>
        </div>
        {data.things.length > 0 &&
          <>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Titel</th>
                  <th>Typen</th>
                  <th>Erstellt am</th>
                  <th>Aktionen</th>
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
                      <div className="button-group button-group--block">
                        <button className="button button--secondary button--block button--sm" onClick={() => showThing(row.id)}>Show</button>
                        <button className="button button--secondary button--block button--sm" onClick={() => deleteThing(row.id)}>Eigenschaften</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div>
              <div className="button-group button-group--block">
                <button disabled={!hasPrevious} className="button button--secondary button--sm" onClick={onPrevious}>{"<"}</button>
                <button disabled={!hasNext} className="button button--secondary button--sm" onClick={onNext}>{">"}</button>
              </div>
            </div>
          </>
        }
        {data.things.length === 0 &&
          <p>Keine Thing Descriptions gefunden.</p>
        }
      </div>
    </>
  );
}
