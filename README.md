# catalogsvc

## Getting Started

These instructions will allow you to run catalog service

## Requirements

Go (golang) : 1.11.2

mongodb as docker container

zipkin as docker container (optional)

## Instructions

1. Clone this repository 

2. You will notice the following directory structure

``` 
├── db.go
├── go.mod
├── go.sum
├── images
│   ├── catsocks_1.jpg
│   ├── cross_1.jpeg
│   ├── product2.jpg
│   ├── puma_1.jpeg
│   ├── slide1.jpg
│   ├── weave1.jpg
│   └── youtube_1.jpeg
├── main.go
├── products.json
├── README.md
└── service.go

```

3. Set GOPATH appropriately as per the documentation - https://github.com/golang/go/wiki/SettingGOPATH

4. Build the go application from the root of the folder

``` go build -o bin/catalog ```

5. Run a mongodb docker container

```sudo docker run -d -p 27017:27017 --name mgo -e MONGO_INITDB_ROOT_USERNAME=mongoadmin -e MONGO_INITDB_ROOT_PASSWORD=secret mongo```

6. Execute this command to import the ```products.json``` file 

   ```sudo docker cp products.json {mongodb_container_id}:/```


7. **Login into the mongodb container**

    
    ```sudo docker exec -it {mongodb_container_id} bash```

8. Import the products file into the database 
    
   ```mongoimport --db catalog --collection products --file products.json -u mongoadmin -p secret --authenticationDatabase=admin```

9. Export CATALOG_HOST/CATALOG_PORT (port and ip) as ENV variable. You may choose any used port as per your environment setup.
    
    ```export CATALOG_HOST=0.0.0.0```
    ```export CATALOG_PORT=:8082```

10. Also, export ENV variables related to the database

    ```
    export CATALOG_DB_USERNAME=mongoadmin
    export CATALOG_DB_PASSWORD=secret
    export CATALOG_DB_HOST=0.0.0.0
    ```

11. Run the catalog service

```./bin/catalog```


### Additional Info
