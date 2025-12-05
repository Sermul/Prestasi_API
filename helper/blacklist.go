package helper

import "sync"

// blacklist token hanya disimpan di memori (RAM)
var blacklist = struct {
    tokens map[string]bool
    sync.RWMutex
}{tokens: make(map[string]bool)}

// Tambahkan token ke blacklist
func AddToBlacklist(token string) {
    blacklist.Lock()
    defer blacklist.Unlock()
    blacklist.tokens[token] = true
}

// Cek apakah token sudah diblacklist (tidak boleh dipakai lagi)
func IsBlacklisted(token string) bool {
    blacklist.RLock()
    defer blacklist.RUnlock()
    return blacklist.tokens[token]
}
