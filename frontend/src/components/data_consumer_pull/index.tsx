import {
  Datagrid,
  ArrayField,
  useShowController,
  TextField,
  Show,
  SimpleShowLayout,
} from "react-admin";
import PropTypes from "prop-types";
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";
import useTheme from "@mui/material/styles/useTheme";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { subDays } from "date-fns";
import PreviewIcon from "@mui/icons-material/Preview";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import { useEffect } from "react";

// everything related to visualization
const dateFormatter = (date: number): string =>
  new Date(date).toLocaleString("de-DE");

const NoVisualisation = () => {
  return (
    <Box textAlign="center">
      <PreviewIcon sx={{ width: "9em", height: "9em" }} color="primary" />
      <Typography variant="h6" paragraph>
        No Preview
      </Typography>
    </Box>
  );
};

const Visualization = ({ record }) => {
  if (typeof record.data !== "undefined" && record.data !== null) {
    if (record.data.constructor === Array && record.data.length > 0) {
      //    case array of json, telemetry
      if (record.data[0].ts) {
        return <TimeseriesVisualisation record={record} />;
      }
      // case array of json
      else {
        return <ArrayVisualisation record={record} />;
      }
    }
    // case json
    else if (record.data.constructor === Object) {
      return <SimpleVisualisation record={record} />;
    }
  }

  // case impossible to visualize

  return <NoVisualisation />;
};

Visualization.propTypes = {
  record: PropTypes.object,
};

const SimpleVisualisation = ({ record }) => {
  let labels = Object.keys(record.data);
  labels = labels.filter((label) => typeof record.data[label] !== "object");
  return (
    <SimpleShowLayout>
      {labels.map((label) => (
        <TextField source={"data." + label} key={label} />
      ))}
    </SimpleShowLayout>
  );
};

SimpleVisualisation.propTypes = {
  record: PropTypes.object,
};

const ArrayVisualisation = ({ record }) => {
  let labels = Object.keys(record.data[0]);
  labels = labels.filter((label) => typeof record.data[0][label] !== "object");
  return (
    <ArrayField source="data">
      <Datagrid rowClick={false} bulkActionButtons={false}>
        {labels.map((label) => (
          <TextField source={label} key={label} />
        ))}
      </Datagrid>
    </ArrayField>
  );
};

ArrayVisualisation.propTypes = {
  record: PropTypes.object,
};

const TimeseriesVisualisation = ({ record }) => {
  // const { record } = useShowController();
  let processedData = record.data;
  let labels = Object.keys(processedData[0]);
  // find another label
  let yLabel = labels.filter((item) => item !== "ts");
  processedData.sort((a, b) => a.ts - b.ts);
  // convert to time
  const transformedData = processedData.map((element) => {
    return { ...element, ts: new Date(element.ts).getTime() };
  });

  const theme = useTheme();
  return (
    <ResponsiveContainer width="85%" height={400}>
      <AreaChart data={transformedData}>
        <defs>
          <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
            <stop
              offset="5%"
              stopColor={theme.palette.primary.main}
              stopOpacity={0.8}
            />
            <stop
              offset="95%"
              stopColor={theme.palette.primary.main}
              stopOpacity={0}
            />
          </linearGradient>
        </defs>
        <XAxis
          dataKey="ts"
          name="Date"
          scale="time"
          domain={[subDays(new Date(), 1).getTime(), new Date().getTime()]}
          tickFormatter={dateFormatter}
        />
        <YAxis dataKey={yLabel[0]} name={yLabel[0]} />
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip labelFormatter={(label: any) => dateFormatter(label)} />
        <Area
          type="monotone"
          dataKey={yLabel[0]}
          stroke={theme.palette.primary.main}
          strokeWidth={2}
          fill="url(#colorUv)"
        />
      </AreaChart>
    </ResponsiveContainer>
  );
};

TimeseriesVisualisation.propTypes = {
  record: PropTypes.object,
};

//  for visualization of data
export const DataConsumerPullShow = () => {
  const { record } = useShowController();

  return (
    <Show>
      <SimpleShowLayout>
        <TextField source="id" />
        <Visualization record={record} />
      </SimpleShowLayout>
    </Show>
  );
};

// everything related to downloading of raw data
const Download = (record) => {
  // file object + define type of data from incoming header
  const file = new Blob([record.bodyString], { type: record.contentType });
  // anchor link
  const element = document.createElement("a");
  element.href = URL.createObjectURL(file);
  element.download = record.id + "_" + Date.now();
  // simulate link click
  document.body.appendChild(element); // Required for this to work in FireFox
  element.click();
};

const HandleDownload = (record) => {
  // download always when refresh
  useEffect(() => {
    if (record) {
      Download(record);
    }
  }, [record]);
};

// download show
export const RawDataDownloadShow = () => {
  const { record } = useShowController();
  HandleDownload(record);
  return (
    <Show>
      <SimpleShowLayout>
        <Box textAlign="center">
          <CheckCircleIcon
            sx={{ width: "9em", height: "9em" }}
            color="primary"
          />
          <Typography variant="h6" paragraph>
            Downloaded
          </Typography>
        </Box>
      </SimpleShowLayout>
    </Show>
  );
};

RawDataDownloadShow.propTypes = {
  record: PropTypes.object,
};
