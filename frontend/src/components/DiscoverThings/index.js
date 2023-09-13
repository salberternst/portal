import React from "react";
import './index.module.css';

export default function DiscoverThings({ onThingsDiscovered, onError }) {
  const [error, setError] = React.useState()
  const [result, setResult] = React.useState()

  const discoverThings = async () => {
    try {
      const response = await fetch('/api/registry/discovery/thingsboard-importer/scan', {
        method: 'POST'
      })
      setResult(await response.json())
    } catch (e) {
      setError(e)
    }
  }

  if(result !== undefined && onThingsDiscovered !== undefined) {
    onThingsDiscovered(result)
  } else if(error !== undefined && onError !== undefined) {
    onError(error)
  }

  return (
    <>
      <button onClick={discoverThings}>Discover Things</button>
    </>

  );
}
