import {
  List,
  Datagrid,
  TextField,
  SimpleShowLayout,
  ArrayField,
  useRecordContext,
  ImageField,
  FunctionField,
} from "react-admin";
import Accordion from "@mui/material/Accordion";
import AccordionDetails from "@mui/material/AccordionDetails";
import AccordionSummary from "@mui/material/AccordionSummary";
import Typography from "@mui/material/Typography";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import MuiButton from "@mui/material/Button";
import { Link } from "react-router-dom";

const CreateContractNegotiationButton = () => {
  const record = useRecordContext();
  return (
    <MuiButton
      component={Link}
      size="large"
      fullWidth
      variant="contained"
      to={{
        pathname: "/contractnegotiations/create",
      }}
      state={{
        record: {
          policy: {
            "@type": "http://www.w3.org/ns/odrl/2/Offer",
            "@id": record["odrl:hasPolicy"]["@id"],
            assigner: record.participantId,
            obligation: record["odrl:obligation"],
            permission: record["odrl:permission"],
            prohibition: record["odrl:prohibition"],
            target: record.id,
          },
          counterPartyAddress: record["dcat:service"]["dcat:endpointUrl"],
        },
      }}
    >
      Negotiate
    </MuiButton>
  );
};

const PoliciesInformation = () => {
  return (
    <SimpleShowLayout>
      <ArrayField source="odrl:permission" label="Permissions">
        <Datagrid bulkActionButtons={false}>
          <TextField source="@type" label="Type" />
          <TextField source="odrl:action" label="Action" />
          <TextField source="odrl:target" label="Target" />
          <TextField source="odrl:leftOperand" label="Left Operand" />
          <TextField source="odrl:rightOperand" label="Right Operand" />
        </Datagrid>
      </ArrayField>
      <ArrayField source="odrl:prohibition" label="Prohibitions">
        <Datagrid bulkActionButtons={false}>
          <TextField source="@type" />
          <TextField source="odrl:action" label="Action" />
          <TextField source="odrl:target" label="Target" />
          <TextField source="odrl:leftOperand" label="Left Operand" />
          <TextField source="odrl:rightOperand" label="Right Operand" />
        </Datagrid>
      </ArrayField>
      <ArrayField source="odrl:obligation" label="Obligations">
        <Datagrid bulkActionButtons={false}>
          <TextField source="@type" />
          <TextField source="odrl:action" label="Action" />
          <TextField source="odrl:target" label="Target" />
          <TextField source="odrl:leftOperand" label="Left Operand" />
          <TextField source="odrl:rightOperand" label="Right Operand" />
        </Datagrid>
      </ArrayField>
    </SimpleShowLayout>
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

export const FederatedCatalogList = () => {
  const accordionSummaryStyle = {
    "&.MuiAccordionSummary-root": {
      minHeight: 30,
    },
  };

  return (
    <List empty={false} exporter={false}>
      <Datagrid bulkActionButtons={false} hover={false} header={<></>}>
        <SimpleShowLayout>
          <ImageField
            source="image"
            label={<></>}
            sx={{
              "& .RaImageField-image": {
                right: 16,
                position: "absolute",
              },
            }}
          />
          <FunctionField
            source="name"
            render={(record) => (
              <Typography variant="h6">{record.name}</Typography>
            )}
          />
          <TextField source="participantId" sortable={false} emptyText="-" />
          <TextField source="type" emptyText="-" sortable={false} />
          <TextField source="description" sortable={false} emptyText="-" />
          <TextField
            label="Content Type"
            source="contenttype"
            sortable={false}
            emptyText="-"
          />
          <Accordion square elevation={4} disableGutters>
            <AccordionSummary
              expandIcon={<ArrowDropDownIcon />}
              sx={accordionSummaryStyle}
            >
              <Typography>Service</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <ServiceInformation />
            </AccordionDetails>
          </Accordion>
          <Accordion square elevation={4} disableGutters>
            <AccordionSummary
              expandIcon={<ArrowDropDownIcon />}
              sx={accordionSummaryStyle}
            >
              <Typography>Policies</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <PoliciesInformation />
            </AccordionDetails>
          </Accordion>
          <CreateContractNegotiationButton />
        </SimpleShowLayout>
      </Datagrid>
    </List>
  );
};
