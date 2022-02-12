##GO Lang based scaffold for creating a project with multiple microservices. 

This scaffold implements the **hexagonal architecture** with microservices. It clearly separates the business/domain logic from the transport, DB and external app services layer. The project structure is a mono-repo design and hence aims at accommodating multiple microservices in a single project.

The structure of the project is as follows:

root - /

service - /service - The parent directory for a microservice. If needed another folder can be created at the same level as shown bellow

service2 - /service2 - The second microservice.

Within a microservice, the following folders are created:

cmd - /service/cmd - This contains the main.go with a main function that imports and invokes the code from the other folders. This is the entry point of the microservice.

entrypoints - /service/cmd/entrypoints - This contains the entrypoint for the entities within the service. This function initializes the repo/data services, http routes and endpoints of a single entity.

db - /service/db - This contains the logic that implements the database connection and CRUD operations needed for the domain/business logic

domain - /service/domain - This contains the domain/business logic of the microservice.

external-services - /service/external-services - This contains the external app/data services that are needed for the microservice. ex. And 3rd party API's etc.

sql - /service/sql - This contains the SQL DDL that are needed for creating the tables in the database.

test - /service/test - This contains the test files that are needed for testing the microservice.

transport - /service/transport - This contains the transport layer logic - the http/gRPC routes and endpoints of the microservice.

endpoints - /service/transport/endpoints - This contains the endpoints for the http/gRPC routes of the microservice.

http - /service/transport/http - This contains the http routes of the microservice.