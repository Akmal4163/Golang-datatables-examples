package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Barang struct {
	ID 				int 	`json:"id"`
	KodeBarang 		string	`json:"kode_barang"`
	NamaBarang 		string	`json:"nama_barang"`
	Satuan 			string	`json:"satuan"`
	HargaSatuan 	float32	`json:"harga_satuan"`
	Jumlah 			float32	`json:"jumlah"`
	Keterangan 		string	`json:"keterangan"`
}

type DataTables struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            []Barang   `json:"data"`
}

type DatatableHandlers struct {
	DB *sql.DB
}

func (h *DatatableHandlers) GetData(w http.ResponseWriter, r *http.Request) {
	draw, _ := strconv.Atoi(r.URL.Query().Get("draw"))
	start, _ := strconv.Atoi(r.URL.Query().Get("start"))
	length, _ := strconv.Atoi(r.URL.Query().Get("length"))
	searchQuery := r.URL.Query().Get("search[value]")
	orderColumnIndex := r.URL.Query().Get("order[0][column]")
	orderColumnDir := r.URL.Query().Get("order[0][dir]")

	if length <= 0 {
		length = 10
	}

	columnMap := map[string]string {
		"0": "id",
		"1": "kode_barang",
		"2": "nama_barang",
		"3": "satuan",
		"4": "harga_satuan",
		"5": "jumlah",
		"6": "keterangan",
	}

	orderByCol, validOrder := columnMap[orderColumnIndex]
	if !validOrder {
		orderByCol = "id"
	}

	if orderColumnDir != "asc" && orderColumnDir != "desc" {
		orderColumnDir = "asc"
	}

	var recordsTotal int
	err := h.DB.QueryRow("SELECT COUNT(id) FROM master_barang").Scan(&recordsTotal)
	if err != nil {
		http.Error(w, "Error:", http.StatusInternalServerError)
		return
	}

	querySelect := fmt.Sprintf(`
		SELECT id, kode_barang, nama_barang, satuan, harga_satuan, jumlah, keterangan 
		FROM master_barang 
		WHERE kode_barang LIKE ? OR nama_barang LIKE ? OR satuan LIKE ? OR harga_satuan LIKE ?
		OR jumlah LIKE ? OR keterangan LIKE ? ORDER BY %s %s
		LIMIT ? OFFSET ?;`, orderByCol, orderColumnDir)

	searchPattern := "%" + searchQuery + "%"
	rows, err := h.DB.Query(querySelect, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, length, start)
	if err != nil {
		http.Error(w, "Error:", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	listBarang := []Barang{}
	for rows.Next() {
		var b Barang
		err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.Satuan, &b.HargaSatuan, &b.Jumlah, &b.Keterangan)
		if err != nil {
			http.Error(w, "Error:", http.StatusInternalServerError)
			return
		}
		listBarang = append(listBarang, b)
	}

	
	queryCountFilter := `SELECT COUNT(id) FROM master_barang 
		WHERE kode_barang LIKE ? OR nama_barang LIKE ? OR satuan LIKE ? OR harga_satuan LIKE ?
		OR jumlah LIKE ? OR keterangan LIKE ?;`
	
	var recordsFiltered int
	err = h.DB.QueryRow(queryCountFilter, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&recordsFiltered)
	if err != nil {
		http.Error(w, "Error:", http.StatusInternalServerError)
		return
	}

	response := DataTables{
		Draw: draw,
		RecordsTotal: recordsTotal,
		RecordsFiltered: recordsFiltered,
		Data: listBarang,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}