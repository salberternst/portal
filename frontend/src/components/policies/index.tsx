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
  BooleanInput,
  useFieldValue,
  FunctionField,
} from "react-admin";
import { countries, getEmojiFlag } from "countries-list";
import Grid from "@mui/material/Grid";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import Typography from "@mui/material/Typography";
import AccordionDetails from "@mui/material/AccordionDetails";
import { useState } from "react";

const countriesList = Object.keys(countries).map((key) => ({
  id: key,
  name: `${getEmojiFlag(key)} ${countries[key].name}`,
}));

const IdentityBasedPermission = () => {
  const [identityBasedPermission, setIdentityBasedPermission] = useState(false);
  const onClick = () => {
    setIdentityBasedPermission(!identityBasedPermission);
  };
  const {
    field: { value: operatorValue },
  } = useInput({ source: "permissions.identity.operator" });
  const selectMultiple =
    operatorValue === "odrl:isAnyOf" ||
    operatorValue === "odrl:isAllOf" ||
    operatorValue === "odrl:isNoneOf";
  const renderInput = () => {
    if (selectMultiple) {
      return (
        <ArrayInput source="permissions.identity.rightOperand">
          <SimpleFormIterator inline>
            <TextInput />
          </SimpleFormIterator>
        </ArrayInput>
      );
    }
    return <TextInput source="permissions.identity.rightOperand" />;
  };

  return (
    <Accordion
      expanded={identityBasedPermission}
      slotProps={{ transition: { unmountOnExit: true } }}
    >
      <AccordionSummary>
        <span>
          <Typography gutterBottom={false}>
            Enable Identity Based Permission
          </Typography>
          <Typography variant="caption">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography>
        </span>
        <span style={{ flex: 1 }} />
        <BooleanInput
          source="permissions.identity.enable"
          onClick={onClick}
          sx={{ alignSelf: "flex-end" }}
          label={false}
        />
      </AccordionSummary>
      <AccordionDetails>
        <TextInput
          source="permissions.identity.leftOperand"
          defaultValue="edc:edc_identity"
          label={false}
          sx={{ display: "none" }}
        />
        <SelectInput
          source="permissions.identity.operator"
          choices={[
            { id: "odrl:eq", name: "Equals" },
            { id: "odrl:neq", name: "Not Equals" },
            { id: "odrl:isAnyOf", name: "Is Any Of" },
            { id: "odrl:isAllOf", name: "Is All Of" },
            { id: "odrl:isNoneOf", name: "Is None Of" },
          ]}
          helperText="Select the operator"
          validate={[required()]}
        />
        {renderInput()}
      </AccordionDetails>
    </Accordion>
  );
};

const TimeBasedPermission = () => {
  const [timeBasedPermission, setTimeBasedPermission] = useState(false);
  const onClick = () => {
    setTimeBasedPermission(!timeBasedPermission);
  };

  return (
    <Accordion
      expanded={timeBasedPermission}
      slotProps={{ transition: { unmountOnExit: true } }}
    >
      <AccordionSummary>
        <span>
          <Typography gutterBottom={false}>
            Enable Time Based Permission
          </Typography>
          <Typography variant="caption">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography>
        </span>
        <span style={{ flex: 1 }} />
        <BooleanInput
          source="permissions.time.enable"
          onClick={onClick}
          sx={{ alignSelf: "flex-end" }}
          label={false}
        />
      </AccordionSummary>
      <AccordionDetails>
        <TextInput
          source="permissions.time.leftOperand"
          defaultValue="edc:POLICY_EVALUATION_TIME"
          label={false}
          sx={{ display: "none" }}
        />
        <SelectInput
          source="permissions.time.operator"
          choices={[
            { id: "odrl:gt", name: "After" },
            { id: "odrl:lt", name: "Before" },
          ]}
          helperText="Select the operator"
          validate={[required()]}
        />
        <DateTimeInput
          source="permissions.time.rightOperand"
          helperText="Select the date"
          validate={[required()]}
        />
      </AccordionDetails>
    </Accordion>
  );
};

const LocationBasedPermission = () => {
  const [locationBasedPermission, setLocationBasedPermission] = useState(false);
  const onClick = () => {
    setLocationBasedPermission(!locationBasedPermission);
  };
  const {
    field: { value: operatorValue },
  } = useInput({ source: "permissions.location.operator" });
  const selectMultiple =
    operatorValue === "odrl:isAnyOf" ||
    operatorValue === "odrl:isAllOf" ||
    operatorValue === "odrl:isNoneOf";

  const renderInput = () => {
    if (selectMultiple) {
      return (
        <AutocompleteArrayInput
          source="permissions.location.rightOperand"
          choices={countriesList}
          label="Select the countries"
          helperText="Select the countries"
        />
      );
    }
    return (
      <AutocompleteInput
        source="permissions.location.rightOperand"
        choices={countriesList}
        label="Select the country"
        helperText="Select the country"
      />
    );
  };

  return (
    <Accordion
      expanded={locationBasedPermission}
      slotProps={{ transition: { unmountOnExit: true } }}
    >
      <AccordionSummary>
        <span>
          <Typography gutterBottom={false}>
            Enable Location Based Permission
          </Typography>
          <Typography variant="caption">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit.
          </Typography>
        </span>
        <span style={{ flex: 1 }} />
        <BooleanInput
          source="permissions.location.enable"
          onClick={onClick}
          sx={{ alignSelf: "flex-end" }}
          label={false}
        />
      </AccordionSummary>
      <AccordionDetails>
        <TextInput
          source="permissions.location.leftOperand"
          defaultValue="edc:location"
          label={false}
          sx={{ display: "none" }}
        />
        <SelectInput
          source="permissions.location.operator"
          choices={[
            { id: "odrl:eq", name: "Equals" },
            { id: "odrl:neq", name: "Not Equals" },
            { id: "odrl:isAnyOf", name: "Is Any Of" },
            { id: "odrl:isAllOf", name: "Is All Of" },
            { id: "odrl:isNoneOf", name: "Is None Of" },
          ]}
          helperText="Select the operator"
          validate={[required()]}
        />
        {renderInput()}
      </AccordionDetails>
    </Accordion>
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
  <List empty={false} hasCreate={true} exporter={false}>
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

export const PolicyShow = () => (
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
        <ArrayField source="policy.odrl:permission" label="Permissions">
          <Datagrid bulkActionButtons={false} rowClick={false} hover={false}>
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
              <Datagrid bulkActionButtons={false} rowClick={false} style={{ tableLayout: "fixed" }} hover={false}>
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
      </SimpleShowLayout>
    </SimpleShowLayout>
  </Show>
);

export const PolicyCreate = (props: any) => (
  <Create {...props}>
    <SimpleForm>
      <Grid container spacing={2}>
        <Grid item xs={3}>
          <SelectInput
            source="policy.@type"
            label="Type"
            choices={[{ id: "odrl:Set", name: "Set" }]}
            fullWidth
            defaultValue={"odrl:Set"}
            validate={[required()]}
            readOnly
          />
        </Grid>
        <Grid item xs={9}>
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
          <IdentityBasedPermission />
        </Grid>
        <Grid item xs={12}>
          <TimeBasedPermission />
        </Grid>
        <Grid item xs={12}>
          <LocationBasedPermission />
        </Grid>
      </Grid>
    </SimpleForm>
  </Create>
);
