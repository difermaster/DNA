package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

type Report struct {
	NoMutant int `json:"count_mutant_dna"`
	Mutant   int `json:"count_human_dna"`
}

const limit int = 2
const noOfChars int = 256

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido al API de validación del ADN")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/mutant", isMutant).Methods("POST")
	router.HandleFunc("/stats", stats).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func stats(w http.ResponseWriter, r *http.Request) {
	report, _ := FindAll()
	json.NewEncoder(w).Encode(report)
}

func isMutant(w http.ResponseWriter, r *http.Request) { //dna []string) bool {
	//dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"} //Mutante
	dna := []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"} //No-Mutante

	var count int = 0
	var len int = len(dna)
	var i int = 0
	var j int = 0
	var k int = 0
	var l int = len - 1
	var m int = l

	for count < limit && i < len {
		var row string = dna[i]
		var col string = ""
		var obLR string = ""
		var obRL string = ""
		count += Search(ToCharArray(row), "Horizontal")

		for count < limit && j < len {
			col += (ToCharArray(dna[j]))[i]

			if len-4 >= k && m >= j+k {
				obLR += (ToCharArray(dna[j]))[j+k]
				obRL += (ToCharArray(dna[j]))[l-j]
			}

			j++
		}

		if count < limit {
			count += Search(ToCharArray(col), "Vertical")
		}

		if count < limit && obLR != "" {
			count += Search(ToCharArray(obLR), "Oblicuo Izquierda-Derecha")
		}

		if count < limit && obRL != "" {
			count += Search(ToCharArray(obRL), "Oblicuo Derecha-Izquierda")
		}

		j = 0
		i++
		k++
		l--
	}

	w.Header().Set("Content-Type", "application/json")

	var reply bool = count >= limit

	result := Create(reply, dna)
	fmt.Println("Result: ", result)

	if reply {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func Search(txt []string, orientation string) int {
	patterns := []string{"AAAA", "CCCC", "GGGG", "TTTT"}
	var count int = 0
	var p int = 0

	for {
		var pat []string = ToCharArray(patterns[p])
		var m int = len(pat)
		var n int = len(txt)

		badchar := BadCharHeuristic(pat, m)

		var s int = 0

		for count < limit && s <= (n-m) {
			var j int = m - 1

			for j >= 0 && pat[j] == txt[s+j] {
				j--
			}

			if j < 0 {
				fmt.Println("Los patrones se producen en el turno = ", s, ", orientacion = "+orientation+", Combinacion = "+patterns[p])
				count++

				if s+m < n {
					s += m - badchar[int(txt[s+m][0])]
				} else {
					s += 1
				}
			} else {
				s += Max(1, j-badchar[int(txt[s+j][0])])
			}
		}

		p++

		if p >= len(patterns) || count > limit {
			break
		}
	}

	return count
}

func BadCharHeuristic(str []string, size int) [noOfChars]int {
	var badchar [noOfChars]int
	var i int

	for i = 0; i < noOfChars; i++ {
		badchar[i] = -1
	}

	for i = 0; i < size; i++ {
		badchar[int(str[i][0])] = i
	}

	return badchar
}

func ToCharArray(str string) []string {
	var ar []string

	for _, r := range str {
		ar = append(ar, string(r))
	}

	return ar
}

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func GetDB() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", ".db/dna.db")
	return
}

func FindAll() ([]Report, error) {
	db, err := GetDB()

	if err != nil {
		return nil, err
	} else {
		rows, err2 := db.Query("select (select count(*) from report where ismutant = true) as Mutant, (select count(*) from report where ismutant = false) as NoMutant")
		if err2 != nil {
			return nil, err
		} else {
			var reports []Report

			for rows.Next() {
				var report Report
				rows.Scan(&report.Mutant, &report.NoMutant)
				reports = append(reports, report)
				break
			}

			return reports, nil
		}
	}
}

func Create(isMutant bool, DNA []string) bool {
	db, err := GetDB()

	out, _ := json.Marshal(DNA)
	result, err := db.Exec("insert into report(isMutant, DNA) values (?,?)", isMutant, string(out))
	if err != nil {
		return false
	}

	rowsAffected, err2 := result.RowsAffected()
	if err2 != nil {
		return false
	}

	return rowsAffected > 0
}
