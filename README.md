#  Portal

This repository contains the frontend and backend code for the SmartLivingNEXT portal. 

## How to run this image
```shell
docker run -it --rm --name portal registry.int.smartlivingnext.de/registry/public/portal:latest
```

## Configuration

The portal can be configured with the following environment variables:

| Name                          | Description                                                                 |
|-------------------------------|-----------------------------------------------------------------------------|
| ENABLE_DEVICE_API             | Enables device management through ThingsBoard.                             |
| ENABLE_FUSEKI_BACKEND         | Enables the creation of Fuseki datasets for user data management.          |
| ENABLE_USERS_API              | Enables user management functionalities through Keycloak.                  |
| EDC_URL                       | Base URL of the Eclipse Data Connector (EDC).                              |
| EDC_API_KEY                   | API key used to authenticate requests to the EDC.                          |
| FEDERATED_CATALOG_URL         | Base URL of the federated catalog for data discovery.                      |
| FEDERATED_CATALOG_API_KEY     | API key for accessing the federated catalog.                               |
| FUSEKI_USERNAME               | Username for authenticating with the Fuseki server.                       |
| FUSEKI_PASSWORD               | Password for authenticating with the Fuseki server.                       |
| THINGSBOARD_URL               | Base URL of the ThingsBoard instance for IoT management.                  |
| THINGSBOARD_EXCHANGE_TOKEN_URL| URL used for exchanging tokens with the ThingsBoard authentication service.|
| KEYCLOAK_CLIENT_ID            | Client ID for authenticating with Keycloak.                               |
| KEYCLOAK_CLIENT_SECRET        | Client secret for authenticating with Keycloak.                           |
| KEYCLOAK_REALM                | The realm in Keycloak that organizes and manages user identities.          |
| KEYCLOAK_URL                  | Base URL of the Keycloak server.                                           |

## Users

User identification is managed via the following custom headers:

| Header Name              | Description                          |
|--------------------------|--------------------------------------|
| X-User-Id               | The unique identifier of the user.  |
| X-User-Full-Name        | The full name of the user.          |
| X-User-Email            | The email address of the user.      |
| X-User-Tenant-Id        | The tenant ID associated with the user. |
| X-User-Customer-Id      | The customer ID associated with the user. |
| X-User-Roles            | A comma-separated list of user roles (admin,customer) |

If any headers are missing, the following environment variables provide default values:

| Environment Variable          | Default Value Description                |
|-------------------------------|------------------------------------------|
| DEFAULT_USER_ID              | Default value for the user ID.           |
| DEFAULT_USER_FULL_NAME       | Default value for the user's full name.  |
| DEFAULT_USER_EMAIL           | Default value for the user's email.      |
| DEFAULT_USER_TENANT_ID       | Default value for the tenant ID.         |
| DEFAULT_USER_CUSTOMER_ID     | Default value for the customer ID.       |
| DEFAULT_USER_ROLES           | Default value for the user roles (comma-separated). |

## Authors

Sebastian Alberternst <sebastian.alberternst@dfki.de>

## License

MIT 

