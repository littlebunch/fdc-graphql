# fdc-graphql
A GraphQL prototype for Food Data Central datasets. It is built atop the [fdc-api](https://github.com/littlebunch/fdc-api) libraries and assumes you have a database built with [fdc-ingest](https://github.com/littlebunch/fdc-ingest).  The server schema is built with [graphql-go](https://github.com/graphql-go/graphql) with a playground provided by [gqlgen](https://github.com/99designs/gqlgen/handler).  A demo playground using the FDC [Branded Food Products](https://fdc.nal.usda.gov/data-documentation.html) dataset is available at https://go.littlebunch.com/graphql/.
## Building    
If you are not using the available Docker image at littlebunch/fdcgql, then you need to have [Go](https://golang.org/dl/) version 1.11 or higher installed to build and run the server.     
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
go run main.go -c config.yml -p 8000
```
Or, run from docker.io (Note: You will need docker installed. You will also need to pass in the Couchbase configuration as environment variables described above. The easiest way to do this is in a file, e.g. 'docker.env', of which a sample is provided in the repo's docker path.) :
```
docker run --rm -it -p 8000:8000 --env-file=./docker.env littlebunch/fdcgql
```
    
### Usage
Some queries to run from the [playground](https://go.littlebunch.com/graphql/) include:

Query for a food by FDC id:
```
query  {
   food(id:"356425"){
        fdcId
        foodDescription
        dataSource
        servingSizes{
            nutrientBasis
            servingUnit
            value
        }
        foodGroup{
          description
        }
    } 
}
  
```
```
curl -g 'https://go.littlebunch.com/graphql?query={food(id:"356427"){fdcId,foodDescription,company,ingredients,servingSizes{nutrientBasis,servingUnit,value},foodGroup{description}}}'
```
```
curl -XPOST -H "Content-type:application-json" https://go.littlebunch.com/graphql -d '{"query":"{food(id:\"356427\"){fdcId,foodDescription,company,ingredients,servingSizes{servingUnit,nutrientBasis,value},foodGroup{description}}}"}'
```
Browse foods:
```
{
   foodsBrowse(browse:{page:0,max:50,sort:"foodDescription"}){
        fdcId
        foodDescription
        company
        ingredients
        servingSizes{
           nutrientBasis
           servingUnit
           value
         }
    }
}
```
```
curl -g 'https://go.littlebunch.com/graphql?query={foodsBrowse(browse:{page:0,max:50,sort:"foodDescription"}){fdcId,foodDescription,company,ingredients,servingSizes{nutrientBasis,servingUnit,value}}}' 
```
```
curl -XPOST -H "Content-type:application/json" https://go.littlebunch.com/graphql -d '{"query":"{foodsBrowse(browse:{page:0,max:50,sort:\"foodDescription\"}){fdcId,foodDescription,company,ingredients,servingSizes{nutrientBasis, servingUnit,value}}}"}'
```
A list of foods given a list of FDC id's:
```
{
   foods(fdcids:["344604","344605","344606"]){
        fdcId
        foodDescription
        company
        ingredients
    }
}
```
Search for foods:
```
{
   foodsSearch(search:{terms:"broccoli rabe",type:"PHRASE",field:"ingredients"}){
        fdcId
        foodDescription
        company
        ingredients
    }
}
```
```
curl -g 'https://go.littlebunch.com/graphql?query={foodsSearch(search:{terms:"broccoli rabe",type:"PHRASE",field:"ingredients"}){fdcId,foodDescription,company,ingredients}}'    
```      
Nutrient data for a food:    
```
{
   food(id:"356425"){
        fdcId
        foodDescription
        dataSource
        servingSizes{
            nutrientBasis
            servingUnit
            value
         }
    }
    nutrientdata(fdcids:["356425"],nutids:[]){
        nutrient
        nutrientno
        value
    }
}
```
```
curl -g 'https://go.littlebunch.com/graphql?query={food(id:"356425"){fdcId,foodDescription,dataSource,servingSizes{nutrientBasis,servingUnit,value}}nutrientdata(fdcids:["356425"],nutids:[]){nutrient,nutrientno,value}}'
```
Nutrient data for an individual nutrient in a food:
```
{
   food(id:"356425"){
        fdcId
        foodDescription
        dataSource
        servingSizes{
            nutrientBasis
            servingUnit
            value
        }
    } 
    nutrientdata(fdcids:["356425"],nutids:[203,204]){
        nutrient
        nutrientno
        value
    }
}
```
```
curl -g 'https://go.littlebunch.com/graphql?query={food(id:"356425"){fdcId,foodDescription,dataSource,servingSizes{nutrientBasis,servingUnit,value}}nutrientdata(fdcids:["356425"],nutids:[203,204]){nutrient,nutrientno,value}}'
```
```
curl -XPOST -H "Content-type:application/json" https://go.littlebunch.com/graphql -d '{"query":"{food(id:"356425"){fdcId,foodDescription,dataSource,servingSizes{nutrientBasis,servingUnit,value}}nutrientdata(fdcids:["356425"],nutids:[203,204]){nutrient,nutrientno,value}}"}'
```
Get a list nutrients from the database:
```
{
  nutrients{
     nutrientno
     name
     unit
   }
}
```
```
curl -g 'https://go.littlebunch.com/graphql?query={nutrients{nutrientno,name,unit}}'
```
```
curl -XPOST -H "Content-type:application/json" https://go.littlebunch.com/graphql -d '{"query":"{nutrients{nutrientno,name,unit}}"}'
```

