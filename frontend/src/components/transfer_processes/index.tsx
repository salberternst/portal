import { useEffect } from "react";
import {
  List,
  Datagrid,
  TextField,
  Show,
  SimpleShowLayout,
  Create,
  SimpleForm,
  TextInput,
  Labeled,
  DateField,
  useShowController,
  ReferenceField,
  SelectInput,
  required,
  FormDataConsumer,
  Button,
  useRedirect,
  useRecordContext,
  TopToolbar,
  useRefresh,
} from "react-admin";
import Alert from "@mui/material/Alert";
import { Link } from "react-router-dom";
import CancelIcon from "@mui/icons-material/Cancel";
import { useFormContext, useWatch } from "react-hook-form";

export const TransferProcessesList = () => (
  <List empty={false} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <ReferenceField
        label="Asset"
        source="assetId"
        reference="federatedcatalog"
        sortable={false}
      >
        <TextField source="name" />
      </ReferenceField>
      <TextField source="type" sortable={false} />
      <TextField source="state" sortable={false} />
      <TextField source="transferType" sortable={false} label="Transfer Type" />
      <TextField
        label="Data Destination Type"
        source="dataDestination.type"
        sortable={false}
      />
      <DateField
        label="State Timestamp"
        source="stateTimestamp"
        sortable={true}
        showTime
      />
    </Datagrid>
  </List>
);

const TransferProcessesShowBar = () => {
  const record = useRecordContext();
  return (
    <TopToolbar>
      <Button
        component={Link}
        color="error"
        to={{
          pathname: `/transferprocesses/${record?.["@id"]}/terminate`,
        }}
        disabled={record?.state !== "STARTED"}
        label="Terminate Transfer Process"
        startIcon={<CancelIcon />}
      />
    </TopToolbar>
  );
};

export const TransferProcessesShow = () => {
  const { record } = useShowController();
  const refresh = useRefresh();

  useEffect(() => {
    const interval = setInterval(refresh, 5000);
    return () => clearInterval(interval);
  }, [refresh]);

  return (
    <Show actions={<TransferProcessesShowBar />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="type" />
        <TextField source="transferType" label="Transfer Type" />
        <DateField source="stateTimestamp" showTime label="State Timestamp" />
        <TextField source="state" />
        {record?.errorDetail && (
          <Labeled label="Error Detail" fullWidth>
            <Alert severity="error">{record?.errorDetail}</Alert>
          </Labeled>
        )}
        {record?.dataDestination && (
          <Labeled label="Data Destination">
            <SimpleShowLayout>
              <TextField source="dataDestination.type" label="Type" />
            </SimpleShowLayout>
          </Labeled>
        )}
        <TextField source="correlationId" />
        <ReferenceField
          label="Contract Agreement"
          source="contractId"
          reference="contractagreements"
          link="show"
        >
          <TextField source="id" />
        </ReferenceField>
        {record?.transferType === "HttpData-PULL" &&
          record?.type === "CONSUMER" &&
          record?.state !== "TERMINATED" && (
            <ReferenceField
              source="id"
              reference="datarequests"
              link="show"
              label="Data Request"
            >
              <TextField source="id" />
            </ReferenceField>
          )}
      </SimpleShowLayout>
    </Show>
  );
};

const TransferTypeHandler = () => {
  const transferType = useWatch({ name: "transferType" });
  const formContext = useFormContext();

  useEffect(() => {
    switch (transferType) {
      case "HttpData-PULL":
        formContext.setValue("dataDestination.type", "HttpData");
        break;
      case "HttpData-PUSH":
        formContext.setValue("dataDestination.type", "HttpData");
        break;
      case "AmazonS3-PUSH":
        formContext.setValue("dataDestination.type", "AmazonS3");
        break;
      default:
        break;
    }
  }, [transferType, formContext]);

  return null;
};

const HttpDataPush = () => (
  <>
    <TextInput
      source="dataDestination.baseUrl"
      label="Base URL"
      helperText="The base URL where the data will be pushed"
      validate={[required()]}
      fullWidth
    />
  </>
);

const AmazonS3Push = () => (
  <>
    <TextInput
      source="dataDestination.region"
      label="Region"
      helperText="The region of the Amazon S3 bucket"
      validate={[required()]}
      fullWidth
    />
    <TextInput
      source="dataDestination.endpointOverride"
      label="Endpoint Override"
      helperText="The endpoint override of the Amazon S3 bucket"
      fullWidth
    />
    <TextInput
      source="dataDestination.bucketName"
      label="Bucket Name"
      helperText="The name of the Amazon S3 bucket"
      validate={[required()]}
      fullWidth
    />
    <TextInput
      source="dataDestination.objectName"
      label="Object Name"
      helperText="The name of the object in the Amazon S3 bucket"
      validate={[required()]}
      fullWidth
    />
    <TextInput
      source="dataDestination.accessKeyId"
      label="Access Key Id"
      helperText="The access key id of the Amazon S3 bucket"
      validate={[required()]}
      fullWidth
    />
    <TextInput
      source="dataDestination.secretAccessKey"
      label="Secret Access Key"
      helperText="The secret access key of the Amazon S3 bucket"
      validate={[required()]}
      fullWidth
    />
  </>
);

export const TransferProcessesCreate = () => {
  return (
    <Create>
      <SimpleForm>
        <TextInput
          label="Counter Party Address"
          source="counterPartyAddress"
          helperText="The address of the counter party"
          fullWidth
        />
        <TextInput
          source="contractId"
          helperText="The contract agreement id"
          fullWidth
        />
        <TextInput source="assetId" helperText="The asset id" fullWidth />
        <TextInput
          source="protocol"
          defaultValue="dataspace-protocol-http"
          helperText="The dataspace protocol to use"
          fullWidth
        />
        <SelectInput
          source="transferType"
          label="Transfer Type"
          validate={[required()]}
          choices={[
            { id: "HttpData-PULL", name: "HttpData-PULL" },
            { id: "HttpData-PUSH", name: "HttpData-PUSH" },
            { id: "AmazonS3-PUSH", name: "AmazonS3-PUSH" },
          ]}
          helperText="The type of transfer"
        />
        <TextInput
          source="dataDestination.type"
          defaultValue="HttpData"
          helperText="The type of data destination"
          fullWidth
          readOnly
        />
        <TransferTypeHandler />
        <FormDataConsumer>
          {({ formData }) => {
            switch (formData.transferType) {
              case "HttpData-PUSH":
                return <HttpDataPush />;
              case "AmazonS3-PUSH":
                return <AmazonS3Push />;
              default:
                return null;
            }
          }}
        </FormDataConsumer>
      </SimpleForm>
    </Create>
  );
};

export const TransferProcessTerminate = () => {
  const { record } = useShowController();

  const redirect = useRedirect();
  const onRedirect = () => {
    return redirect("show", "/transferprocesses", record?.id);
  };
  return (
    <Create resource="terminatetransferprocess" redirect={onRedirect}>
      <SimpleForm>
        <TextInput source="id" defaultValue={record?.id} readOnly />
        <TextInput source="reason" multiline rows={4} />
      </SimpleForm>
    </Create>
  );
};
