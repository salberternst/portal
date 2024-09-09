import { List, useListController } from "react-admin";
import Box from "@mui/material/Box";
import { CatalogList } from "../catalog_list";

export const FederatedCatalogList = () => {
  const { data } = useListController();

  return (
    <List empty={true} exporter={false} component={Box} actions={false}>
      <CatalogList record={data} />
    </List>
  );
};
