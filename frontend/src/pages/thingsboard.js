import React from "react";
import Layout from "@theme/Layout";

export default function Thingsboard() {
    const [url, setUrl] = React.useState()

    React.useEffect(() => {
        setUrl(new URL(window.location.href))
    }, [setUrl]);

    if(url === undefined) {
        return null
    }
    
    return (
        <Layout title="Thingsboard" description="Landing Page">
            <iframe src={`${url.protocol}//thingsboard.${url.host}/`} style={{ position: "absolute", width: "100%", height: "100%" }} />
        </Layout>
    );
}
