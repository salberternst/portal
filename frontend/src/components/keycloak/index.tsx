import { useState, useEffect } from "react";

export default function Keycloak() {
  const [url, setUrl] = useState();

  useEffect(() => {
    setUrl(new URL(window.location.href));
  }, [setUrl]);

  if (url === undefined) {
    return null;
  }

  return (
    <iframe
      src={`${url.protocol}//${url.host}/admin/master/console/`}
      style={{ width: "100%", height: "100%" }}
    />
  );
}
