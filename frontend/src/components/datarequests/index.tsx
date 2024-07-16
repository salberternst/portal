import {
  TextField,
  Show,
  SimpleShowLayout
} from "react-admin";

export const DataRequestShow = () => {
  return (
    <Show>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="endpoint" />
        <TextField source="authType" label="Authentication Type" />
        <TextField source="authorization" label="Authentication" />
      </SimpleShowLayout>
    </Show>
  );
};
