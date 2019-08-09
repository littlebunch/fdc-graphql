# fdc-graphql
A GraphQL prototype for Food Data Central datasets. It is built atop the [fdc-api](https://github.com/littlebunch/fdc-api) libraries and assumes you have a database build with [fdc-ingest](https://github.com/littlebunch/fdc-ingest).  The server schema is built with [graphql-go](https://github.com/graphql-go/graphql) with a playground provided by [gqlgen](https://github.com/99designs/gqlgen/handler).    
## Building    
You need to have go version 1.11 or higher installed.     
### Step 1: Clone this repo
Clone the repo into a path *other* than your $GOPATH:
```
git clone git@github.com:littlebunch/fdc-graphql.git
```
### Step 2
cd to the fdc-graphql path and configure for your datastore by creating and populating a config.yml:


```
couchdb:   
  url:  localhost   
  bucket: gnutdata   //default  bucket    
  fts: fd_food  // default full-text index   
  user: <your_user>    
  pwd: <your_password>    

```
or set these vars in the environment:
```
COUCHBASE_URL=localhost   
COUCHBASE_BUCKET=gnutdata   
COUCHBASE_FTSINDEX=fd_food   
COUCHBASE_USER=user_name   
COUCHBASE_PWD=user_password   
```
### Step 3: Start the server.
```
go run main.go -c config.yml -p 8000 -r graphql
```
### Usage
A playground is available at http://localhost:8000/graphql/.  Some queries to run include:

Query for an food by FDC id:
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356427"){fdcId,description,company,ingredients,servings{Description,Nutrientbasis,Servingamount}}}'
```
A list of foods:
```
curl -g 'http://localhost:8000/graphql?query={foods{fdcId,description,company,ingredients,servings{Nutrientbasis, Description,Servingamount}}}'
```
Nutrient data for a food:
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356425"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}}nutrientdata(fdcid:"356425",nutid:0){Nutrient,Nutrientno,Value}}'
```
Data for an individual nutrient for a food:
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356425"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}} nutrientdata(fdcid:\"356425\",nutid:203){Nutrient,Nutrientno,Value}}'
```
Get a list nutrients from the database:
```
curl -g 'http://localhost:8000/graphql?query={nutrients{Nutrientno,Name,Unit}}'
```

