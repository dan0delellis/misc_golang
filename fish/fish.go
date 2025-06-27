package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
)

func main() {

	tsv, _ := os.Open("fishies")
	scn := bufio.NewScanner(tsv)
	var lines []string

	for scn.Scan() {
		lines = append(lines, scn.Text())
	}
	var fishies []Fish
	for _, val := range lines {
		var temp Fish
		temp = parseFish(val)
		fishies = append(fishies, temp)
	}

	fmt.Println("made fishes, now opening db")
	db, err := sql.Open("mysql", "golang:donkeyboner@tcp(127.0.0.1:3306)/nh")
	fmt.Println("open db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("able to ping db")
	for _, val := range fishies {
		id := insertFish(db, val)
		if id == 0 {
			continue
		}

		fmt.Println("setting time data for")
		fmt.Println(val)
		setTemporalData(id, val, db)
	}

}

func setTemporalData(id int64, f Fish, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	st, err := tx.Prepare("insert into creatures_time (time_id, creature_id) values (?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, val := range f.Times {
		if val {
			_, err = st.Exec(key+1, id)
			if err != nil {
				fmt.Println(err)
				tx.Rollback()
				return
			}
		}
	}

	st, err = tx.Prepare("insert into creatures_months (month_id, creature_id) values (?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, val := range f.Months {
		if val {
			_, err = st.Exec(key+1, id)
			if err != nil {
				fmt.Println(err)
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
}

func insertFish(db *sql.DB, f Fish) (id int64) {
	st, err := db.Prepare("Insert into creatures (type, name, price, location, location_sub, size, fin) values (?, ?, ?, ?, ?, ?, ?)")
	res, err := st.Exec(f.Type, f.Name, f.Price, f.Location.Main, f.Location.Sub, f.Size.Size, f.Size.Fin)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	id, err = res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return 0
	}
	return

}

func parseFish(s string) (f Fish) {
	p := strings.Split(s, "\t")
	f.Type = p[0]
	f.Name = p[1]
	f.Price = parseCost(p[2])
	f.Location = getLocation(p[3])
	f.Size = getSize(p[4])
	f.Times = parseTimes(p[5])
	f.Months = parseMonths(p[6])
	return
}

type Fish struct {
	Type     string
	Name     string
	Price    int64
	Location Location
	Size     Shadow
	Times    [24]bool
	Months   [12]bool
}

type Location struct {
	Main string
	Sub  string
}

type Shadow struct {
	Size int64
	Fin  bool
}

func parseCost(s string) (c int64) {
	s = strings.Replace(s, ",", "", -1)
	c, _ = strconv.ParseInt(s, 10, 64)
	return
}

func getLocation(code string) (location Location) {
	loc := make(map[string]Location)
	loc["0"] = Location{"river", ""}
	loc["0.1"] = Location{"river", "mouth"}
	loc["0.2"] = Location{"river", "cliff"}
	loc["1"] = Location{"lake", ""}
	loc["2"] = Location{"sea", ""}
	loc["2.1"] = Location{"sea", "pier"}
	loc["2.2"] = Location{"sea", "rain"}
	loc["-1"] = Location{"water", ""}
	loc["-2"] = Location{"water", "fresh"}

	location = loc[code]
	return
}

func getSize(code string) (size Shadow) {
	s := strings.Split(code, ".")
	size.Size, _ = strconv.ParseInt(s[0], 10, 64)
	if len(s) == 2 {
		size.Fin = true
	} else {
		size.Fin = false
	}
	return
}

func parseTimes(s string) (t [24]bool) {
	p := strings.Split(s, ",")

	for _, val := range p {
		x, _ := strconv.ParseInt(val, 10, 8)
		t[x] = true
	}

	return
}

func parseMonths(s string) (m [12]bool) {
	p := strings.Split(s, ",")

	for key, val := range p {
		if val == "TRUE" {
			m[key] = true
		}
	}
	return
}
