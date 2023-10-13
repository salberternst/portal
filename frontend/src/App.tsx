import {
  Admin,
  Resource,
  Labeled,
  List,
  Datagrid,
  DateField,
  BooleanField,
  TextField,
  Show,
  Edit,
  SimpleForm,
  TextInput,
  Create,
  ArrayField,
  BooleanInput,
  ArrayInput,
  SimpleFormIterator,
  CustomRoutes,
  Layout,
  Menu,
  FileInput,
  FileField,
  useRecordContext,
  TabbedShowLayout,
  SimpleShowLayout,
  UserMenu,
  AppBar,
  TopToolbar,
  EditButton,
  Button
} from "react-admin";
import { Route } from "react-router-dom";
import {
  Divider,
  Typography,
  Container,
} from "@mui/material";
import DeviceHub from "@mui/icons-material/DeviceHub";
import QueryStatsIcon from "@mui/icons-material/QueryStats";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { EditorState } from "@codemirror/state";
import lzs from "lz-string"

import dataSource from "./data-source";
import authProvider from './auth-provider';
import { SparqlPage } from "./custom_pages/sparql";

const ThingList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <DateField showTime={true} source="createdAt" />
      <TextField source="title" />
      <TextField source="types" />
    </Datagrid>
  </List>
);

const ThingShowProperties = () => {
  const record = useRecordContext();
  return (
    <>
      <Typography variant="h6" sx={{ marginTop: 2 }}>
        Properties
      </Typography>
      <Divider />
      {record.properties?.map((_, index) => (
        <div key={record.properties[index].id}>
          <Labeled fullWidth label="Name">
            <TextField source={`properties.${index}.name`} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`properties.${index}.title`}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField
              label="description"
              source={`properties.${index}.description`}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField source={`properties.${index}.unit`} emptyText="-" />
          </Labeled>
          <ArrayField source={`properties.${index}.forms`}>
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
              <BooleanField source="public" />
            </Datagrid>
          </ArrayField>
        </div>
      ))}
    </>
  );
};

const ThingShowActions = () => {
  const record = useRecordContext();
  return (
    <>
      <Typography variant="h6" sx={{ marginTop: 2 }}>
        Actions
      </Typography>
      <Divider />
      {record.actions?.map((_, index) => (
        <div key={record.actions[index].id}>
          <Labeled fullWidth label="Name">
            <TextField source={`actions.${index}.name`} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`actions.${index}.title`}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField source={`actions.${index}.description`} emptyText="-" />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField source={`actions.${index}.unit`} emptyText="-" />
          </Labeled>
          <ArrayField source={`actions.${index}.forms`}>
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
              <BooleanField source="public" />
            </Datagrid>
          </ArrayField>
        </div>
      ))}
    </>
  );
};

const ThingShowEvents = () => {
  const record = useRecordContext();
  if (record.events?.length === 0) {
    return null;
  }

  return (
    <>
      <Typography variant="h6" sx={{ marginTop: 2 }}>
        Events
      </Typography>
      <Divider />
      {record.events?.map((_, index) => (
        <div key={record.events[index].id}>
          <Labeled fullWidth label="Name">
            <TextField source={`events.${index}.name`} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`events.${index}.title`}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField source={`events.${index}.description`} emptyText="-" />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField source={`events.${index}.unit`} emptyText="-" />
          </Labeled>
          <ArrayField source={`events.${index}.forms`}>
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
              <BooleanField source="public" />
            </Datagrid>
          </ArrayField>
        </div>
      ))}
    </>
  );
};

const ThingShowDescription = () => {
  const record = useRecordContext();
  return (
    <CodeMirror
      value={JSON.stringify(record.description, null, 4)}
      extensions={[json(), EditorState.readOnly.of(true)]}
      basicSetup={{
        lineNumbers: false,
        foldGutter: false,
      }}
      maxHeight="100%"
    />
  );
};

const ThingShowCredentials = () => {
  return (
    <>
      <Typography variant="h6" sx={{ marginTop: 2 }}>Security Definitions</Typography>
      <Divider />
      <ArrayField source="securityDefinitions">
        <Datagrid bulkActionButtons={false}>
          <TextField source="name" />
          <TextField source="scheme" />
          <TextField source="description" emptyText="-" />
        </Datagrid>
      </ArrayField>
    </>
  );
};

const ThingShowLinks = () => (
  <>
    <Typography variant="h6" sx={{ marginTop: 2 }}>Links</Typography>
    <Divider />
    <ArrayField source={`description.links`}>
      <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
        <TextField source="rel" label="Relation" emptyText="-" />
        <TextField source="type" label="Type" emptyText="-" />
        <TextField source="href" label="Link" emptyText="-" />
      </Datagrid>
    </ArrayField>
  </>
)

