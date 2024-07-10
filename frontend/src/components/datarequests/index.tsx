import {
    TextField,
    Show,
    SimpleShowLayout,
    useShowController,
    Button,
} from "react-admin";
import Typography from '@mui/material/Typography';

export const DataRequestShow = () => {
    const { record } = useShowController();

    return (
        <Show>
            <SimpleShowLayout>
                <TextField source="id" />
                <TextField source="contractId" label="Contract Id" />
                <TextField source="endpoint" />
                <TextField source="authKey" label="Authentication Key" />
                <TextField source="authCode" label="Authentication Code" />
                <Typography variant="h6">Copy Curl Command:</Typography>
                <Typography variant="body1">
                    curl -X GET -H "{record?.authKey}: {record?.authCode}" {record?.endpoint}
                </Typography>
            </SimpleShowLayout>
        </Show>
    );
}