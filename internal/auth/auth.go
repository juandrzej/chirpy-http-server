package auth

func HashPassword(password string) (string, error)

func CheckPasswordHash(password, hash string) (bool, error)
