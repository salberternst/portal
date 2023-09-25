import React from "react";
import styles from './index.module.css';
import { useHistory } from '@docusaurus/router';

export default function Things({ data }) {
  const history = useHistory()

  console.log(shit)
  
  return (
    <>
      <div className={styles.container}>
        <div className={styles.actionBar}>
          <button className="button button--primary" onClick={() => history.push('/things/create')}>Hinzufügen</button>
        </div>
        {data.data.length > 0 &&
          <>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Label</th>
                  <th>Model</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {data.data.map((row) => (
                  <tr key={row.id.id}>
                    <td>{row.id.id}</td>
                    <td>{row.name}</td>
                    <td>{row.type}</td>
                    <td>{row.label}</td>
                    <td>{row.model}</td>
                    <td>
                      <div className="button-group button-group--block">
                        <button className="button button--secondary button--sm" onClick={() => history.push('/things/create')}>Info</button>
                        <button className="button button--secondary button--sm">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div>
              <div className="button-group button-group--block">
                <button className="button button--secondary button--sm">{"<"}</button>
                <button className="button button--secondary button--sm">{">"}</button>
              </div>
            </div>
          </>
        }
        {data.data.length === 0 &&
          <p>Keine Things gefunden.</p>
        }
      </div>
    </>
  );
}
