import { useEffect, useState } from "react";
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
  BooleanInput,
  BooleanField,
  SelectInput,
  required,
  FunctionField,
  useShowController,
} from "react-admin";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import SearchIcon from "@mui/icons-material/Search";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import { MarkdownField, MarkdownInput } from "../markdown";

const ThingEndpointsQuery = `
PREFIX iot: <http://iotschema.org/>
PREFIX td: <https://www.w3.org/2019/wot/td#>

SELECT ?thing ?title ?target ?name ?type WHERE {
  GRAPH ?g {
    ?thing ?type ?affordance .
    ?affordance td:name ?name .
    ?affordance td:hasForm ?form .
    ?form <https://www.w3.org/2019/wot/hypermedia#hasTarget> ?target .
    ?thing td:hasActionAffordance ?affordance .
    ?thing td:title ?title .
  }
} 
`;

const SelectThingEndpoints = () => {
  const [endpoints, setEndpoints] = useState([]);
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("/api/registry/sparql", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: new URLSearchParams({ query: ThingEndpointsQuery }),
      });
      const data = await response.json();
      setEndpoints(data.results.bindings);
    };
    fetchData();
  }, []);

  return (
    <SelectInput
      validate={required()}
      source="dataAddress.baseUrl"
      label="Thing Endpoint"
      choices={endpoints.map((e) => ({
        id: e.target.value,
        name: e.target.value + " - " + e.title.value,
      }))}
    />
  );
};

export const AssetsList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid
      style={{ tableLayout: "fixed" }}
      bulkActionButtons={false}
      rowClick="show"
    >
      <TextField source="id" sortable={false} />
      <TextField source="properties.name" label="Name" sortable={false} />
      <TextField
        source="properties.type"
        label="Type"
        sortable={false}
        defaultValue="-"
      />
      <TextField
        source="dataAddress.type"
        label="Data Address Type"
        sortable={false}
      />
    </Datagrid>
  </List>
);

const AssetShowBar = () => {
  return (
    <TopToolbar>
      <DeleteButton mutationMode="pessimistic" />
    </TopToolbar>
  );
};

export const AssetShow = () => {
  const { record } = useShowController();
  return (
    <Show actions={<AssetShowBar />}>
      <SimpleShowLayout>
        <FunctionField
          source="name"
          render={(record) => (
            <>
              <Typography variant="h6">{record.properties.name}</Typography>
              <Typography variant="caption">{record.id}</Typography>
            </>
          )}
        />
        <Labeled fullWidth label="Description">
          <MarkdownField source="properties.description" />
        </Labeled>
        <Labeled fullWidth label="Type">
          <TextField source="properties.type" defaultValue="-" />
        </Labeled>
        <Labeled fullWidth label="Content Type">
          <TextField source="properties.contenttype" />
        </Labeled>
        <Labeled fullWidth label="Type">
          <TextField source="dataAddress.type" />
        </Labeled>
        <Labeled fullWidth label="Base URL">
          <TextField source="dataAddress.baseUrl" />
        </Labeled>
        <Labeled fullWidth label="Accept Header">
          <TextField source="dataAddress.header:Accept" emptyText="-" />
        </Labeled>
        <Labeled fullWidth label="Proxy Path">
          <FunctionField
            source="dataAddress.proxyPath"
            render={(record) => (
              <BooleanField
                record={{ value: record.dataAddress.proxyPath === "true" }}
                source="value"
              />
            )}
          />
        </Labeled>
        <Labeled fullWidth label="Proxy Query Params">
          <FunctionField
            source="dataAddress.proxyQueryParams"
            render={(record) => (
              <BooleanField
                record={{
                  value: record.dataAddress.proxyQueryParams === "true",
                }}
                source="value"
              />
            )}
          />
        </Labeled>
        <Labeled fullWidth label="Proxy Body">
          <FunctionField
            source="dataAddress.proxyBody"
            render={(record) => (
              <BooleanField
                record={{ value: record.dataAddress.proxyBody === "true" }}
                source="value"
              />
            )}
          />
        </Labeled>
        <Labeled fullWidth label="Proxy Method">
          <FunctionField
            source="dataAddress.proxyMethod"
            render={(record) => (
              <BooleanField
                record={{
                  value: record.dataAddress.proxyMethod === "true",
                }}
                source="value"
              />
            )}
          />
        </Labeled>
      </SimpleShowLayout>
    </Show>
  );
};

export const AssetCreate = () => {
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Create>
      <SimpleForm>
        <Dialog open={open} onClose={handleClose} fullWidth>
          <DialogTitle>Select Endpoint</DialogTitle>
          <DialogContent>
            <SelectThingEndpoints />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>OK</Button>
          </DialogActions>
        </Dialog>
        <TextInput
          source="properties.name"
          label="Asset Name"
          fullWidth
          validate={required()}
        />
        <Labeled fullWidth label="Description">
          <MarkdownInput source="properties.description" label="Description" />
        </Labeled>
        <TextInput source="properties.type" label="Asset Type" fullWidth />
        <TextInput
          source="properties.contenttype"
          label="Content Type"
          defaultValue="application/json"
          helperText="The content type of the asset."
          fullWidth
          validate={required()}
        />
        <TextInput
          source="dataAddress.type"
          label="Data Address Type"
          defaultValue="HttpData"
          helperText="The type of the data address e.g. HttpData"
          fullWidth
          validate={required()}
          readOnly
        />
        <TextInput
          source="dataAddress.baseUrl"
          label="Base URL"
          helperText="The base URL of the data address e.g. http://example.com/api/v1/"
          fullWidth
          validate={required()}
          InputProps={{
            endAdornment: window.config.showQuery ? (
              <InputAdornment position="end">
                <IconButton>
                  <SearchIcon onClick={handleClickOpen} />
                </IconButton>
              </InputAdornment>
            ) : null,
          }}
        />
        <TextInput
          source="dataAddress.header:Accept"
          label="Accept Header"
          helperText="The accept header of the data address e.g. application/json"
          fullWidth
        />
        <BooleanInput
          source="dataAddress.proxyPath"
          helperText="Allows specifying additional path segments."
          defaultValue={"false"}
          parse={(v) => (v ? "true" : "false")}
          format={(v) => v === "true"}
        />
        <BooleanInput
          source="dataAddress.proxyQueryParams"
          helperText="Allows specifying query params."
          defaultValue={"false"}
          parse={(v) => (v ? "true" : "false")}
          format={(v) => v === "true"}
        />
        <BooleanInput
          source="dataAddress.proxyBody"
          helperText="Allows attaching a body."
          defaultValue={"false"}
          parse={(v) => (v ? "true" : "false")}
          format={(v) => v === "true"}
        />
        <BooleanInput
          source="dataAddress.proxyMethod"
          helperText="Allows specifying the Http Method (default `GET`)"
          defaultValue={"false"}
          parse={(v) => (v ? "true" : "false")}
          format={(v) => v === "true"}
        />
      </SimpleForm>
    </Create>
  );
};
