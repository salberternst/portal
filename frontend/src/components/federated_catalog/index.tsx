import { List, SearchInput, useListController } from "react-admin";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import useTheme from "@mui/material/styles/useTheme";
import InboxIcon from "@mui/icons-material/Inbox";
import { CatalogList } from "../catalog_list";

// eslint-disable-next-line react/jsx-key
const postFilters = [<SearchInput source="name" alwaysOn />];

const EmptyFederatedCatalogList = () => {
  const theme = useTheme();
  return (
    <Box textAlign="center">
      <InboxIcon
        sx={{ width: "9em", height: "9em", color: theme.palette.primary.main }}
      />
      <Typography variant="h6" paragraph>
        No Records in Federated Catalog yet
      </Typography>
    </Box>
  );
};

export const FederatedCatalogList = () => {
  const { data } = useListController();

  return (
    <List
      empty={<EmptyFederatedCatalogList />}
      component={Box}
      actions={false}
      filters={postFilters}
    >
      <CatalogList record={data} />
    </List>
  );
};
