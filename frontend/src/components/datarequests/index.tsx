import {
  ReferenceField,
  TextField,
  Show,
  SimpleShowLayout,
  TopToolbar,
  Button,
} from "react-admin";
import { PasswordField } from "../password_field";
import TimelineIcon from "@mui/icons-material/Timeline";
import DownloadIcon from "@mui/icons-material/Download";

const DataRequestShowBar = () => {
  return (
    <TopToolbar>
      <ReferenceField
        source="id"
        reference="rawdataconsumerpull"
        link="show"
        label="Download"
      >
        <Button
          // color="success"
          label="Download"
          startIcon={<DownloadIcon />}
        />
      </ReferenceField>
      <ReferenceField
        source="id"
        reference="dataconsumerpull"
        link="show"
        label="Visualization"
      >
        <Button
          // color="info"
          label="Visualization"
          startIcon={<TimelineIcon />}
        />
      </ReferenceField>
    </TopToolbar>
  );
};
export const DataRequestShow = () => {
  return (
    <Show actions={<DataRequestShowBar />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="endpoint" />
        <TextField source="authType" label="Authentication Type" />
        <PasswordField source="authorization" label="Authentication" />
      </SimpleShowLayout>
    </Show>
  );
};
