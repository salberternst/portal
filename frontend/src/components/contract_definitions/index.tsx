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
  ArrayInput,
  SimpleFormIterator,
  ReferenceInput,
  AutocompleteInput,
  required,
  ArrayField,
  ReferenceField,
} from "react-admin";

export const ContractDefinitionsList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid
      bulkActionButtons={false}
      rowClick="show"
      style={{ tableLayout: "fixed" }}
    >
      <TextField source="id" sortable={false} />
      <TextField
        label="Name"
        source="privateProperties.name"
        sortable={false}
      />
    </Datagrid>
  </List>
);

const ContractDefinitionShowBar = () => (
  <TopToolbar>
    <DeleteButton mutationMode="pessimistic" />
  </TopToolbar>
);

export const ContractDefinitionShow = () => (
  <Show actions={<ContractDefinitionShowBar />}>
    <SimpleShowLayout>
      <TextField label="Name" source="privateProperties.name" />
      <TextField source="id" />
      <TextField
        label="Description"
        source="privateProperties.description"
        emptyText="-"
      />
      <TextField label="Type" source="@type" />
      <ReferenceField
        source="accessPolicyId"
        reference="policies"
        link="show"
        label="Access Policy"
      >
        <TextField source="id" />
      </ReferenceField>
      <ReferenceField
        source="contractPolicyId"
        reference="policies"
        label="Contract Policy"
        link="show"
      >
        <TextField source="id" />
      </ReferenceField>
      <ArrayField label="Asset Selector" source="assetsSelector">
        <Datagrid bulkActionButtons={false} rowClick={false} hover={false}>
          <TextField source="operandLeft" label="Operand Left" />
          <TextField source="operator" label="Operator" />
          <ReferenceField
            source="operandRight"
            label="Operand Right"
            reference="assets"
            link="show"
          >
            <TextField source="properties.name" />
          </ReferenceField>
        </Datagrid>
      </ArrayField>
    </SimpleShowLayout>
  </Show>
);

export const ContractDefinitionCreate = (props: any) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput
        source="privateProperties.name"
        label="Name"
        fullWidth
        validate={[required()]}
      />
      <TextInput
        source="privateProperties.description"
        label="Description"
        fullWidth
        multiline
        rows={4}
      />
      <ReferenceInput source="accessPolicyId" reference="policies">
        <AutocompleteInput
          label="Access Policy"
          optionText="privateProperties.name"
          fullWidth
          validate={[required()]}
        />
      </ReferenceInput>
      <ReferenceInput source="contractPolicyId" reference="policies">
        <AutocompleteInput
          label="Contract Policy"
          optionText="privateProperties.name"
          fullWidth
          validate={[required()]}
        />
      </ReferenceInput>
      <ArrayInput
        source="assetsSelector"
        label="Asset Selector"
        fullWidth
        validate={[required()]}
      >
        <SimpleFormIterator inline fullWidth>
          <TextInput
            source="operandLeft"
            defaultValue="https://w3id.org/edc/v0.0.1/ns/id"
            fullWidth
          />
          <TextInput source="operator" defaultValue="=" fullWidth />
          <ReferenceInput source="operandRight" reference="assets">
            <AutocompleteInput
              optionText="properties.name"
              validate={[required()]}
              fullWidth
            />
          </ReferenceInput>
        </SimpleFormIterator>
      </ArrayInput>
    </SimpleForm>
  </Create>
);
