import {
  Labeled,
  List,
  Datagrid,
  DateField,
  TextField,
  Show,
  Edit,
  SimpleForm,
  Create,
  ArrayField,
  useRecordContext,
  TabbedShowLayout,
  SimpleShowLayout,
  TopToolbar,
  EditButton,
  Button,
  useInput,
  useShowController,
  DeleteWithConfirmButton,
} from "react-admin";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { EditorState } from "@codemirror/state";
import lzs from "lz-string";
import SwaggerUI from "swagger-ui-react";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";

export const ThingDescriptionList = () => (
  <List
    empty={false}
    hasCreate={true}
    exporter={false}
    sort={{
      field: "createdAt",
      order: "DESC",
    }}
  >
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <DateField showTime={true} source="createdAt" label="Created At" />
      <DateField showTime={true} source="updatedAt" label="Updated At" />
      <TextField source="title" />
      <TextField source="types" sortable={false} />
    </Datagrid>
  </List>
);

export const ThingDescriptionShowProperties = () => {
  const record = useRecordContext();
  if (Object.keys(record.description.properties || []).length === 0) {
    return null;
  }

  return (
    <>
      {Object.keys(record.description.properties).map((name) => (
        <div
          key={record.description.properties[name]}
          style={{ marginTop: 15 }}
        >
          <Labeled fullWidth label="Name">
            <TextField source={`name`} record={{ name }} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`title`}
              record={record?.description.properties[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField
              source={`description`}
              record={record?.description.properties[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField
              source={`unit`}
              record={record?.description.properties[name]}
              emptyText="-"
            />
          </Labeled>
          <ArrayField
            source={`forms`}
            record={record?.description.properties[name]}
          >
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
            </Datagrid>
          </ArrayField>
          <Divider sx={{ mt: 4, mb: 4 }} />
        </div>
      ))}
    </>
  );
};

export const ThingDescriptionShowActions = () => {
  const record = useRecordContext();
  if (Object.keys(record.description.actions || []).length === 0) {
    return null;
  }

  return (
    <>
      {Object.keys(record.description.actions).map((name) => (
        <div key={record.description.actions[name]}>
          <Labeled fullWidth label="Name">
            <TextField source={`name`} record={{ name }} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`title`}
              record={record?.description.actions[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField
              source={`description`}
              record={record?.description.actions[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField
              source={`unit`}
              record={record?.description.actions[name]}
              emptyText="-"
            />
          </Labeled>
          <ArrayField
            source={`forms`}
            record={record?.description.actions[name]}
          >
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
            </Datagrid>
          </ArrayField>
          <Divider sx={{ mt: 4, mb: 4 }} />
        </div>
      ))}
    </>
  );
};

export const ThingDescriptionShowEvents = () => {
  const record = useRecordContext();
  if (Object.keys(record.description.events || []).length === 0) {
    return null;
  }

  return (
    <>
      {Object.keys(record.description.events).map((name) => (
        <div key={record.description.events[name]}>
          <Labeled fullWidth label="Name">
            <TextField source={`name`} record={{ name }} />
          </Labeled>
          <Labeled fullWidth label="Title">
            <TextField
              source={`title`}
              record={record?.description.events[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Description">
            <TextField
              source={`description`}
              record={record?.description.events[name]}
              emptyText="-"
            />
          </Labeled>
          <Labeled fullWidth label="Unit">
            <TextField
              source={`unit`}
              record={record?.description.events[name]}
              emptyText="-"
            />
          </Labeled>
          <ArrayField
            source={`forms`}
            record={record?.description.events[name]}
          >
            <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
              <TextField source="op" label="Operation" emptyText="-" />
              <TextField source="href" label="Target" />
            </Datagrid>
          </ArrayField>
        </div>
      ))}
    </>
  );
};

export const ThingDescriptionShowDescription = () => {
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

export const ThingDescriptionShowLinks = () => (
  <>
    <ArrayField source={`description.links`}>
      <Datagrid bulkActionButtons={false} hover={false} sx={{}}>
        <TextField source="rel" label="Relation" emptyText="-" />
        <TextField source="type" label="Type" emptyText="-" />
        <TextField source="href" label="Link" emptyText="-" />
      </Datagrid>
    </ArrayField>
  </>
);

export const ThingDescriptionShowTitle = () => {
  const record = useRecordContext();
  return (
    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
      {record.description.title}
    </Typography>
  );
};

export const ThingDescriptionShowActionBar = () => {
  const record = useRecordContext();
  const onClick = () => {
    if (record) {
      const data = "td" + "json" + JSON.stringify(record.description, null, 4);
      const compressed = lzs.compressToEncodedURIComponent(data);
      window.open(`https://playground.thingweb.io/#${compressed}`);
    }
  };

  return (
    <TopToolbar>
      <EditButton />
      <Button
        color="primary"
        onClick={onClick}
        label="TD Playground"
        startIcon={<OpenInNewIcon />}
      />
      <DeleteWithConfirmButton />
    </TopToolbar>
  );
};

export const ThingDescriptionShow = () => {
  const { record } = useShowController();
  return (
    <Show actions={<ThingDescriptionShowActionBar />}>
      <SimpleShowLayout>
        <ThingDescriptionShowTitle />
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
        <TabbedShowLayout.Tab label="Summary">
          <TabbedShowLayout syncWithLocation={false}>
            <TabbedShowLayout.Tab label="Properties" syncWithLocation={false}>
              <ThingDescriptionShowProperties />
            </TabbedShowLayout.Tab>
            <TabbedShowLayout.Tab label="Actions" syncWithLocation={false}>
              <ThingDescriptionShowActions />
            </TabbedShowLayout.Tab>
            <TabbedShowLayout.Tab label="Events" syncWithLocation={false}>
              <ThingDescriptionShowEvents />
            </TabbedShowLayout.Tab>
            <TabbedShowLayout.Tab label="Links" syncWithLocation={false}>
              <ThingDescriptionShowLinks />
            </TabbedShowLayout.Tab>
          </TabbedShowLayout>
        </TabbedShowLayout.Tab>
        <TabbedShowLayout.Tab label="Thing Description">
          <ThingDescriptionShowDescription />
        </TabbedShowLayout.Tab>
        <TabbedShowLayout.Tab label="OpenAPI">
          <SwaggerUI url={`/api/registry/things/${record?.id}/openapi`} />
        </TabbedShowLayout.Tab>
      </TabbedShowLayout>
    </Show>
  );
};

const ThingDescriptionEditor = () => {
  const record = useRecordContext();
  const { field } = useInput({ source: "description" });
  return (
    <CodeMirror
      {...field}
      extensions={[json()]}
      basicSetup={{
        lineNumbers: true,
        foldGutter: true,
      }}
      maxHeight="100%"
      value={JSON.stringify(record.description, null, 4)}
    />
  );
};

export const ThingDescriptionEdit = () => {
  return (
    <Edit mutationMode="pessimistic">
      <SimpleForm sx={{ p: 0 }}>
        <ThingDescriptionEditor />
      </SimpleForm>
    </Edit>
  );
};

const ThingDescriptionCreateEditor = () => {
  const { field } = useInput({ source: "description" });
  return (
    <CodeMirror
      {...field}
      extensions={[json()]}
      basicSetup={{
        lineNumbers: true,
        foldGutter: true,
      }}
      width="100%"
      minHeight="400px"
      style={{
        width: "100%",
        minHeight: "400px",
      }}
    />
  );
};

export const ThingDescriptionCreate = () => (
  <Create redirect="show">
    <SimpleForm>
      <ThingDescriptionCreateEditor />
    </SimpleForm>
  </Create>
);
