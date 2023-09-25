import React from "react";
import clsx from 'clsx';
import useSWR from "swr";
import siteConfig from '@generated/docusaurus.config';

import styles from './styles.module.css';

const fetcher = (url) => fetch(url).then((res) => res.json());

export default function UserWelcome() {
  const userInfo = siteConfig.customFields.env === "s" ? useSWR("/api/user/info", fetcher) : {
    data: {
      name: "Example User"
    }
  };

  if (userInfo.error) {
    return <>Failed to load</>;
  }

  if (!userInfo.data) {
    return null;
  }

  return (
    <div className="container">
      <h2 className={clsx('hero__title', styles.container)}>Hello {userInfo.data.name}!</h2>
      <div className={clsx('row', styles.container)}>
      </div>
    </div>
  );
}
