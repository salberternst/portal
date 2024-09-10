import { useState } from "react";
import {
  Datagrid,
  TextField,
  SimpleShowLayout,
  ArrayField,
  useRecordContext,
  FunctionField,
  Labeled,
  Button,
} from "react-admin";
import Accordion from "@mui/material/Accordion";
import AccordionDetails from "@mui/material/AccordionDetails";
import AccordionSummary from "@mui/material/AccordionSummary";
import Typography from "@mui/material/Typography";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import { Link } from "react-router-dom";
import { MarkdownField } from "../markdown";
import Box from "@mui/material/Box";
import MenuItem from "@mui/material/MenuItem";
import MuiTextField from "@mui/material/TextField";
import Paper from "@mui/material/Paper";
import HandshakeIcon from "@mui/icons-material/Handshake";

interface CreateContractNegotiationButtonProps {
  selectedPolicy: number;
}

const CreateContractNegotiationButton = ({
  selectedPolicy,
}: CreateContractNegotiationButtonProps) => {
  const record = useRecordContext();
  const policy = record?.["odrl:hasPolicy"][selectedPolicy];
  const datasetId = record?.["@id"];
  const counterPartyAddress = record?.["dcat:service"]?.["dcat:endpointUrl"];
  const participantId = record?.["participantId"];

  return (
    <Button
      component={Link}
      color="secondary"
      to={{
        pathname: "/contractnegotiations/create",
      }}
      sx={{
        float: "right",
      }}
      state={{
        record: record && {
          policy: {
            "@type": policy["@type"],
            "@id": policy["@id"],
            assigner: participantId,
            obligation: policy?.["odrl:obligation"],
            permission: policy?.["odrl:permission"],
            prohibition: policy?.["odrl:prohibition"],
            target: datasetId,
          },
          counterPartyAddress,
        },
      }}
      label="Start Negotiation"
      startIcon={<HandshakeIcon />}
    />
  );
};

const ServiceInformation = () => {
  return (
    <SimpleShowLayout>
      <TextField source="dcat:service.@id" label="Id" />
      <TextField source="dcat:service.@type" label="Type" />
      <TextField source="dcat:service.dcat:endpointUrl" label="Endpoint URL" />
      <TextField source="dcat:service.dct:terms" label="Terms" />
      <TextField
        source="dcat:service.dcat:endpointDescription"
        label="Endpoint Description"
      />
    </SimpleShowLayout>
  );
};

const PolicySelect = ({ record, onPolicySelect }) => {
  const choices = record["odrl:hasPolicy"].map((policy, index) => {
    return {
      id: index,
      name: policy["@id"],
    };
  });
  return (
    <MuiTextField
      label="Selected Policy"
      helperText="Select a policy for negotiation"
      size="small"
      fullWidth
      defaultValue={0}
      select
      onChange={onPolicySelect}
      InputLabelProps={{ sx: { p: 0 } }}
    >
      {choices.map((choice) => (
        <MenuItem key={choice.id} value={choice.id}>
          {choice.name}
        </MenuItem>
      ))}
    </MuiTextField>
  );
};

const Permissions = ({ record }) => {
  const accordionSummaryStyle = {
    padding: 0,
    "&.MuiAccordionSummary-root": {
      minHeight: 30,
    },
  };
  if (record?.["odrl:permission"].length === 0) {
    return null;
  }
  return (
    <Accordion square elevation={0} disableGutters>
      <AccordionSummary
        expandIcon={<ArrowDropDownIcon />}
        sx={accordionSummaryStyle}
      >
        <Typography variant="caption">Permissions</Typography>
      </AccordionSummary>
      <AccordionDetails sx={{ p: 0 }}>
        <ArrayField
          source="odrl:permission"
          label="Permissions"
          emptyText="No Permissions"
          record={record}
        >
          <Datagrid
            bulkActionButtons={false}
            rowClick={false}
            hover={false}
            size="small"
          >
            <TextField
              source="odrl:action.@id"
              label="Action"
              sortable={false}
            />
            <ArrayField source="odrl:constraint" label="Constraints">
              <Datagrid
                hover={false}
                bulkActionButtons={false}
                style={{ tableLayout: "fixed" }}
                size="small"
              >
                <TextField
                  source="odrl:leftOperand.@id"
                  label="Left Operand"
                  sortable={false}
                />
                <TextField
                  source="odrl:operator.@id"
                  label="Operator"
                  sortable={false}
                />
                <TextField
                  source="odrl:rightOperand"
                  label="Right Operand"
                  sortable={false}
                />
              </Datagrid>
            </ArrayField>
          </Datagrid>
        </ArrayField>
      </AccordionDetails>
    </Accordion>
  );
};

const PoliciesShow = () => {
  const [selectedPolicy, setSelectedPolicy] = useState(0);
  const record = useRecordContext();
  const handlePolicyChange = (event) => {
    setSelectedPolicy(event.target.value);
  };

  return (
    <>
      <PolicySelect record={record} onPolicySelect={handlePolicyChange} />
      <Permissions record={record?.["odrl:hasPolicy"][selectedPolicy]} />
      <CreateContractNegotiationButton selectedPolicy={selectedPolicy} />
    </>
  );
};

export const CatalogList = ({ record }) => {
  const accordionSummaryStyle = {
    padding: 0,
    "&.MuiAccordionSummary-root": {
      minHeight: 30,
    },
  };

  return (
    <Box display="flex" flexDirection="column" gap={2} sx={{ maxWidth: 1200 }}>
      {record?.map((dataset, index) => (
        <Paper key={dataset?.["@id"] + index}>
          <SimpleShowLayout record={dataset}>
            <FunctionField
              source="name"
              render={(record) => (
                <Typography variant="h6">{record.name}</Typography>
              )}
            />
            <TextField source="participantId" sortable={false} emptyText="-" />
            <TextField source="type" emptyText="-" sortable={false} />
            <Accordion square elevation={0} disableGutters>
              <AccordionSummary
                expandIcon={<ArrowDropDownIcon />}
                sx={accordionSummaryStyle}
              >
                <Typography variant="caption">Details</Typography>
              </AccordionSummary>
              <AccordionDetails sx={{ p: 0, marginLeft: 1 }}>
                <Labeled fullWidth label="Description">
                  <MarkdownField source="description" />
                </Labeled>
                <TextField
                  label="Content Type"
                  source="contenttype"
                  sortable={false}
                  emptyText="-"
                />
                <Accordion square elevation={0} disableGutters>
                  <AccordionSummary
                    expandIcon={<ArrowDropDownIcon />}
                    sx={accordionSummaryStyle}
                  >
                    <Typography variant="caption">Service</Typography>
                  </AccordionSummary>
                  <AccordionDetails sx={{ p: 0, marginLeft: 1 }}>
                    <ServiceInformation />
                  </AccordionDetails>
                </Accordion>
              </AccordionDetails>
            </Accordion>
            <PoliciesShow />
          </SimpleShowLayout>
        </Paper>
      ))}
    </Box>
  );
};
