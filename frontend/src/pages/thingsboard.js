import React from "react";
import Layout from "@theme/Layout";

export default function Thingsboard() {
    return (
        <Layout title="Thingsboard" description="Landing Page">
            <iframe src="http://thingsboard.example-tenant.192-168-178-60.nip.io/" style={{ position: "absolute", width: "100%", height: "100%" }} />
        </Layout>
    );
}
