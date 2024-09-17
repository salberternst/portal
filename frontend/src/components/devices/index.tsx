import {
  List,
  Datagrid,
  TextField,
  Show,
  SimpleShowLayout,
  TopToolbar,
  Create,
  TextInput,
  SimpleForm,
  DateField,
  BooleanInput,
  BooleanField,
  Labeled,
  AutocompleteInput,
  Edit,
  EditButton,
  useEditController,
  useGetList,
  required,
  useGetIdentity,
  ReferenceField,
  ListActions,
  DeleteWithConfirmButton,
} from "react-admin";

const ThingModelSelector = ({ defaultValue = "" }) => {
  return (
    <AutocompleteInput
      source="thingModel"
      choices={[
        { id: "", name: "None" },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/carbon_dioxide.json",
          name: "Carbon Dioxide",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/carbon_monoxide.json",
          name: "Carbon Monoxide",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/connectivity.json",
          name: "Connectivity",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/data_rate.json",
          name: "Data Rate",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/energy.json",
          name: "Energy",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/enum.json",
          name: "Enum",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/garage.json",
          name: "Garage",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/gas.json",
          name: "Gas",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/generic.json",
          name: "Generic",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/human_activity.json",
          name: "Human Activity",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/humidity.json",
          name: "Humidity",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/moisture.json",
          name: "Moisture",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/motion.json",
          name: "Motion",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/outlet.json",
          name: "Outlet",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/power.json",
          name: "Power",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/solar.json",
          name: "Solar",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/temperature.json",
          name: "Temperature",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/home_assistant/timestamp.json",
          name: "Timestamp",
        },
        {
          id: "https://gitlab.int.smartlivingnext.de/data-space/components/resources/thing-models/-/raw/main/stadtwerke_sb/iMSys.json",
          name: "iMSys (Intelligentes Messsystem)",
        },
      ]}
      defaultValue={defaultValue}
      fullWidth
    />
  );
};

const CustomerSelector = () => {
  const { data, isLoading } = useGetList("customers");
  if (isLoading) return null;

  return (
    <AutocompleteInput
      source="customerId"
      choices={data}
      optionText="name"
      optionValue="id"
      emptyText="None"
      fullWidth
      isLoading={isLoading}
    />
  );
};

export const DeviceEdit = () => {
  const { record } = useEditController();
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }
  const isAdmin = identity?.groups.includes("role:admin");

  return (
    <Edit mutationMode="pessimistic">
      <SimpleForm>
        <TextInput source="id" fullWidth disabled />
        <TextInput source="name" fullWidth disabled />
        <TextInput
          source="title"
          fullWidth
          helperText="Enter a title for the device"
        />
        <TextInput
          source="description"
          multiline
          fullWidth
          rows={3}
          helperText="Enter a description for the device"
        />
        <TextInput
          source="category"
          fullWidth
          helperText="Enter a category for the device e.g. dishwasher, fridge, etc."
        />
        <TextInput
          source="manufacturer"
          fullWidth
          helperText="Enter the manufacturer of the device e.g. Samsung, LG, etc."
        />
        <TextInput
          source="model"
          fullWidth
          helperText="Enter the model of the device e.g. 1234, 1234X, etc."
        />
        <BooleanInput
          source="gateway"
          label="Gateway"
          defaultValue={record?.gateway}
        />
        <ThingModelSelector defaultValue={record?.thingModel} />
        {isAdmin && <CustomerSelector />}
      </SimpleForm>
    </Edit>
  );
};

export const DeviceCreate = () => {
  return (
    <Create redirect="show">
      <SimpleForm>
        <TextInput source="name" fullWidth validate={[required()]} />
        <CustomerSelector />
        <BooleanInput source="gateway" label="Gateway" />
        <TextInput
          source="title"
          fullWidth
          helperText="Enter a title for the device"
        />
        <TextInput
          source="description"
          multiline
          fullWidth
          rows={3}
          helperText="Enter a description for the device"
        />
        <ThingModelSelector />
        <TextInput
          source="category"
          fullWidth
          helperText="Enter a category for the device e.g. dishwasher, fridge, etc."
        />
        <TextInput
          source="manufacturer"
          fullWidth
          helperText="Enter the manufacturer of the device e.g. Samsung, LG, etc."
        />
        <TextInput
          source="model"
          fullWidth
          helperText="Enter the model of the device e.g. 1234, 1234X, etc."
        />
      </SimpleForm>
    </Create>
  );
};

const DeviceShowBar = () => {
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }
  const isAdmin = identity?.groups.includes("role:admin");

  return (
    <TopToolbar>
      <EditButton />
      {isAdmin && <DeleteWithConfirmButton />}
    </TopToolbar>
  );
};

export const DeviceShow = () => {
  const { record } = useEditController();
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }
  const isAdmin = identity?.groups.includes("role:admin");

  return (
    <Show actions={<DeviceShowBar />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="name" />
        <DateField source="createdAt" label="Created At" showTime />
        <TextField source="title" label="Title" emptyText="-" />
        <TextField source="description" label="Description" emptyText="-" />
        {record?.customerId && isAdmin && (
          <ReferenceField source="customerId" reference="customers" link="show">
            <TextField source="name" label="Name" />
          </ReferenceField>
        )}
        <TextField source="category" label="Category" emptyText="-" />
        <TextField source="manufacturer" label="Manufacturer" emptyText="-" />
        <TextField source="model" label="Model" emptyText="-" />
        <TextField source="thingModel" label="Thing Model" emptyText="-" />
        <BooleanField source="gateway" label="Gateway" defaultValue={false} />
        {record?.syncStatus && (
          <Labeled label="Sync Status">
            <SimpleShowLayout>
              <TextField source="syncStatus.message" label="Message" />
              <TextField source="syncStatus.status" label="Status" />
              <DateField source="syncStatus.ts" label="Timestamp" showTime />
            </SimpleShowLayout>
          </Labeled>
        )}
        <Labeled label="Credentials">
          <SimpleShowLayout>
            <TextField source="credentials.type" label="Type" />
            <TextField source="credentials.credentials" label="Token" />
          </SimpleShowLayout>
        </Labeled>
      </SimpleShowLayout>
    </Show>
  );
};

export const DevicesList = () => {
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }
  const isAdmin = identity?.groups.includes("role:admin");

  return (
    <List
      empty={false}
      actions={<ListActions hasCreate={isAdmin} />}
      exporter={false}
    >
      <Datagrid bulkActionButtons={false} rowClick="show">
        <TextField source="id" sortable={false} />
        <TextField source="name" label="Name" sortable={false} />
        <BooleanField
          source="additionalInfo.gateway"
          label="Gateway"
          defaultValue={false}
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
};
