import {
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
  SelectInput,
  DateField,
  ArrayField,
  required,
} from "react-admin";
import Grid from "@mui/material/Grid";

const PolicyShowBar = () => {
  return (
    <TopToolbar>
      <DeleteButton mutationMode="pessimistic" />
    </TopToolbar>
  );
};

export const PoliciesList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid
      style={{ tableLayout: "fixed" }}
      bulkActionButtons={false}
      rowClick="show"
    >
      <TextField source="id" sortable={false} />
      <TextField
        source="privateProperties.name"
        label="Name"
        sortable={false}
      />
      <DateField
        showTime={true}
        source="createdAt"
        label="Created At"
        sortable={false}
      />
    </Datagrid>
  </List>
);

export const PolicyShow = () => (
  <Show actions={<PolicyShowBar />}>
    <SimpleShowLayout>
      <SimpleShowLayout>
        <TextField source="id" />
        <DateField source="createdAt" label="Created At" showTime />
        <TextField label="Type" source="@type" />
        <TextField label="Policy Type" source="policy.@type" />
        <TextField label="Name" source="privateProperties.name" />
        <TextField label="Description" source="privateProperties.description" />
        <ArrayField source="policy.odrl:permission">
          <Datagrid>
            <TextField source="subject" label="Subject" />
            <TextField source="permission" label="Permission" />
          </Datagrid>
        </ArrayField>
        <ArrayField source="policy.odrl:prohibition">
          <Datagrid>
            <TextField source="subject" label="Subject" />
            <TextField source="prohibition" label="Prohibition" />
          </Datagrid>
        </ArrayField>
        <ArrayField source="policy.odrl:obligation">
          <Datagrid>
            <TextField source="subject" label="Subject" />
            <TextField source="obligation" label="Obligation" />
          </Datagrid>
        </ArrayField>
      </SimpleShowLayout>
    </SimpleShowLayout>
  </Show>
);

export const PolicyCreate = (props: any) => (
  <Create {...props}>
    <SimpleForm>
      <Grid container spacing={2}>
        <Grid item xs={3}>
          <SelectInput
            source="policy.@type"
            label="Type"
            choices={[{ id: "odrl:Set", name: "Set" }]}
            fullWidth
            defaultValue={"odrl:Set"}
            validate={[required()]}
            readOnly
          />
        </Grid>
        <Grid item xs={9}>
          <TextInput
            source="privateProperties.name"
            label="Name"
            required
            fullWidth
          />
        </Grid>
      </Grid>
      <TextInput
        source="privateProperties.description"
        label="Description"
        fullWidth
        multiline
        rows={3}
      />
    </SimpleForm>
  </Create>
);
