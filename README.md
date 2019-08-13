# fdc-graphql
A GraphQL prototype for Food Data Central datasets. It is built atop the [fdc-api](https://github.com/littlebunch/fdc-api) libraries and assumes you have a database built with [fdc-ingest](https://github.com/littlebunch/fdc-ingest).  The server schema is built with [graphql-go](https://github.com/graphql-go/graphql) with a playground provided by [gqlgen](https://github.com/99designs/gqlgen/handler).    
## Building    
You need to have go version 1.11 or higher installed.     
### Step 1: Clone this repo
Clone the repo into a path *other* than your $GOPATH:
```
git clone git@github.com:littlebunch/fdc-graphql.git
```
### Step 2 Configure the database
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
go run main.go schema.go -c config.yml -p 8000 -r graphql
```
Or, run from docker.io (Note: You will need docker installed. You will also need to pass in the Couchbase configuration as environment variables described above. The easiest way to do this is in a file, e.g. 'docker.env', of which a sample is provided in the repo's docker path.) :
```
docker run --rm -it -p 8000:8000 --env-file=./docker.env littlebunch/fdcgql
```
    
### Usage
A playground is available at http://localhost:8000/graphql/.  Some queries to run include:

Query for a food by FDC id:
```
{
  food(id:"356427"){
      fdcId 
      description
      company
      ingredients
      servings{
        Description
        Nutrientbasis
        Servingamount
     }
   }
}
```
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356427"){fdcId,description,company,ingredients,servings{Description,Nutrientbasis,Servingamount}}}'
```
```
curl -XPOST -H "Content-type:application-json" http://localhost:8000/graphql -d '{"query":"{food(id:\"356427\"){fdcId,description,company,ingredients,servings{Description,Nutrientbasis,Servingamount}}}"}'
```
A list of foods:
```
{
   foods{
        fdcId
        description
        company
        ingredients
        servings{
           Nutrientbasis
           Description
           Servingamount
         }
    }
}
```
```
curl -g 'http://localhost:8000/graphql?query={foods{fdcId,description,company,ingredients,servings{Nutrientbasis, Description,Servingamount}}}'
```
```
curl -XPOST -H "Content-type:application/json" http://localhost:8000/graphql -d '{"query":"{foods{fdcId,description,company,ingredients,servings{Nutrientbasis, Description,Servingamount}}}"}'
```
Nutrient data for a food:
```
{
   food(id:"356425"){
        fdcId
        description
        dataSource
        servings{
            Nutrientbasis
            Description
            Servingamount
         }
    }
    nutrientdata(fdcid:"356425",nutid:[]){
        Nutrient
        Nutrientno
        Value
    }
}
```
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356425"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}}nutrientdata(fdcid:"356425",nutids:[]){nutrient,nutrientno,value}}'
```
```
curl -XPOST -H "Content-type:application/json" http://localhost:8000/graphql -d '{"query":"{food(id:"356425"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}}nutrientdata(fdcid:"356425",nutids:[]){nutrient,nutrientno,value}}"}'
```
Nutrient data for an individual nutrient in a food:
```
{
   food(id:"356425"){
        fdcId
        description
        dataSource
        servings{
            Nutrientbasis
            Description
            Servingamount
        }
    } 
    nutrientdata(fdcid:"356425",nutids:[203,204]){
        nutrient
        nutrientno
        value
    }
}
```
```
curl -g 'http://localhost:8000/graphql?query={food(id:"356425"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}}nutrientdata(fdcid:"356425",nutids:[203,204]){nutrient,nutrientno,value}}'
```
```
curl -XPOST -H "Content-type:application/json" http://localhost:8000/graphql -d '{"query":"{food(id:\"356425\"){fdcId,description,dataSource,servings{Nutrientbasis,Description,Servingamount}}nutrientdata(fdcid:"356425",nutids:[203,204]){nutrient,nutrientno,value}}"}'
```
Get a list nutrients from the database:
```
{
  nutrients{
     Nutrientno
     Name
     Unit
   }
}
```
```
curl -g 'http://localhost:8000/graphql?query={nutrients{Nutrientno,Name,Unit}}'
```
```
curl -XPOST -H "Content-type:application/json" http://localhost:8000/graphql -d '{"query":"{nutrients{Nutrientno,Name,Unit}}"}'
```

