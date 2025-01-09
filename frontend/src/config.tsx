declare global {
  // eslint-disable-next-line no-unused-vars
  interface Window {
    config: {
      federatedCatalogUrl: string;
      showDevices: boolean;
      showThingDescriptions: boolean;
      showQuery: boolean;
      showCustomers: boolean;
      // edc settings
      showAssets: boolean;
      showPolicies: boolean;
      showContractDefinitions: boolean;
      showCatalog: boolean;
      showFederatedCatalog: boolean;
      showContractAgreements: boolean;
      showTransferProcesses: boolean;
      // edc settings
      showThingsboard: boolean;
      showKeycloak: boolean;
      // portal settings
      useInGeoLogo: boolean;
    };
  }
}

export {};
