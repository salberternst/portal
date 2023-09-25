import React from "react";
import styles from './index.module.css';
import { useForm } from "react-hook-form"

export default function ThingCreate() {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm()
  const onSubmit = async (data) => {
    try {
      const response = await fetch("/api/things/", {
        method: "POST",
        body: JSON.stringify(data)
      })

      console.log(response.statusText)
    }
    catch (e) {
      console.log(e)
    }
  }

  return (
    <>
      <div className={styles.container}>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div style={{ paddingBottom: 12 }}>
            <label>Name</label>
            <input className={styles.textfield} {...register("name", { required: true })} required />
            <p style={{ fontSize: "0.75rem" }}>Name of the thing</p>
          </div>
          <div>
            <label>Label</label>
            <input className={styles.textfield} {...register("label", { required: true })} required />
            <p style={{ fontSize: "0.75rem" }}>Label of the thing</p>
          </div>
          <div>
            <label>Thing Model</label>
            <input className={styles.textfield} {...register("model", { required: true })} />
            <p style={{ fontSize: "0.75rem" }}>URL to the thing model e.g. https://github.com/test.json</p>
          </div>
          <div className={styles.actionBar}>
            <input type="submit" value="Erstellen" className="button button--primary" />
          </div>
        </form>
      </div>
    </>
  );
}
