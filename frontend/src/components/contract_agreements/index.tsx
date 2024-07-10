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
} from "react-admin";
import { Link } from "react-router-dom";

export const ContractAgreementShow = () => {
  const { record } = useShowController();
  return (
    <Show>
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
      </SimpleShowLayout>
      <Button
        component={Link}
        to={{
          pathname: "/transferprocesses/create",
        }}
        state={{
          record: {
            counterPartyAddress: record?.negotiation.counterPartyAddress,
            contractId: record?.id,
            assetId: record?.dataset.id,
          },
        }}
        variant="contained"
        fullWidth
        label="Transfer"
      />
      <SimpleShowLayout>
        <Labeled label="Asset">
          <SimpleShowLayout>
            <TextField label="Id" source="dataset.@id" />
            <TextField label="Name" source="dataset.name" />
            <TextField label="Content Type" source="dataset.contenttype" />
          </SimpleShowLayout>
        </Labeled>
        <Labeled label="Policy">
          <SimpleShowLayout>
            <TextField label="Type" source="contractAgreement.policy.@type" />
            <TextField
              label="Target"
              source="contractAgreement.policy.odrl:target.@id"
            />
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
      </SimpleShowLayout>
      <Button
        component={Link}
        to={{
          pathname: `/contractnegotiations/${record?.negotiation["@id"]}/terminate`,
        }}
        color="warning"
        variant="contained"
        label="Terminate"
        fullWidth
      />
    </Show>
  );
};

export const ContractAgreementsList = () => (
  <List empty={false} hasCreate={true} exporter={false}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" sortable={false} />
      <TextField
        source="contractAgreement.assetId"
        sortable={false}
        label="Asset"
      />
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
