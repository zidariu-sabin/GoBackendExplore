package store

// A file containing methods for pulling data through the TMDB api and populating a PostgresSQL database
import (
	// "GoBackendExploreMovieTracker/internal/store"
	"GoBackendExploreMovieTracker/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// passed (?,?) for as many pairs of data for columns as necessary
func buildSqlInsertPlaceholders(stmt, bindVars string, len int) string {
	bindVars += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(bindVars, len))
	return strings.TrimSuffix(stmt, ",")
}

// method for building psql params in relation to the amount of data following to be added in a table
func buildPostgresInsertPlaceholders(rowCount, colCount int) string {
	placeholders := make([]string, rowCount)
	arg := 1
	for i := 0; i < rowCount; i++ {
		row := make([]string, colCount)
		for j := 0; j < colCount; j++ {
			row[j] = fmt.Sprintf("$%d", arg)
			arg++
		}
		placeholders[i] = fmt.Sprintf("(%s)", strings.Join(row, ", "))
	}
	return strings.Join(placeholders, ", ")
}

// Store used for postgres database operations
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresSeedStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

type BodyGenre struct {
	Genres []Genre `json:"genres"`
}

type BodyMovie struct {
	Results []Movie `json:"results"`
}

func (pg *PostgresStore) PullGenres() {
	url := "https://api.themoviedb.org/3/genre/movie/list"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.TMDB_API_READ_ACCESS_TOKEN)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var body BodyGenre

	err := json.NewDecoder(res.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(body.Genres)
	vals := []interface{}{}
	sqlStr := "INSERT INTO genres(id, name) VALUES %s"

	for _, row := range body.Genres {
		vals = append(vals, row.ID, row.Name)
	}
	// fmt.Println(vals)
	// sqlStr = setupBindVars(sqlStr, "(?, ?)", len(body.Genres))
	sqlStr = fmt.Sprintf(
		"INSERT INTO genres(id, name) VALUES %s",
		buildPostgresInsertPlaceholders(len(body.Genres), 2),
	)

	// fmt.Println(sqlStr)
	_, err = pg.db.Exec(sqlStr, vals...)

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(sqlStr)
}

func (pg *PostgresStore) PullMovies() {

	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc&with_genres=28%20AND%%2012%20AND%2016"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.TMDB_API_READ_ACCESS_TOKEN)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var body BodyMovie

	err := json.NewDecoder(res.Body).Decode(&body)

	// fmt.Println(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	//
	vals := []interface{}{}

	for _, row := range body.Results {
		vals = append(vals, row.ID, row.Title, row.ReleaseDate, row.Overview, row.PosterPath, row.GenreIds)
	}
	// fmt.Println(vals)
	// sqlStr = setupBindVars(sqlStr, "(?, ?)", len(body.Genres))
	sqlStr := fmt.Sprintf(
		"INSERT INTO movies(id, title, release_date, overview, poster_path, genre_ids) VALUES %s",
		buildPostgresInsertPlaceholders(len(body.Results), 6),
	)

	// fmt.Println(sqlStr)
	_, err = pg.db.Exec(sqlStr, vals...)

	if err != nil {
		fmt.Println(err)
		return
	}
}
