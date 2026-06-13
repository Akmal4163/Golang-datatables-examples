#!/usr/bin/python3

import csv
import sqlite3

conn = sqlite3.connect('data_barang.db')
cursor = conn.cursor()

cursor.execute('''
    CREATE TABLE IF NOT EXISTS master_barang (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        kode_barang VARCHAR(200) NOT NULL UNIQUE,
        nama_barang TEXT NOT NULL,
        satuan VARCHAR(200) NOT NULL,
        harga_satuan REAL NOT NULL,
        jumlah INTEGER NOT NULL,
        keterangan TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
''')

with open('data_barang.csv', 'r', encoding='utf-8') as file:
    reader = csv.DictReader(file)
    data_to_insert = [(row['kode_barang'], row['nama_barang'], row['satuan'], row['harga_satuan'], row['jumlah'], row['keterangan']) for row in reader]
    query_insert = "INSERT INTO master_barang (kode_barang, nama_barang, satuan, harga_satuan, jumlah, keterangan) VALUES (?, ?, ?, ?, ?, ?)"
    
    cursor.executemany(query_insert, data_to_insert)

conn.commit()
conn.close()

print("Success")
