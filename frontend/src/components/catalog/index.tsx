import { useState } from "react";
import { Show, useRecordContext, useStore } from "react-admin";
import MuiTextField from "@mui/material/TextField";
import MuiButton from "@mui/material/Button";
import Alert from "@mui/material/Alert";
import { CatalogList } from "../catalog_list";
import LinkIcon from "@mui/icons-material/Link";
import LinkOffIcon from "@mui/icons-material/LinkOff";

const CatalogShow = () => {
  const record = useRecordContext();
  return <CatalogList record={record?.datasets} />;
};

const CatalogConnect = () => {
  const [counterPartyAddress, setCounterPartyAddress] = useStore(
    "counterPartyAddress",
    null
  );
  const [inputValue, setInputValue] = useState("/api/dsp");
  const connect = () => {
    setCounterPartyAddress(inputValue);
  };
  const CatalogConnectButton = () => {
    if (counterPartyAddress) {
      return (
        <MuiButton
          onClick={() => setCounterPartyAddress(null)}
          variant="text"
          startIcon={<LinkOffIcon />}
        >
          Disconnect
        </MuiButton>
      );
    } else {
      return (
        <MuiButton onClick={connect} variant="text" startIcon={<LinkIcon />}>
          Connect
        </MuiButton>
      );
    }
  };
  return (
    <MuiTextField
      label="EDC Address"
      value={inputValue}
      disabled={counterPartyAddress !== null}
      onChange={(e) => setInputValue(e.target.value)}
      InputProps={{
        endAdornment: <CatalogConnectButton />,
      }}
      fullWidth
    />
  );
};

export const Catalog = () => {
  const [counterPartyAddress] = useStore("counterPartyAddress");
  const [error, setError] = useState();

  const onError = (e) => {
    setError(e);
  };

  const onSuccess = () => {
    setError(null);
  };

  return (
    <>
      <CatalogConnect />
      {error && (
        <Alert severity="error">
          Unable to fetch catalog {counterPartyAddress}
        </Alert>
      )}
      {counterPartyAddress && (
        <Show
          resource="catalog"
          id={counterPartyAddress}
          queryOptions={{
            onError,
            onSuccess,
          }}
        >
          <CatalogShow />
        </Show>
      )}
    </>
  );
};
