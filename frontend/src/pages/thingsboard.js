import React from "react";
import Layout from "@theme/Layout";
import BrowserOnly from '@docusaurus/BrowserOnly';

export default function Thingsboard() {
    return (
        <Layout title="Thingsboard" description="Landing Page">
            <BrowserOnly>
                {() => {
                    const url = new URL(window.location.href)
                    return <iframe src={`${url.protocol}://thingsboard.${url.host}/`} style={{ position: "absolute", width: "100%", height: "100%" }} />
                }}
            </BrowserOnly>
        </Layout>
    );
}
