package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "time"
	"strconv"

    _ "github.com/mattn/go-sqlite3"
)

func getXboardUsers(db *sql.DB) ([]XboardUser, error) {
    query := `SELECT id, email, uuid, u, d, transfer_enable, group_id, plan_id 
              FROM v2_user WHERE is_admin = 0`
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []XboardUser
    for rows.Next() {
        var user XboardUser
        err := rows.Scan(&user.ID, &user.Email, &user.UUID, &user.U, &user.D, 
            &user.TransferEnable, &user.GroupID, &user.PlanID)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, rows.Err()
}

func syncUsers(xboardDB, xuiDB *sql.DB) error {

    users, err := getXboardUsers(xboardDB)
    if err != nil {
        return fmt.Errorf("ERR: Read xboard resp fail: %w", err)
    }
	
	log.Println("debug: getid num=",len(users))

    inbound, err := getInbound(xuiDB)
    if err != nil {
        return err
    }

    if inbound == nil {
        return nil
    }

    var settings XuiSettings
    if err := json.Unmarshal([]byte(inbound.Settings), &settings); err != nil {
        return err
    }
	
    existingClients := make(map[string]bool)
    for _, client := range settings.Clients {
        existingClients[client.Email] = true
    }
	
    addedCount := 0
    for _, user := range users {
        if !existingClients[user.Email] {
            client := XuiClient{
                ID:         user.UUID,
                Security:   "",
                Password:   user.UUID,
                Flow:       xuiFlow,
                Email:      "xboard_"+ user.UUID+strconv.Itoa(xuiID),
                LimitIp:    0,
                TotalGB:    user.TransferEnable / (1024 * 1024 * 1024),
                ExpiryTime: 0,
                Enable:     true,
                TgId:       0,
                SubId:      user.UUID,
                Comment:    "",
                Reset:      0,
                CreatedAt:  time.Now().UnixMilli(),
                UpdatedAt:  time.Now().UnixMilli(),
            }

            settings.Clients = append(settings.Clients, client)
            
            if err := addClientTraffic(xuiDB, inbound.ID, user.Email); err != nil {
                log.Print(err)
            }
            
            addedCount++
        }
    }

    if addedCount > 0 {
        newSettings, err := json.MarshalIndent(settings, "", "  ")
        if err != nil {
            return err
        }

        if err := updateInboundSettings(xuiDB, inbound.ID, string(newSettings)); err != nil {
            return err
        }
    } else {
    }

    return nil
}

func updateXboardUserTraffic(db *sql.DB, email string, up, down int64) error {
    query := `UPDATE v2_user SET u = ?, d = ?, updated_at = ? WHERE email = ?`
    
    now := time.Now().Unix()
    result, err := db.Exec(query, up, down, now, email)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("ERR: request to xboard fail")
    }

    return nil
}