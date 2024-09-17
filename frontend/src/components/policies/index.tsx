import React from "react";
import {
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
  SelectInput,
  DateField,
  ArrayField,
  required,
  AutocompleteInput,
  ArrayInput,
  SimpleFormIterator,
  DateTimeInput,
  useInput,
  AutocompleteArrayInput,
  FormDataConsumer,
  useArrayInput,
} from "react-admin";
import { countries, getEmojiFlag } from "countries-list";
import Grid from "@mui/material/Grid";
import Menu from "@mui/material/Menu";
import Accordion from "@mui/material/Accordion";
import AccordionDetails from "@mui/material/AccordionDetails";
import AccordionSummary from "@mui/material/AccordionSummary";
import IconButton from "@mui/material/IconButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import MenuItem from "@mui/material/MenuItem";
import Typography from "@mui/material/Typography";
import AddCircleOutlineIcon from "@mui/icons-material/AddCircleOutline";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import PermIdentityIcon from "@mui/icons-material/PermIdentity";
import LocationOnIcon from "@mui/icons-material/LocationOn";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";

const countriesList = Object.keys(countries).map((key) => ({
  id: key,
  name: `${getEmojiFlag(key)} ${countries[key].name}`,
}));

const CustomAddButton = () => {
  const context = useArrayInput();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };
  const handleAdd = (event) => {
    context.append({ type: event.target.getAttribute("value") });
    setAnchorEl(null);
  };

  return (
    <>
      <IconButton color="primary" onClick={handleClick}>
        <AddCircleOutlineIcon />
      </IconButton>
      <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
        <MenuItem value="time" onClick={handleAdd} disableRipple>
          <ListItemIcon>
            <AccessTimeIcon />
          </ListItemIcon>
          Time Based
        </MenuItem>
        <MenuItem value="identity" onClick={handleAdd} disableRipple>
          <ListItemIcon>
            <PermIdentityIcon />
          </ListItemIcon>
          Identity Based
        </MenuItem>
        <MenuItem value="location" onClick={handleAdd} disableRipple>
          <ListItemIcon>
            <LocationOnIcon />
          </ListItemIcon>
          Location Based
        </MenuItem>
      </Menu>
    </>
  );
};

const IdentityBasedPermission = () => {
  const {
    field: { value: operatorValue },
  } = useInput({ source: "operator" });
  const selectMultiple =
    operatorValue === "odrl:isNoneOf" || operatorValue === "odrl:isPartOf";

  const renderInput = () => {
    if (selectMultiple) {
      return (
        <ArrayInput source="rightOperand" label="Identities">
          <SimpleFormIterator inline>
            <TextInput label="Identity" />
          </SimpleFormIterator>
        </ArrayInput>
      );
    }
    return <TextInput source="rightOperand" label="Identity" />;
  };

  return (
    <>
      <Typography sx={{ mt: 2 }}>Identity Based Permission</Typography>
      <TextInput
        source="leftOperand"
        defaultValue="edc:participantId"
        label={false}
        sx={{ display: "none" }}
      />
      <SelectInput
        source="operator"
        label="Operator"
        choices={[
          { id: "odrl:eq", name: "Equals" },
          { id: "odrl:neq", name: "Not Equals" },
          { id: "odrl:isNoneOf", name: "Is None Of" },
          { id: "odrl:isPartOf", name: "Is Part Of" },
        ]}
        helperText="Select the operator"
        validate={[required()]}
      />
      {renderInput()}
    </>
  );
};

const TimeBasedPermission = () => {
  return (
    <>
      <Typography sx={{ mt: 2 }}>Time Based Permission</Typography>
      <TextInput
        source="leftOperand"
        defaultValue="edc:policyEvaluationTime"
        label={false}
        sx={{ display: "none" }}
      />
      <SelectInput
        source="operator"
        choices={[
          { id: "odrl:gt", name: "After" },
          { id: "odrl:lt", name: "Before" },
        ]}
        helperText="Select the operator"
        validate={[required()]}
      />
      <DateTimeInput
        source="rightOperand"
        helperText="Select the date"
        validate={[required()]}
      />
    </>
  );
};

