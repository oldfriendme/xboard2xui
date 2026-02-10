package main

import (
    "database/sql"
    "log"
	"strconv"

    _ "github.com/mattn/go-sqlite3"
)

func syncTraffic(xboardDB, xuiDB *sql.DB) error {
    traffics, err := getXuiClientTraffics(xuiDB)
    if err != nil {
        return err
    }
	
	log.Println("debug: getid num=",len(traffics))

    updatedCount := 0
    for _, traffic := range traffics {
        if err := updateXboardUserTraffic(xboardDB, traffic.Email, traffic.Up, traffic.Down); err != nil {
            log.Print(err)
            continue
        }
        updatedCount++
    }
	log.Println("debug: getid num=",updatedCount)
    return nil
}

func getInbound(db *sql.DB) (*XuiInbound, error) {
    query := `SELECT id, port, protocol, settings, remark FROM inbounds LIMIT `+strconv.Itoa(xuiID)
    
    var inbound XuiInbound
	
    err := db.QueryRow(query).Scan(&inbound.ID, &inbound.Port, &inbound.Protocol, 
        &inbound.Settings, &inbound.Remark)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &inbound, nil
}

func updateInboundSettings(db *sql.DB, inboundID int, settings string) error {
    query := `UPDATE inbounds SET settings = ? WHERE id = ?`
    _, err := db.Exec(query, settings, inboundID)
    return err
}

func addClientTraffic(db *sql.DB, inboundID int, email string) error {
    var count int
    checkQuery := `SELECT COUNT(*) FROM client_traffics WHERE inbound_id = ? AND email = ?`
    err := db.QueryRow(checkQuery, inboundID, email).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return nil
    }
    query := `INSERT INTO client_traffics (inbound_id, enable, email, up, down, all_time, expiry_time, total, reset, last_online)
              VALUES (?, 1, ?, 0, 0, 0, 0, 0, 0, 0)`
    
    _, err = db.Exec(query, inboundID, email)
    return err
}

func getXuiClientTraffics(db *sql.DB) ([]XuiClientTraffic, error) {
    query := `SELECT inbound_id, email, up, down, enable FROM client_traffics`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var traffics []XuiClientTraffic
    for rows.Next() {
        var traffic XuiClientTraffic
        var enable int
        err := rows.Scan(&traffic.InboundID, &traffic.Email, &traffic.Up, &traffic.Down, &enable)
        if err != nil {
            return nil, err
        }
        traffic.Enable = enable == 1
        traffics = append(traffics, traffic)
    }

    return traffics, rows.Err()
}