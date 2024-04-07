package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"context"
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:	 "localhost:6379",
        Password: "password123",
        DB:		  0, // use default DB
    })

    // Crea un contesto per il comando PING
    ctx := context.Background()

    // Invia un comando PING al server Redis
    pong, err := client.Ping(ctx).Result()
    if err != nil {
        fmt.Println("Errore nella connessione a Redis:", err)
        return
    }

    // Se la risposta è "PONG", la connessione è stata stabilita con successo
    if pong == "PONG" {
        fmt.Println("Connessione a Redis stabilita con successo")
    } else {
        fmt.Println("Risposta inaspettata dal server Redis:", pong)
    }
}
