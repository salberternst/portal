import React from "react";
import styles from './index.module.css';
import { useHistory } from '@docusaurus/router';

const fetcher = (url) => fetch(url).then((res) => res.json());


export default function ThingInfo({ id }) {
  const history = useHistory()
  const thing = useSWR(`/api/thing/${id}`, fetcher) 
  const credentials = useSWR(`/api/thing/${id}/credentials`, fetcher) 

  console.log(thing.data, credentials.data)

  return (
    <>
      <div className={styles.container}>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div style={{ paddingBottom: 12 }}>
            <label>Name</label>
            <input className={styles.textfield} readOnly required />
            <p style={{ fontSize: "0.75rem" }}>Name of the thing</p>
          </div>
          <div>
            <label>Label</label>
            <input className={styles.textfield} readOnly required />
            <p style={{ fontSize: "0.75rem" }}>Label of the thing</p>
          </div>
          <div>
            <label>Thing Model</label>
            <input className={styles.textfield} readOnly />
            <p style={{ fontSize: "0.75rem" }}>URL to the thing model e.g. https://github.com/test.json</p>
          </div>
        </form>
      </div>
    </>
  );
}
