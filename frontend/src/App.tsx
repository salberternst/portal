import {
  Admin,
  Resource,
  CustomRoutes,
  Layout,
  Menu,
  AppBar,
  useGetIdentity,
} from "react-admin";
import { Route } from "react-router-dom";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import DeviceHub from "@mui/icons-material/DeviceHub";
import QueryStatsIcon from "@mui/icons-material/QueryStats";
import DashboardIcon from "@mui/icons-material/Dashboard";
import PeopleOutlineIcon from "@mui/icons-material/PeopleOutline";
import PolicyIcon from "@mui/icons-material/Policy";
import GroupsIcon from "@mui/icons-material/Groups";
import GavelIcon from "@mui/icons-material/Gavel";
import AutoStoriesIcon from "@mui/icons-material/AutoStories";
import AutoModeIcon from "@mui/icons-material/AutoMode";
import ShieldIcon from "@mui/icons-material/Shield";
import DescriptionIcon from "@mui/icons-material/Description";
import WebStoriesIcon from "@mui/icons-material/WebStories";
import InventoryIcon from "@mui/icons-material/Inventory";
import Toolbar from "@mui/material/Toolbar";
import { useLocation } from "react-router-dom";
import dataSource from "./data-source";
import authProvider from "./auth-provider";
import { SparqlPage } from "./custom_pages/sparql";
import Thingsboard from "./components/thingsboard";
import {
  ThingDescriptionCreate,
  ThingDescriptionEdit,
  ThingDescriptionList,
  ThingDescriptionShow,
} from "./components/thing_description";
import { AssetCreate, AssetShow, AssetsList } from "./components/assets";
import {
  CustomerCreate,
  CustomerShow,
  CustomersList,
  CustomerUpdate,
} from "./components/customers";
import { UserCreate, UserShow } from "./components/users";
import { PoliciesList, PolicyCreate, PolicyShow } from "./components/policies";
import {
  ContractDefinitionCreate,
  ContractDefinitionShow,
  ContractDefinitionsList,
} from "./components/contract_definitions";
import { Catalog } from "./components/catalog";
import { FederatedCatalogList } from "./components/federated_catalog";
import {
  ContractNegotationCreate,
  ContractNegotationShow,
  ContractNegotiationTerminate,
} from "./components/contract_negotiations";
import {
  ContractAgreementShow,
  ContractAgreementsList,
} from "./components/contract_agreements";
import {
  TransferProcessesCreate,
  TransferProcessesList,
  TransferProcessesShow,
  TransferProcessTerminate,
} from "./components/transfer_processes";
import {
  DeviceCreate,
  DeviceEdit,
  DeviceShow,
  DevicesList,
} from "./components/devices";
import Keycloak from "./components/keycloak";
import { DataRequestShow } from "./components/datarequests";
import {
  DataConsumerPullShow,
  RawDataDownloadShow,
} from "./components/data_consumer_pull";
import { darkTheme, theme } from "./theme";
import SvgIcon from "@mui/material/SvgIcon";
import SmartLivingNextIcon from "./assets/borlabs-cookie-icon-dynamic.svg?react";
import Box from "@mui/material/Box";

const CustomUserMenu = () => {
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }

  return (
    <>
      <Typography variant="button">{identity?.email}</Typography>
    </>
  );
};

const CustomAppBar = () => {
  const logo = window.config.logo || "SmartLivingNEXT-bw-1.svg";

  return (
    <AppBar userMenu={<CustomUserMenu />}>
      <Toolbar>
        <Box
          sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
        >
          <SvgIcon
            color="primary"
            fontSize="large"
            style={{
              marginRight: 10,
              animation: "spin 20s linear infinite",
            }}
          >
            <SmartLivingNextIcon />
          </SvgIcon>
          <Box
            sx={{
              width: 200,
              height: 35,
              mixBlendMode: "difference",
            }}
          >
            <img src={`src/assets/${logo}`} alt="Logo" style={{ width: '100%', height: '100%' }} />
          </Box>
        </Box>
      </Toolbar>
      <span style={{ flex: 1 }} />
    </AppBar>
  );
};

const CustomMenu = () => {
  const { isLoading, identity } = useGetIdentity();
  if (isLoading) {
    return null;
  }

  const isAdmin = identity?.roles.includes("admin");

  return (
    <Menu dense={false} sx={{ pt: 1 }}>
      {window.config.showDevices && <Menu.ResourceItem name="devices" />}
      {window.config.showThingDescriptions && (
        <Menu.ResourceItem name="thingDescriptions" />
      )}
      {isAdmin && window.config.showCustomers && (
        <Menu.ResourceItem name="customers" />
      )}
      {window.config.showQuery && (
        <Menu.Item
          to="/sparql"
          primaryText="Query"
          leftIcon={<QueryStatsIcon />}
        />
      )}
      <Divider />
      {window.config.showAssets && <Menu.ResourceItem name="assets" />}
      {window.config.showPolicies && <Menu.ResourceItem name="policies" />}
      {window.config.showContractDefinitions && (
        <Menu.ResourceItem name="contractdefinitions" />
      )}
      {isAdmin && window.config.showCatalog && (
        <Menu.Item
          to="/catalog"
          primaryText="Catalog"
          leftIcon={<AutoStoriesIcon />}
        />
      )}
      {isAdmin && window.config.showFederatedCatalog && (
        <Menu.ResourceItem name="federatedcatalog" />
      )}
      {window.config.showContractAgreements && (
        <Menu.ResourceItem name="contractagreements" />
      )}
      {window.config.showTransferProcesses && (
        <Menu.ResourceItem name="transferprocesses" />
      )}
      <Divider />
      {window.config.showThingsboard && (
        <Menu.Item
          to="/thingsboard"
          primaryText="Thingsboard"
          leftIcon={<DashboardIcon />}
        />
      )}
      {isAdmin && window.config.showKeycloak && (
        <Menu.Item
          to="/keycloak"
          primaryText="Keycloak"
          leftIcon={<ShieldIcon />}
        />
      )}
    </Menu>
  );
};

