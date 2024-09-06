import {
  List,
  Datagrid,
  TextField,
  SimpleShowLayout,
  ArrayField,
  useRecordContext,
  ImageField,
  FunctionField,
  Labeled,
} from "react-admin";
import Accordion from "@mui/material/Accordion";
import AccordionDetails from "@mui/material/AccordionDetails";
import AccordionSummary from "@mui/material/AccordionSummary";
import Typography from "@mui/material/Typography";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import MuiButton from "@mui/material/Button";
import { Link } from "react-router-dom";
import { MarkdownField } from "../markdown";

const CreateContractNegotiationButton = ({
  datasetId,
  counterPartyAddress,
  participantId,
}) => {
  const record = useRecordContext();
  return (
    <MuiButton
      component={Link}
      variant="contained"
      color="secondary"
      fullWidth
      to={{
        pathname: "/contractnegotiations/create",
      }}
      state={{
        record: {
          policy: {
            "@type": record["@type"].replace("odrl:", ""),
            "@id": record["@id"],
            assigner: participantId,
            obligation: record["odrl:obligation"],
            permission: record["odrl:permission"],
            prohibition: record["odrl:prohibition"],
            target: datasetId,
          },
          counterPartyAddress,
        },
      }}
    >
      Negotiate
    </MuiButton>
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

const PoliciesShow = () => {
  const record = useRecordContext();
  const EmptyHeader = () => null;
  return (
    <ArrayField source="odrl:hasPolicy" label="Policies">
      <Datagrid
        bulkActionButtons={false}
        rowClick={false}
        hover={false}
        header={EmptyHeader}
      >
        <div>
          <ArrayField
            source="odrl:permission"
            label="Permissions"
            emptyText="No Permissions"
          >
            <Datagrid bulkActionButtons={false} rowClick={false} hover={false}>
              <ArrayField source="odrl:constraint" label="Constraints">
                <Datagrid
                  bulkActionButtons={false}
                  style={{ tableLayout: "fixed" }}
                >
                  <TextField
                    source="odrl:leftOperand.@id"
                    label="Left Operand"
                  />
                  <TextField source="odrl:operator.@id" label="Operator" />
                  <TextField source="odrl:rightOperand" label="Right Operand" />
                </Datagrid>
              </ArrayField>
            </Datagrid>
          </ArrayField>
          <CreateContractNegotiationButton
            datasetId={record["@id"]}
            counterPartyAddress={record["dcat:service"]["dcat:endpointUrl"]}
            participantId={record["participantId"]}
          />
        </div>
      </Datagrid>
    </ArrayField>
  );
};

export const FederatedCatalogList = () => {
  const accordionSummaryStyle = {
    "&.MuiAccordionSummary-root": {
      minHeight: 30,
    },
  };
  const EmptyHeader = () => null;

  return (
    <List empty={false} exporter={false}>
      <Datagrid bulkActionButtons={false} hover={false} header={EmptyHeader}>
        <SimpleShowLayout>
          <ImageField
            source="image"
            label={false}
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
          <Labeled fullWidth label="Description">
            <MarkdownField source="description" />
          </Labeled>
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
          <Accordion square elevation={4} disableGutters defaultExpanded>
            <AccordionSummary
              expandIcon={<ArrowDropDownIcon />}
              sx={accordionSummaryStyle}
            >
              <Typography>Policies</Typography>
            </AccordionSummary>
            <AccordionDetails>
              <PoliciesShow />
            </AccordionDetails>
          </Accordion>
        </SimpleShowLayout>
      </Datagrid>
    </List>
  );
};
