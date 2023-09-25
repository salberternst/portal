import React from "react";
import Layout from "@theme/Layout";
import ThingDescriptions from "../components/ThingDescriptions";
import useSWR from "swr";
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

import styles from './thing-registry.module.css';
import ThingInbox from "../components/ThingInbox";
import DiscoverThings from "../components/DiscoverThings";

const fetcher = (url) => fetch(url).then((res) => res.json());

export default function ThingRegistry() {
    const [thingsPage, setThingsPage] = React.useState(1)
    const swrInbox = useSWR("/api/registry/inbox", fetcher);
    const swrThings = useSWR(`/api/registry/things?page_size=10&page=${thingsPage}`, fetcher);
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
    const onNextThingsPage = () => {
        setThingsPage(thingsPage + 1)
    }
    const onPreviousPage = () => {
        setThingsPage(thingsPage - 1)
    }

    console.log(swrInbox.isLoading === true, swrThings.isLoading === true)

    if (swrInbox.isLoading === true || swrThings.isLoading === true) {
        return <progress style={{ width: "100%" }} />;
    }

    if (swrInbox.error || swrThings.error) {
        return <>{swrInbox.error || swrThings.error}</>;
    }

    return (
        <Layout title="Thing Registry" description="">
            <main className="main">
                <div>
                    <Tabs>
                        <TabItem value="thing-descriptions" label="Thing Descriptions" default>
                            <ThingDescriptions
                                data={swrThings.data}
                                onThingDeleted={onThingDeleted}
                                onNext={onNextThingsPage}
                                onPrevious={onPreviousPage}
                                hasNext={swrThings.data.totalPages > swrThings.data.page}
                                hasPrevious={swrThings.data.page > 1}
                            />
                        </TabItem>
                        <TabItem value="inbox" label="Inbox">
                            <ThingInbox data={swrInbox.data} onThingApproved={onThingApproved} onThingsDiscovered={onThingsDiscovered} />
                        </TabItem>
                    </Tabs>
                </div>
            </main>
        </Layout>
    );
}
