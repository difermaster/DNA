# DNA

1. En la siguiente URL https://dna-mutant-2021.herokuapp.com/mutant será posible consumir a través de HTTP POST el servicio que evalúa el ADN e identifica si es o no mutante; de esta manera si es Mutante devuelve: HTTP 200-OK o si es humano devuelve 403-Forbidden

La siguiente es la entrada JSON:

POST → /mutant/
{
“dna”:["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"]
}


2. En la siguiente URL https://dna-mutant-2021.herokuapp.com/stats será posible consumir a través de HTTP GET el servicio que devuelve un JSON con las estadísticas de las verificaciones de ADN realizadas.

