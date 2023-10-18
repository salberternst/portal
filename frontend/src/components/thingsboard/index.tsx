import { useState, useEffect } from "react";

export default function Thingsboard() {
  const [url, setUrl] = useState();

  useEffect(() => {
    setUrl(new URL(window.location.href));
  }, [setUrl]);

  if (url === undefined) {
    return null;
  }

  return (
    <iframe
      src={`${url.protocol}//thingsboard.${url.host}/`}
      style={{ width: "100%", height: "100%" }}
    />
  );
}
