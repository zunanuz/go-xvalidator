package main
import "fmt"
func calculateChecksum(id string) int {
    sum := 0
    for i := 0; i < 12; i++ {
        digit := int(id[i] - '0')
        sum += digit * (13 - i)
    }
    checksum := (11 - (sum % 11)) % 11
    if checksum == 10 {
        checksum = 0
    }
    return checksum
}
func main() {
    ids := []string{"110370016611", "310112345678", "101010101010"}
    for _, id := range ids {
        checksum := calculateChecksum(id)
        fullID := id + fmt.Sprintf("%d", checksum)
        fmt.Printf("%s checksum=%d full=%s\n", id, checksum, fullID)
    }
}
