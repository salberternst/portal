import React from "react";
import Layout from "@theme/Layout";
import ThingCreate from "../../components/ThingCreate";


export default function Thingsboard() {
    const [url, setUrl] = React.useState()
    const data = { things: [] }

    React.useEffect(() => {
        setUrl(new URL(window.location.href))
    }, [setUrl]);

    if(url === undefined) {
        return null
    }
    
    return (
        <Layout title="Things" description="">
            <main className="main">
                <div>
                    <ThingCreate />
                </div>
            </main>
        </Layout>
    );
}
