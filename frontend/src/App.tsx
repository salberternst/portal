import {
  Admin,
  Resource,
  ListGuesser,
  EditGuesser,
  ShowGuesser,
  List,
  Datagrid,
  DateField,
  TextField,
  Show,
  SimpleShowLayout,
} from "react-admin";
import dataSource from './data-source';

const ThingList = () => (
  <List empty={false} hasCreate={true} exporter={false}> 
      <Datagrid bulkActionButtons={false} rowClick="show">
          <TextField source="id.id" title="Id"/>
          <TextField source="name" />
          <TextField source="label" />
          {/* <DatagridActionsCol label="Actions">
              <InboxRowActions />
          </DatagridActionsColumn> */}
      </Datagrid>
  </List>
);

const ThingShow = () => (
  <Show>
      <SimpleShowLayout>
          
      </SimpleShowLayout>
  </Show>
)

export const App = () => (
  <Admin dataProvider={dataSource}>
    <Resource
      name="thing"
      list={ThingList}
      show={ShowGuesser}
      edit={EditGuesser}
    />
    {/* <Resource
      name="thingDescription"
      list={ListGuesser}
      edit={EditGuesser}
      show={ShowGuesser}
    /> */}
  </Admin>
);