const LocationBasedPermission = () => {
  const {
    field: { value: operatorValue },
  } = useInput({ source: "operator" });
  const selectMultiple =
    operatorValue === "odrl:isNoneOf" || operatorValue === "odrl:isPartOf";

  const renderInput = () => {
    if (selectMultiple) {
      return (
        <AutocompleteArrayInput
          source="rightOperand"
          choices={countriesList}
          label="Select the countries"
        />
      );
    }
    return (
      <AutocompleteInput
        source="rightOperand"
        choices={countriesList}
        label="Select the country"
      />
    );
  };

  return (
    <>
      <Typography sx={{ mt: 2 }}>Location Based Permission</Typography>
      <TextInput
        source="leftOperand"
        defaultValue="edc:country"
        label={false}
        sx={{ display: "none" }}
      />
      <SelectInput
        source="operator"
        choices={[
          { id: "odrl:eq", name: "Equals" },
          { id: "odrl:neq", name: "Not Equals" },
          { id: "odrl:isNoneOf", name: "Is None Of" },
          { id: "odrl:isPartOf", name: "Is Part Of" },
        ]}
        validate={[required()]}
      />
      {renderInput()}
    </>
  );
};

const PolicyShowBar = () => {
  return (
    <TopToolbar>
      <DeleteButton mutationMode="pessimistic" />
    </TopToolbar>
  );
};

export const PoliciesList = () => (
  <List empty={false} exporter={false}>
    <Datagrid
      style={{ tableLayout: "fixed" }}
      bulkActionButtons={false}
      rowClick="show"
    >
      <TextField source="id" sortable={false} />
      <TextField
        source="privateProperties.name"
        label="Name"
        sortable={false}
      />
      <DateField
        showTime={true}
        source="createdAt"
        label="Created At"
        sortable={false}
      />
    </Datagrid>
  </List>
);

export const PolicyShow = () => {
  const accordionSummaryStyle = {
    padding: 0,
    "&.MuiAccordionSummary-root": {
      minHeight: "auto",
    },
    "& .MuiAccordionSummary-content": {
      margin: 0,
    },
  };

  return (
    <Show actions={<PolicyShowBar />}>
      <SimpleShowLayout>
        <SimpleShowLayout>
          <TextField source="id" />
          <TextField label="Name" source="privateProperties.name" />
          <TextField
            label="Description"
            source="privateProperties.description"
            emptyText="-"
          />
          <DateField source="createdAt" label="Created At" showTime />
          <TextField label="Type" source="@type" />
          <TextField label="Policy Type" source="policy.@type" />
          <Accordion square elevation={0} disableGutters defaultExpanded>
            <AccordionSummary
              expandIcon={<ArrowDropDownIcon />}
              sx={accordionSummaryStyle}
            >
              <Typography variant="caption">Permissions</Typography>
            </AccordionSummary>
            <AccordionDetails sx={{ p: 0 }}>
              <ArrayField source="policy.odrl:permission" label="Permissions">
                <Datagrid
                  bulkActionButtons={false}
                  rowClick={false}
                  hover={false}
                >
                  <TextField
                    source="odrl:action.@id"
                    label="Action"
                    sortable={false}
                  />
                  <ArrayField
                    source="odrl:constraint"
                    label="Constraint"
                    sortable={false}
                  >
                    <Datagrid
                      bulkActionButtons={false}
                      rowClick={false}
                      style={{ tableLayout: "fixed" }}
                      hover={false}
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
        </SimpleShowLayout>
      </SimpleShowLayout>
    </Show>
  );
};

const PermissionCreate = () => {
  const renderPermission = ({ type }) => {
    if (type === "time") {
      return <TimeBasedPermission />;
    }
    if (type === "identity") {
      return <IdentityBasedPermission />;
    }
    if (type === "location") {
      return <LocationBasedPermission />;
    }
    return null;
  };

  return (
    <SimpleFormIterator addButton={<CustomAddButton />}>
      <FormDataConsumer>
        {({ scopedFormData }) => {
          return renderPermission({ type: scopedFormData.type });
        }}
      </FormDataConsumer>
    </SimpleFormIterator>
  );
};

export const PolicyCreate = (props: any) => (
  <Create {...props}>
    <SimpleForm>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <TextInput
            source="privateProperties.name"
            label="Name"
            required
            fullWidth
          />
        </Grid>
        <Grid item xs={12}>
          <TextInput
            source="privateProperties.description"
            label="Description"
            fullWidth
            multiline
            rows={3}
          />
        </Grid>
        <Grid item xs={12}>
          <ArrayInput source="permissions">
            <PermissionCreate />
          </ArrayInput>
        </Grid>
      </Grid>
    </SimpleForm>
  </Create>
);
