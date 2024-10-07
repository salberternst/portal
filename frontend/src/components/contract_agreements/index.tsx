import {
  Labeled,
  List,
  Datagrid,
  TextField,
  Show,
  SimpleShowLayout,
  DateField,
  useShowController,
  Button,
  ReferenceField,
  TopToolbar,
  useRecordContext,
} from "react-admin";
import { Link } from "react-router-dom";
import DownloadIcon from "@mui/icons-material/Download";

const ContractAgreementShowBar = () => {
  const record = useRecordContext();
  const enabled = record?.negotiation.state !== "TERMINATED";
  return (
    <TopToolbar>
      <Button
        component={Link}
        to={{
          pathname: "/transferprocesses/create",
        }}
        state={{
          record: {
            counterPartyAddress: record?.negotiation.counterPartyAddress,
            contractId: record?.id,
            assetId: record?.contractAgreement.assetId,
          },
        }}
        disabled={!enabled}
        label="Transfer Dataset"
        startIcon={<DownloadIcon />}
      />
    </TopToolbar>
  );
};

export const ContractAgreementShow = () => {
  const { record, isPending } = useShowController();
  if (isPending) {
    return <div>Loading...</div>;
  }

  return (
    <Show emptyWhileLoading={false} actions={<ContractAgreementShowBar />}>
      <SimpleShowLayout>
        <TextField label="Id" source="id" />
        <TextField label="Type" source="contractAgreement.@type" />
        <TextField label="Asset Id" source="contractAgreement.assetId" />
        <TextField label="Consumer Id" source="contractAgreement.consumerId" />
        <TextField label="Provider Id" source="contractAgreement.providerId" />
        <DateField
          label="Contract Signing Date"
          source="contractAgreement.contractSigningDate"
          transform={(v: number) => new Date(v * 1000)}
          showTime
        />
        <Labeled label="Policy">
          <SimpleShowLayout>
            <TextField label="Type" source="contractAgreement.policy.@type" />
            <TextField
              label="Target"
              source="contractAgreement.policy.odrl:target.@id"
            />
            <TextField
              label="Assignee"
              source="contractAgreement.policy.odrl:assignee"
              emptyText="-"
            />
            <TextField
              label="Assigner"
              source="contractAgreement.policy.odrl:assigner"
              emptyText="-"
            />
          </SimpleShowLayout>
        </Labeled>
        <Labeled label="Asset">
          <SimpleShowLayout>
            <TextField label="Id" source="dataset.@id" />
            <TextField label="Name" source="dataset.name" />
            <TextField label="Content Type" source="dataset.contenttype" />
          </SimpleShowLayout>
        </Labeled>
        <Labeled label="Negotiation">
          <SimpleShowLayout>
            <TextField label="Id" source="negotiation.@id" />
            <TextField label="Type" source="negotiation.@type" />
            <TextField
              label="Contract Agreement Id"
              source="negotiation.contractAgreementId"
            />
            <TextField
              label="Counter Party Address"
              source="negotiation.counterPartyAddress"
            />
            <TextField
              label="Counter Party Id"
              source="negotiation.counterPartyId"
            />
            <TextField label="Protocol" source="negotiation.protocol" />
            <TextField label="State" source="negotiation.state" />
          </SimpleShowLayout>
        </Labeled>
        {record?.negotiation.type === "PROVIDER" && (
          <ReferenceField
            source="contractAgreement.assetId"
            label="Asset"
            reference="assets"
            link="show"
          >
            <TextField source="id" />
          </ReferenceField>
        )}

        <ReferenceField
          label="Contract Negotiation"
          source="negotiation.@id"
          reference="contractnegotiations"
          link="show"
        >
          <TextField source="id" />
        </ReferenceField>
      </SimpleShowLayout>
    </Show>
  );
};

export const ContractAgreementsList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <TextField source="dataset.name" sortable={false} label="Asset" />
      {/* <ReferenceField
        reference="assets"
        source="contractAgreement.assetId"
        label="Asset"
      >
        <TextField source="properties.name"></TextField>
      </ReferenceField> */}
      <TextField
        source="contractAgreement.consumerId"
        sortable={false}
        label="Consumer"
      />
      <TextField
        source="contractAgreement.providerId"
        sortable={false}
        label="Provider"
      />
      <DateField
        label="Contract Signing Date"
        source="contractAgreement.contractSigningDate"
        transform={(v: number) => new Date(v * 1000)}
        showTime
        sortable={false}
      />
    </Datagrid>
  </List>
);
