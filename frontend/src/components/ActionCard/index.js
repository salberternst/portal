import * as React from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";

export default function ActionCard({ name, href, to, description } = {}) {
  return (
    <div className="col col--3 margin-bottom--lg">
      <div className={clsx("card")}>
        <div className={clsx("card__image")}>
          <Link to={href}/>
        </div>
        <div className="card__body">
          <h3>{name}</h3>
          <p>{description}</p>
        </div>
        <div className="card__footer">
          <div className="button-group button-group--block">
            <Link className="button button--secondary" href={href} to={to}>
              Open
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
