import { TextField, Show, SimpleShowLayout } from "react-admin";
import { PasswordField } from "../password_field";

export const DataRequestShow = () => {
  return (
    <Show>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="endpoint" />
        <TextField source="authType" label="Authentication Type" />
        <PasswordField source="authorization" label="Authentication" />
      </SimpleShowLayout>
    </Show>
  );
};
