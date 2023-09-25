import React from "react";
import Layout from "@theme/Layout";
import useSWR from "swr";
import Things from "../components/Things";

const fetcher = (url) => fetch(url).then((res) => res.json());


export default function Thingsboard() {
    const { data, error, isLoading } = useSWR("/api/things/", fetcher) 
    
    if (isLoading === true) {
        return <progress style={{ width: "100%" }} />;
    }

    if (error ) {
        return <>{error}</>;
    }


    return (
        <Layout title="Things" description="">
            <main className="main">
                <div>
                    <Things data={data}/>
                </div>
            </main>
        </Layout>
    );
}
