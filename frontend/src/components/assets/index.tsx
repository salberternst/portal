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
} from "react-admin";
import InputAdornment from "@mui/material/InputAdornment";
import IconButton from "@mui/material/IconButton";
import SearchIcon from "@mui/icons-material/Search";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import Button from "@mui/material/Button";

const HistoryEndpointsQuery = `
PREFIX iot: <http://iotschema.org/>
PREFIX td: <https://www.w3.org/2019/wot/td#>

SELECT ?thing ?title ?target ?name  WHERE {
  GRAPH ?g {
    ?affordance iot:recordedBy ?action .
    ?action td:name ?name .
    ?action td:hasForm ?form .
    ?form <https://www.w3.org/2019/wot/hypermedia#hasTarget> ?target .
    ?thing td:hasActionAffordance ?action .
    ?thing td:title ?title .
  }
} 
`;

const SelectHistoryEndpoints = () => {
  const [endpoints, setEndpoints] = useState([]);
  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch("/api/registry/sparql", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: new URLSearchParams({ query: HistoryEndpointsQuery }),
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
      label="History Endpoint"
      choices={endpoints.map((e) => ({
        id: e.target.value,
        name: e.thing.value + " - " + e.title.value,
      }))}
    />
  );
};

export const AssetsList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <TextField source="properties.name" label="Name" sortable={false} />
      <TextField source="properties.type" label="Type" sortable={false} defaultValue="-"/>
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
  return (
    <Show actions={<AssetShowBar />}>
      <SimpleShowLayout>
        <Labeled fullWidth label="ID">
          <TextField source="@id" />
        </Labeled>
        <Labeled label="Properties">
          <SimpleShowLayout>
            <Labeled fullWidth label="Name">
              <TextField source="properties.name" />
            </Labeled>
            <Labeled fullWidth label="Description">
              <TextField source="properties.description" defaultValue="-"/>
            </Labeled>
            <Labeled fullWidth label="Type">
              <TextField source="properties.type" defaultValue="-"/>
            </Labeled>
            <Labeled fullWidth label="Content Type">
              <TextField source="properties.contenttype" />
            </Labeled>
          </SimpleShowLayout>
        </Labeled>
        <Labeled label="Data Address">
          <SimpleShowLayout>
            <Labeled fullWidth label="Type">
              <TextField source="dataAddress.type" />
            </Labeled>
            <Labeled fullWidth label="Base URL">
              <TextField source="dataAddress.baseUrl" />
            </Labeled>
            <Labeled fullWidth label="Proxy Path">
              <BooleanField source="dataAddress.proxyPath" looseValue />
            </Labeled>
            <Labeled fullWidth label="Proxy Query Params">
              <BooleanField source="dataAddress.proxyQueryParams" looseValue />
            </Labeled>
            <Labeled fullWidth label="Proxy Body">
              <BooleanField source="dataAddress.proxyBody" looseValue />
            </Labeled>
            <Labeled fullWidth label="Proxy Method">
              <BooleanField source="dataAddress.proxyMethod" looseValue />
            </Labeled>
          </SimpleShowLayout>
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
            <SelectHistoryEndpoints />
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
        <TextInput
          source="properties.description"
          label="Description"
          fullWidth
          multiline
          rows={4}
        />
        <TextInput
          source="properties.type"
          label="Asset Type"
          fullWidth
        />
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
            endAdornment: (
              <InputAdornment position="end">
                <IconButton>
                  <SearchIcon onClick={handleClickOpen} />
                </IconButton>
              </InputAdornment>
            ),
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
