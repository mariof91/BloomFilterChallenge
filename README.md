# BloomFilterChallenge

#challenge:

START >>> Flo Bloom Filter Challenge
Build a simple REST API backed by a configurable Bloom Filter.
Reference link: https://hur.st/bloomfilter
1. Do not use any filter library, create all structures and logic.
2. You can use any HTTP middleware ( built in, mux, etc )
3. Provide all source code and compiled code for Linux/Windows/MacOS
4. Once app is launched, using port 8080 must be able to execute the HTTP commands with application/json content below
5. set-name = any url safe string
6. item-name = any url safe string
   POST /sets/set-name
   {
   "type": "bloom-filter",
   "config": {
   "size": 256,
   "functions: 3
   }
   }
   PUT /sets/{set-name}/{item-name}
   GET /sets/{set-name}/items/{item-name}
   {
   "exists" : true|false
   }
   GET /sets/{set-name}/stats
   {
   "size": 256,
   "functions": 3,
   "count": 74,
   "falsePositiveProbability": 0.194979311
   }
   END <<< Flo Bloom Filter Challenge

## Install
docker build -t bloomfilter .
docker run -d -p 8080:8080 bloomfilter

you can test the app from postman using the BloomFilter.postman_collection.json file

