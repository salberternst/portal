import { Card, CardContent, CardHeader } from "@mui/material";
import SparqlQuery from "../components/sparql_query";
import { Title } from "react-admin";

export const SparqlPage = () => (
  <Card>
    <Title title="Query" />
    <CardHeader title="SPARQL Endpoints" />
    <CardContent>
      <SparqlQuery />
    </CardContent>
  </Card>
);
