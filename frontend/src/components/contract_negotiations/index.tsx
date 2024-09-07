import {
  Create,
  SimpleForm,
  TextInput,
  useInput,
  Labeled,
  SaveButton,
  Toolbar,
  Show,
  SimpleShowLayout,
  TextField,
  useRefresh,
  DateField,
  useShowController,
  LinearProgress,
  useRedirect,
  ReferenceField,
  Button,
} from "react-admin";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { EditorState } from "@codemirror/state";
import { useEffect } from "react";
import Alert from "@mui/material/Alert";
import { Link } from "react-router-dom";

const ContractNegotiationPolicyInput = () => {
  const { field } = useInput({ source: "policy" });
  return (
    <Labeled label="Policy">
      <CodeMirror
        {...field}
        value={JSON.stringify(field.value, null, 4)}
        extensions={[json(), EditorState.readOnly.of(true)]}
        basicSetup={{
          lineNumbers: false,
          foldGutter: false,
        }}
      />
    </Labeled>
  );
};

const ContractNegotiationsCreateToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton alwaysEnable />
  </Toolbar>
);

export const ContractNegotationCreate = () => {
  return (
    <Create>
      <SimpleForm toolbar={<ContractNegotiationsCreateToolbar />}>
        <TextInput
          source="counterPartyAddress"
          label="Counter Party Address"
          fullWidth
        />
        <ContractNegotiationPolicyInput />
        <TextInput
          source="protocol"
          defaultValue={"dataspace-protocol-http"}
          fullWidth
        />
      </SimpleForm>
    </Create>
  );
};

export const ContractNegotationShow = () => {
  const refresh = useRefresh();
  const { error, isLoading, record } = useShowController();

  useEffect(() => {
    const interval = setInterval(refresh, 5000);
    return () => clearInterval(interval);
  }, [refresh]);

  if (isLoading) {
    return <LinearProgress />;
  }

  if (error) {
    return <div>Error!</div>;
  }

  return (
    <Show>
      <SimpleShowLayout>
        <TextField source="@id" label="Id" />
        {record.createdAt && (
          <DateField source="createdAt" showTime label="Created At" />
        )}
        <TextField source="type" />
        <TextField source="counterPartyAddress" label="Counter Party Adress" />
        <TextField source="counterPartyId" label="Counter Party Id" />
        <TextField source="protocol" />
        <TextField source="state" />
        {record.contractAgreementId && (
          <ReferenceField
            source="contractAgreementId"
            reference="contractagreements"
            link="show"
            label="Contract Agreement"
          >
            <TextField source="id" />
          </ReferenceField>
        )}
        {record?.errorDetail && (
          <Labeled label="Error Detail" fullWidth>
            <Alert severity="error">{record?.errorDetail}</Alert>
          </Labeled>
        )}
      </SimpleShowLayout>
      {record?.state === "FINALIZED" && (
        <Button
          component={Link}
          to={{
            pathname: `/contractnegotiations/${record?.["@id"]}/terminate`,
          }}
          color="warning"
          variant="contained"
          label="Terminate"
          disabled={record?.state === "TERMINATED"}
          fullWidth
        />
      )}
    </Show>
  );
};

export const ContractNegotiationTerminate = () => {
  const { record } = useShowController();
  const redirect = useRedirect();
  const onRedirect = () => {
    return redirect("show", "/contractagreements", record?.contractAgreementId);
  };
  return (
    <Create resource="terminatecontractnegotiation" redirect={onRedirect}>
      <SimpleForm>
        <TextInput source="id" defaultValue={record?.id} disabled />
        <TextInput source="reason" multiline rows={4} />
      </SimpleForm>
    </Create>
  );
};
