import { Card, CardContent, CardHeader } from "@mui/material";
import SparqlQuery from "../components/sparql_query";

export const SparqlPage = () => (
  <Card>
    <CardHeader title="SPARQL Endpoints" />
    <CardContent>
      <SparqlQuery />
    </CardContent>
  </Card>
);
