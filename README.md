# test-assigment
##Get authors:
```
curl -X GET localhost:8080/authors/
```
##Create author:
```
curl -X POST localhost:8080/authors/ -H "Accept: application/json" -d '{"FirstName":"<firstname>","LastName":"<lastname>"}'
```
##Delete author:
```
curl -X DELETE localhost:8080/authors/<author id>
```
##Create movie:
```
curl -X POST localhost:8080/movies/ -H "Accept: application/json" -d '{"MovieName":"<title>","MovieYear":<year>, "AuthorID":"<author id>"}'
```
##Get movies:
```
curl -X GET localhost:8080/movies/
```
##Get movie by id:
```
curl -X GET localhost:8080/movies/<id>
```
##Delete movie:
```
curl -X DELETE localhost:8080/movies/<id>
```
##Get movie by author:
```
curl -X GET localhost:8080/movies/authors/<author id>
```
##Restart app:
```
curl -X GET localhost:8080/triggerpanic/
```
