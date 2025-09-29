package main

import (
  "database/sql"
  "html/template"
  "log"
  "net/http"
  "os"
  "strings"

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
var detailsTmpl = template.Must(template.New("details").Parse(detailsTemplate))

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
  http.HandleFunc("/tambah", tambahHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/delete/", hapusHandler)
  http.HandleFunc("/details/", detailsHandler)
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

// tambah mahasiswa
func tambahHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    t := template.Must(template.New("form").Parse(formTemplate))
    t.Execute(w, nil)
    return
  }

  if err := r.ParseForm(); err != nil { 
    http.Error(w, "Form error", 400)
    return
  }

  nama := r.FormValue("nama")
  npm := r.FormValue("npm")
  kelas := r.FormValue("kelas")
  minat := r.FormValue("minat")

  _, err := db.Exec("INSERT INTO mahasiswa (Nama, NPM, Kelas, Minat) VALUES (?, ?, ?, ?)", nama, npm, kelas, minat)
  if err != nil {
    http.Error(w, "Insert failed: "+err.Error(), 500)
    return
  }

  http.Redirect(w, r, "/", http.StatusSeeOther) 
}

// edit mahasiswa
func editHandler(w http.ResponseWriter, r *http.Request) {
  id := strings.TrimPrefix(r.URL.Path, "/edit/")

  if r.Method == http.MethodGet {
    row := db.QueryRow("SELECT ID, Nama, NPM, Kelas, Minat FROM mahasiswa WHERE ID = ?", id)
    var m Mahasiswa

    err := row.Scan(&m.ID, &m.Nama, &m.NPM, &m.Kelas, &m.Minat)
    if err != nil {
      http.Error(w, "Mahasiswa not found", 404)
      return
    }

    t := template.Must(template.New("form").Parse(formTemplate))
    t.Execute(w, m)
    return
  }

  if err := r.ParseForm(); err != nil { 
    http.Error(w, "Form error", 400)
    return
  }

  nama := r.FormValue("nama")
  npm := r.FormValue("npm")
  kelas := r.FormValue("kelas")
  minat := r.FormValue("minat")

  _, err := db.Exec("UPDATE mahasiswa SET Nama=?, NPM=?, Kelas=?, Minat=? WHERE ID=?", nama, npm, kelas, minat, id)
  if err != nil {
    http.Error(w, "Update failed: "+err.Error(), 500)
    return
  }

  http.Redirect(w, r, "/", http.StatusSeeOther) 
}

// hapus mahasiswa
func hapusHandler(w http.ResponseWriter, r *http.Request) {
  id := strings.TrimPrefix(r.URL.Path, "/delete/")

  _, err := db.Exec("DELETE FROM mahasiswa WHERE ID = ?", id)
  if err != nil {
    http.Error(w, "Delete failed: "+err.Error(), 500)
    return
  }

  http.Redirect(w, r, "/", http.StatusSeeOther)
}

// detail mahasiswa
func detailsHandler(w http.ResponseWriter, r *http.Request) {
  id := strings.TrimPrefix(r.URL.Path, "/details/")

  row := db.QueryRow("SELECT ID, Nama, NPM, Kelas, Minat FROM mahasiswa WHERE ID = ?", id)
  var m Mahasiswa
  err := row.Scan(&m.ID, &m.Nama, &m.NPM, &m.Kelas, &m.Minat)
  if err != nil {
    http.Error(w, "Mahasiswa not found: "+err.Error(), http.StatusNotFound)
    return
  }

  err = detailsTmpl.Execute(w, m)
  if err != nil {
    http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
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

// hal edit, tambah
const formTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{if .ID}}Edit{{else}}Tambah{{end}} Mahasiswa</title>
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
      max-width: 600px;
      margin: auto;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h2 { margin-bottom: 20px; }
    label { display: block; margin-top: 10px; font-weight: bold; }
    input[type="text"] {
      width: 580px;
      padding: 8px;
      margin-top: 5px;
      border: 1px solid #ccc;
      border-radius: 4px;
    }
    .btn {
      padding: 8px 14px;
      margin-top: 15px;
      border-radius: 4px;
      text-decoration: none;
      color: white;
      font-size: 14px;
      display: inline-block;
    }
    .submit { background: #0d6efd; border: none; cursor: pointer; }
    .cancel { background: #6c757d; }
  </style>
</head>
<body>
  <div class="container">
    <h2>{{if .ID}}Edit{{else}}Tambah{{end}} Mahasiswa</h2>
    <form method="POST" action="{{if .ID}}/edit/{{.ID}}{{else}}/tambah{{end}}">
      <label>Nama:</label>
      <input type="text" name="nama" value="{{.Nama}}" required>

      <label>NPM:</label>
      <input type="text" name="npm" value="{{.NPM}}" required>

      <label>Kelas:</label>
      <input type="text" name="kelas" value="{{.Kelas}}" required>

      <label>Minat:</label>
      <input type="text" name="minat" value="{{.Minat}}" required>

      <button type="submit" class="btn submit">{{if .ID}}Update{{else}}Simpan{{end}}</button>
      <a href="/" class="btn cancel">Batal</a>
    </form>
  </div>
</body>
</html>
`

// hal detail
const detailsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Detail Mahasiswa</title>
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
      max-width: 600px;
      margin: auto;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h2 { margin-bottom: 20px; }
    .info { margin-bottom: 10px; }
    .label { font-weight: bold; display: inline-block; width: 80px; }
    .btn {
      padding: 8px 14px;
      margin-top: 15px;
      border-radius: 4px;
      text-decoration: none;
      color: white;
      font-size: 14px;
      display: inline-block;
    }
    .back { background: #6c757d; }
    .edit { background: #ffc107; color: #000; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Detail Mahasiswa</h2>
    <div class="info"><span class="label">ID:</span> {{.ID}}</div>
    <div class="info"><span class="label">Nama:</span> {{.Nama}}</div>
    <div class="info"><span class="label">NPM:</span> {{.NPM}}</div>
    <div class="info"><span class="label">Kelas:</span> {{.Kelas}}</div>
    <div class="info"><span class="label">Minat:</span> {{.Minat}}</div>

    <a href="/" class="btn back">Kembali</a>
    <a href="/edit/{{.ID}}" class="btn edit">Edit</a>
  </div>
</body>
</html>
`