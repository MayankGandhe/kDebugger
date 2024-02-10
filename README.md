# kDebugger
This repo contains code for a docker image that can be used to debug issues on your K8s cluster 


# Running the application

To run the application, follow these steps:

- Create a .env File: Create a file named .env in the root directory of the application. Paste the following content     into the file and adjust the details accordingly:

PORT=3000
VARIABLE_ONE=valueOne
VARIABLE_TWO=valueTwo

MONGO_HOST=localhost
MONGO_USER=root
MONGO_PASSWORD=
MONGO_PORT=27017
MONGO_DATABASE=myDatabase

MYSQL_USER=root
MYSQL_PASSWORD=
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DATABASE=mysql

- Ensure to replace placeholders like valueOne, valueTwo, myDatabase, etc., with actual values relevant to your setup.

- Open Terminal: Open a terminal or command prompt.

- Run the Application: Use the following command to run the application:

    go run main.go
This command will start the application, and it will use the configurations provided in the .env file.

Access the Application: Once the application is running, you can access it by navigating to http://localhost:3000 in your web browser.

That's it! You have successfully run the application with the provided configurations.



# API Endpoints

| Method | Endpoint                                        | Description                       |
|--------|-------------------------------------------------|-----------------------------------|
| POST   | localhost:3000/                                 | Get headers                       |
| POST   | localhost:3000/env                              | Get all env variables             |
| POST   | localhost:3000/env-from-dotenv                  | Get .env file variables           |
| GET    | localhost:3000/env/mongo                        | Get filtered environment variables|
| POST   | localhost:3000/setup-and-check-mysql-connection | Set vars &  check mysql connection|
| GET    | localhost:3000/timeout/200                      | Check timeout                     |
