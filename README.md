# Shakespearean Pokemon API

The Shakespearean Pokemon API is a simple yet scalable Go API capable of returning the description of your favourite 
Pokemon in Shakespeare's style.

## User stories:
A client can request pokemon description given:
- the pokemon name

e.g: 
        
        If the client want's to know Charizard's Shakespearian description: 
            name: charizard
        which will return back
            description: Charizard flies 'round the sky in search of powerful opponents. 't breathes fire of such most wondrous heat yond 't melts aught. However, 't nev'r turns its fiery breath on any opponent weaker than itself.

## Assumptions
* The requested pokemon exists in the PokeAPI ([pokeapi.co pokemon list](https://pokeapi.co/api/v2/pokemon/?limit=1000))
* Only the most updated description of the pokemon in english is required.
* Both external API (defined in the Dependant APIs section) are up and running.

## How to run: 
### Prerequisites: 
- Go 1.13
- Docker 19.03.5
- DockerHub access to pull image

### How run the api:
To pull the docker image from DockerHub, please run the below command
```
docker run --rm -p 8080:8080 sam195/shakespeareanpokemonapi:latest
```
The API will run on ```http://localhost:8080/``` the usage is described in the "Usage" section of this README.

If that fails, open terminal in the root of the project and run the below command:

```
go run main.go
```

## Usage

This can be done using multiple tools such as Postman, Curl or a simple browser, requirements mentioned httpie, 
hence, once the service is up and serving on port 8080, the below request can be executed.
```
http http://localhost:8080/pokemon/charizard
```

### Get pokemon's description in Shakespear's style from API

**Definition**

`GET http://localhost:8080/pokemon/<PokemonName>`

PokemonName: contains the name of the pokemon

**Response**

- `200 OK` on success

```json
{
	"name":"the name of the requested pokemon",
	"description":"the pokemon_domain of the requested pokemon in Shakespear's style"
}
```
e.g.

- For: 
`GET http://localhost:8080/pokemon/charizard`

```
{
    "name":"charizard",
    "description":"Charizard flies 'round the sky in search of powerful opponents. 't breathes fire of such most 
wondrous heat yond 't melts aught. However, 't nev'r turns its fiery breath on any opponent weaker than itself."
}
```

- `400 Bad Request` if any of the fields are invalid, or connection to external api can not be established
- `404 Not Found` if the pokemon was not found, this error will be returned
- `429 Too Many Requests` if the request limit specified in the Dependent APIs section below is hit
- `500 Internal Server Error` if any of the two external API return something that is not expected

## How to test
The project contains both Unit and Integration tests, below are steps to run them

### Unit tests
The unit tests mocks the call to the external API in order to only check the functionality of the 
Shakespearean Pokemon API

To run the unit tests:

Open terminal in the root of the project

```
go test ./... -short
```

### Integration tests
The integration tests test calls between the Shakespearean Pokemon API and the external APIs, for example, this is done 
to check whether the external application has changed request formats

Open terminal in the root of the project, If these are run, the limit of translation requests (defined in the Dependant 
APIs section below) will be reached.

```
go test ./... -run Integration
```

### Dependant APIs
The proposed solution is dependant on two main APIs as mentioned below:
- [PokeAPI](https://pokeapi.co/docs/v2) : which has a 300 requests limit per resource per IP address
- [Shakespeare Translator](https://funtranslations.com/api/shakespeare) : which has a 5 requests limit per hour

### Future work
* Assuming the translation method won't change in time, caching should be used to reduce the number of request to the
FunTranslationAPI, so that we can achieve more than 5 requests per hour, the cache should also be stored in a persistent
location using perhaps docker volumes or even a proper db.
* Logging and metrics to analyse what are the most common requests we get might also be useful in the long term, to see
how the service behaves and what the system's bottlenecks are. It was decided these features were out of scope for this
task.
* In case this system will be used in a bigger micro service architecture, it will be required to have a config file
which makes easy the ports settings and provider's url in case they change. There should also be a liveness and readiness
endpoints to check whether our service is live and readiness to check whether we are ready to process requests. Readiness
will also check whether the dependent systems are live and ready to process too, as if not, our system will be not ready.
