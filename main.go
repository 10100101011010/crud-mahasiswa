package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)


type Mahasiswa struct { 
	ID    int
	Nama  string
	NPM   string
	Kelas string
	Minat string
}

var db *sql.DB
var tmpl = template.Must(template.New("main").Parse(htmlTemplate))

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASS") + "@tcp(" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME")

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Connected to MySQL!")
	http.HandleFunc("/", handler)
	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// list mahasiswa
func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT ID, Nama, NPM, Kelas, Minat FROM mahasiswa")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	var mahasiswaList []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		err := rows.Scan(&m.ID, &m.Nama, &m.NPM, &m.Kelas, &m.Minat)
		if err != nil {
			http.Error(w, "Data error: "+err.Error(), 500)
			return
		}
		mahasiswaList = append(mahasiswaList, m)
	}

	err = tmpl.Execute(w, mahasiswaList)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), 500)
	}
}

// hal utama
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Data Minat Mahasiswa</title>
<style>
    body {
      font-family: Arial, sans-serif;
      background: #f8f9fa;
      margin: 0;
      padding: 20px;
    }
    .container {
      background: #fff;
      padding: 20px;
      border-radius: 6px;
      max-width: 900px;
      margin: auto;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h2 { margin-bottom: 15px; }
    table {
      width: 100%;
      border-collapse: collapse;
    }
    thead {
      border-bottom: 2px solid #343a40;
      color: #343a40;
      text-align: left;
    }
    th, td {
      padding: 10px;
    }
    tbody tr {
      border-bottom: 1px solid #dee2e6;
    }
    tbody tr:last-child {
      border-bottom: none;
    }
    .btn {
      padding: 5px 10px;
      margin-right: 4px;
      border-radius: 4px;
      text-decoration: none;
      color: white;
      font-size: 14px;
    }
    .tambah { background: #0d6efd; margin-bottom: 10px; display: inline-block; }
    .details { background: #28a745; }
    .edit { background: #ffc107; color: #000; }
    .delete { background: #dc3545; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Data Minat Mahasiswa</h2>
    <a href="/tambah" class="btn tambah">Tambah</a>

    <table>
      <thead>
        <tr>
          <th>Id</th>
          <th>Nama</th>
          <th>NPM</th>
          <th>Kelas</th>
          <th>Minat</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {{range .}}
        <tr>
          <td>{{.ID}}</td>
          <td>{{.Nama}}</td>
          <td>{{.NPM}}</td>
          <td>{{.Kelas}}</td>
          <td>{{.Minat}}</td>
          <td>
            <a href="/details/{{.ID}}" class="btn details">Details</a>
            <a href="/edit/{{.ID}}" class="btn edit">Edit</a>
            <a href="/delete/{{.ID}}" class="btn delete" onclick="return confirm('Are you sure you want to delete this?')">Delete</a>
          </td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>
</body>
</html>
`