const ThingShowTitle = () => {
  const record = useRecordContext();
  return (
    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
      {record.description?.title}
    </Typography>
  )
}

const ThingShowActionBar = () => {
  const record = useRecordContext()
  const onClick = () => {
    if(record) {
      const data = "td" + "json" + JSON.stringify(record.description, null, 4)
      const compressed = lzs.compressToEncodedURIComponent(data)
      window.open(`http://plugfest.thingweb.io/playground/#${compressed}`)
    }
  }

  return (
    <TopToolbar>
      <EditButton />
      <Button color="primary" onClick={onClick} label="Open in Editor"/>
    </TopToolbar>
  )
}

const ThingShow = () => {
  return (
    <Show actions={<ThingShowActionBar/>}>
      <SimpleShowLayout>
        <ThingShowTitle />
        <Divider />
        <Labeled fullWidth label="Id">
          <TextField source="description.id" />
        </Labeled>
        <Labeled fullWidth label="Title">
          <TextField source="description.title" />
        </Labeled>
        <Labeled fullWidth label="Description">
          <TextField source="description.description" emptyText="-" />
        </Labeled>
      </SimpleShowLayout>
      <TabbedShowLayout>
        <TabbedShowLayout.Tab label="summary">
          <ThingShowLinks />
          <ThingShowCredentials />
          <ThingShowProperties />
          <ThingShowActions />
          <ThingShowEvents />
        </TabbedShowLayout.Tab>
        <TabbedShowLayout.Tab label="Thing Description">
          <ThingShowDescription />
        </TabbedShowLayout.Tab>
      </TabbedShowLayout>
    </Show>
  );
};

const ThingEdit = () => {
  return (
    <Edit mutationMode="pessimistic">
      <SimpleForm>
        <TextInput fullWidth source="description.id" label="Id" disabled />
        <TextInput fullWidth source="description.title" label="Title" />
        <TextInput
          fullWidth
          source="description.description"
          label="Description"
        />
        <ArrayInput source="properties">
          <SimpleFormIterator fullWidth inline>
            <TextInput fullWidth source="name" />
            <TextInput fullWidth source="description" />
            <ArrayInput source="forms">
              <SimpleFormIterator fullWidth inline>
                <TextInput sx={{ flex: 1 }} source="href" label="Target" />
                <BooleanInput source="public" />
              </SimpleFormIterator>
            </ArrayInput>
          </SimpleFormIterator>
        </ArrayInput>
        <ArrayInput source="actions">
          <SimpleFormIterator fullWidth inline>
            <TextInput fullWidth source="name" />
            <TextInput fullWidth source="description" />
            <ArrayInput source="forms">
              <SimpleFormIterator fullWidth inline>
                <TextInput sx={{ flex: 1 }} source="href" label="Target" />
                <BooleanInput source="public" />
              </SimpleFormIterator>
            </ArrayInput>
          </SimpleFormIterator>
        </ArrayInput>
        <ArrayInput source="events">
          <SimpleFormIterator fullWidth inline>
            <TextInput fullWidth source="name" />
            <TextInput fullWidth source="description" />
            <ArrayInput source="forms">
              <SimpleFormIterator fullWidth inline>
                <TextInput sx={{ flex: 1 }} source="href" label="Target" />
                <BooleanInput source="public" />
              </SimpleFormIterator>
            </ArrayInput>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Edit>
  );
};

const ThingCreate = () => (
  <Create redirect="show">
    <SimpleForm>
      <FileInput source="attachments" accept="application/json">
        <FileField source="src" title="title" />
      </FileInput>
    </SimpleForm>
  </Create>
);

const CustomUserMenu = () => (
  <UserMenu />
);

const CustomAppBar = () => <AppBar userMenu={<CustomUserMenu />} />;

const CustomMenu = () => (
  <Menu>
    <Menu.ResourceItem name="thingDescriptions" />
    <Menu.Item to="/sparql" primaryText="Query" leftIcon={<QueryStatsIcon />} />
  </Menu >
)

const CustomLayout = (props: any) => (
  <Layout menu={CustomMenu} appBar={CustomAppBar}>
    <Container maxWidth="lg">{props.children}</Container>
  </Layout>
);

export const App = () => (
  <Admin loginPage={false} dataProvider={dataSource} layout={CustomLayout} authProvider={authProvider}>
    <CustomRoutes>
      <Route path="/sparql" element={<SparqlPage />} />
    </CustomRoutes>
    <Resource
      name="thingDescriptions"
      options={{ label: "Thing Descriptions" }}
      icon={DeviceHub}
      list={ThingList}
      show={ThingShow}
      edit={ThingEdit}
      create={ThingCreate}
    />
    <Resource name="thingCredentials" icon={DeviceHub} edit={ThingEdit} />
  </Admin>
);
