package main

import (
    "log"
    "github.com/mediocregopher/radix.v2/redis"
    "strconv"
)

type Album struct {
    Likes   int
    Price  float64
    Title  string
    Artist string
}

func main() {
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        log.Println(err.Error())
    }
    defer conn.Close()
    
    //if err := conn.Cmd("HMSET", "album:2", "title", "Black album", "artist", "Metallica", "price", 4.98, "likes", 8).Err; err != nil {
    //    log.Fatalln(err)
    //}

    //if artist, err := conn.Cmd("HGET", "album:2", "artist").Str(); err != nil {
    //    log.Fatalln(err)
    //} else {
    //    log.Printf("success: [%s]\n", artist)
    //}

    result, err := conn.Cmd("HGETALL", "album:2").Map()
    if err != nil {
        log.Fatalln(err)
    }
    ab, err := populateAlbum(result)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println(ab)
}

func populateAlbum(reply map[string]string) (*Album, error) {
    var err error
    ab := new(Album)
    ab.Title = reply["title"]
    ab.Artist = reply["artist"]
    // We need to use the strconv package to convert the 'price' value from a
    // string to a float64 before assigning it.
    ab.Price, err = strconv.ParseFloat(reply["price"], 64)
    if err != nil {
        return nil, err
    }
    // Similarly, we need to convert the 'likes' value from a string to an
    // integer.
    ab.Likes, err = strconv.Atoi(reply["likes"])
    if err != nil {
        return nil, err
    }
    return ab, nil
}