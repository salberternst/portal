import React from "react";
import Layout from "@theme/Layout";
import ThingDescriptions from "../components/ThingDescriptions";
import useSWR from "swr";

import styles from './thing-registry.module.css';
import ThingInbox from "../components/ThingInbox";
import DiscoverThings from "../components/DiscoverThings";

const fetcher = (url) => fetch(url).then((res) => res.json());

export default function ThingRegistry() {
    const swrInbox = useSWR("/api/registry/inbox", fetcher);
    const swrThings = useSWR("/api/registry/things", fetcher);
    const onThingsDiscovered = () => {
        swrInbox.mutate()
    }
    const onThingApproved = () => {
        swrInbox.mutate()
        swrThings.mutate()
    }
    const onThingDeleted = () => {
        swrThings.mutate()
    }

    if (swrInbox.isLoading === true || swrThings.isLoading === true) {
        return <progress style={{ width: "100%" }} />;
    }

    if (swrInbox.error || swrThings.error) {
        return <>{swrInbox.error || swrThings.error}</>;
    }

    return (
        <Layout title="Thing Registry" description="">
            <main className={styles.main}>
                <p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.</p>
                <ThingDescriptions data={swrThings.data} onThingDeleted={onThingDeleted}/>
                <hr/>
                <p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.</p>
                <DiscoverThings onThingsDiscovered={onThingsDiscovered}/>
                <ThingInbox data={swrInbox.data} onThingApproved={onThingApproved}/>
            </main>
        </Layout>
    );
}