const CustomLayout = (props: any) => {
  const location = useLocation();

  if (
    location.pathname === "/thingsboard" ||
    location.pathname === "/keycloak"
  ) {
    return (
      <Layout menu={CustomMenu} appBar={CustomAppBar} sx={{ pt: 2 }}>
        {props.children}
      </Layout>
    );
  }

  return (
    <Layout menu={CustomMenu} appBar={CustomAppBar} sx={{ pt: 2 }}>
      <Container maxWidth="lg">{props.children}</Container>
    </Layout>
  );
};

export const App = () => (
  <Admin
    loginPage={false}
    layout={CustomLayout}
    dataProvider={dataSource}
    authProvider={authProvider}
    theme={theme}
    darkTheme={darkTheme}
    disableTelemetry={true}
  >
    <CustomRoutes>
      {window.config.showQuery && (
        <Route path="/sparql" element={<SparqlPage />} />
      )}
      {window.config.showThingsboard && (
        <Route path="/thingsboard" element={<Thingsboard />} />
      )}
      {window.config.showCatalog && (
        <Route path="/catalog" element={<Catalog />} />
      )}
      {window.config.showKeycloak && (
        <Route path="/keycloak" element={<Keycloak />} />
      )}
    </CustomRoutes>
    {window.config.showThingDescriptions && (
      <Resource
        name="thingDescriptions"
        options={{ label: "Thing Descriptions" }}
        icon={DescriptionIcon}
        list={ThingDescriptionList}
        show={ThingDescriptionShow}
        create={ThingDescriptionCreate}
        edit={ThingDescriptionEdit}
      />
    )}
    {window.config.showDevices && (
      <Resource
        name="devices"
        options={{ label: "Devices" }}
        icon={DeviceHub}
        list={DevicesList}
        show={DeviceShow}
        create={DeviceCreate}
        edit={DeviceEdit}
      />
    )}
    {window.config.showAssets && (
      <Resource
        name="assets"
        options={{ label: "Assets" }}
        icon={InventoryIcon}
        list={AssetsList}
        show={AssetShow}
        create={AssetCreate}
      />
    )}
    {window.config.showCustomers && (
      <Resource
        name="customers"
        options={{ label: "Customers" }}
        icon={GroupsIcon}
        list={CustomersList}
        show={CustomerShow}
        create={CustomerCreate}
        edit={CustomerUpdate}
      />
    )}
    {window.config.showCustomers && (
      <Resource
        name="users"
        options={{ label: "Users" }}
        icon={PeopleOutlineIcon}
        show={UserShow}
        create={UserCreate}
      />
    )}
    {window.config.showPolicies && (
      <Resource
        name="policies"
        options={{ label: "Policies" }}
        icon={PolicyIcon}
        list={PoliciesList}
        show={PolicyShow}
        create={PolicyCreate}
      />
    )}
    {window.config.showContractDefinitions && (
      <Resource
        name="contractdefinitions"
        options={{ label: "Contract Definitions" }}
        icon={GavelIcon}
        list={ContractDefinitionsList}
        show={ContractDefinitionShow}
        create={ContractDefinitionCreate}
      />
    )}
    <Resource
      name="contractnegotiations"
      options={{ label: "Contract Negotiations" }}
      create={ContractNegotationCreate}
      show={ContractNegotationShow}
    >
      <Route path=":id/terminate" element={<ContractNegotiationTerminate />} />
    </Resource>
    {window.config.showContractAgreements && (
      <Resource
        name="contractagreements"
        options={{ label: "Contract Agreements" }}
        list={ContractAgreementsList}
        show={ContractAgreementShow}
      />
    )}
    {window.config.showTransferProcesses && (
      <Resource
        name="transferprocesses"
        icon={AutoModeIcon}
        options={{ label: "Transfer Processes" }}
        list={TransferProcessesList}
        show={TransferProcessesShow}
        create={TransferProcessesCreate}
      >
        <Route path=":id/terminate" element={<TransferProcessTerminate />} />
      </Resource>
    )}
    <Resource
      name="datarequests"
      options={{ label: "Data Requests" }}
      show={DataRequestShow}
    />
    <Resource
      name="dataconsumerpull"
      options={{ label: "Data Consumer Pull" }}
      show={DataConsumerPullShow}
    />
    <Resource
      name="rawdataconsumerpull"
      options={{ label: "Data Download" }}
      show={RawDataDownloadShow}
    ></Resource>
    {window.config.showFederatedCatalog && (
      <Resource
        name="federatedcatalog"
        icon={WebStoriesIcon}
        options={{ label: "Federated Catalog" }}
        list={FederatedCatalogList}
      />
    )}
  </Admin>
);
