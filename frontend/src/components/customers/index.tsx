import {
  Labeled,
  List,
  Datagrid,
  TextField,
  Show,
  SimpleShowLayout,
  TopToolbar,
  DeleteButton,
  Create,
  TextInput,
  SimpleForm,
  EditButton,
  Edit,
  useShowController,
  ArrayField,
  ReferenceField,
} from "react-admin";
import Divider from "@mui/material/Divider";
import Alert from "@mui/material/Alert";

const CustomerShowBar = () => {
  return (
    <TopToolbar>
      <EditButton />
      <DeleteButton mutationMode="pessimistic" />
    </TopToolbar>
  );
};

export const CustomersList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <TextField source="name" label="Name" sortable={false} />
      <TextField source="description" label="Description" sortable={false} />
    </Datagrid>
  </List>
);

export const CustomerShow = () => {
  const { record, isLoading } = useShowController();
  if (isLoading) return null;

  return (
    <Show actions={<CustomerShowBar />}>
      <SimpleShowLayout>
        <Labeled fullWidth label="ID">
          <TextField source="id" />
        </Labeled>
        <Labeled fullWidth label="Name">
          <TextField source="name" />
        </Labeled>
        <Labeled fullWidth label="Description">
          <TextField source="description" emptyText="-" />
        </Labeled>
        <Labeled fullWidth label="Members">
          <ArrayField source="members">
            <Datagrid 
              bulkActionButtons={false} 
              rowClick={false}
            >
              <ReferenceField source="id" reference="users" link="show">
                <TextField source="id" />
              </ReferenceField>
              <TextField source="email" />
              <TextField source="firstName" />
              <TextField source="lastName" />
            </Datagrid>
          </ArrayField>
        </Labeled>
      </SimpleShowLayout>
      {window.config.devicesEnabled && (
        <>
          <Divider>Thingsboard</Divider>
          {record.thingsboard.error !== undefined && (
            <Alert severity="error">{record.thingsboard.error}</Alert>
          )}
          {record.thingsboard.error === undefined && (
            <SimpleShowLayout>
              <TextField source="thingsboard.id" label="ID" />
              <TextField source="thingsboard.title" label="Title" />
              <TextField
                source="thingsboard.country"
                label="Country"
                emptyText="-"
              />
              <TextField source="thingsboard.city" label="City" emptyText="-" />
              <TextField
                source="thingsboard.address"
                label="Address"
                emptyText="-"
              />
              <TextField
                source="thingsboard.phone"
                label="Phone"
                emptyText="-"
              />
              <TextField
                source="thingsboard.email"
                label="Email"
                emptyText="-"
              />
              <TextField source="thingsboard.zip" label="ZIP" emptyText="-" />
            </SimpleShowLayout>
          )}
        </>
      )}
      {window.config.sparqlEnabled && (
        <>
          <Divider>Fuseki</Divider>
          {record.fuseki.error !== undefined && (
            <Alert severity="error">{record.fuseki.error}</Alert>
          )}
          {record.fuseki.error === undefined && (
            <SimpleShowLayout>
              <TextField source="fuseki.name" label="Name" />
              <TextField source="fuseki.state" label="State" />
            </SimpleShowLayout>
          )}
        </>
      )}
    </Show>
  );
};

export const CustomerCreate = () => {
  return (
    <Create redirect="show">
      <SimpleForm>
        <TextInput source="name" label="Name" fullWidth required />
        <TextInput
          source="description"
          label="Description"
          fullWidth
          multiline
          rows={3}
        />
      </SimpleForm>
    </Create>
  );
};

export const CustomerUpdate = () => {
  return (
    <Edit>
      <SimpleForm>
        <Labeled fullWidth label="ID">
          <TextField source="id" />
        </Labeled>
        <Labeled fullWidth label="Name">
          <TextField source="name" />
        </Labeled>
        <TextInput
          source="description"
          label="Description"
          fullWidth
          multiline
          rows={3}
        />
      </SimpleForm>
    </Edit>
  );
};
