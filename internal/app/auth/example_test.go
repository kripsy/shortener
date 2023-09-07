//nolint:testpackage
package auth

import (
	"fmt"
	"log"
)

//nolint:testableexamples
func ExampleBuildJWTString() {
	// Suppose we have user id 123
	userID := 123

	// Create a JWT string for a given user
	tokenString, err := BuildJWTString(userID)
	if err != nil {
		log.Fatalf("Failed to build JWT string: %v", err)
	}

	fmt.Println("Generated JWT token:", tokenString)
}

//nolint:testableexamples,nolintlint
func ExampleGetUserID() {
	// To demonstrate, first create a token using the BuildJWTString function
	userID := 12345
	tokenString, err := BuildJWTString(userID)
	if err != nil {
		fmt.Println("Error creating JWT:", err)

		return
	}

	// Now use this token to get the userID
	extractedUserID, err := GetUserID(tokenString)
	if err != nil {
		fmt.Println("Error extracting userID:", err)

		return
	}

	fmt.Println("Extracted UserID:", extractedUserID)
	// Output: Extracted UserID: 12345
}

//nolint:testableexamples,nolintlint
func ExampleGetExpires() {
	// To demonstrate, first create a token using the BuildJWTString function
	userID := 12345
	tokenString, err := BuildJWTString(userID)
	if err != nil {
		fmt.Println("Error creating JWT:", err)

		return
	}

	// Now use this token to get the expiration time
	expirationTime, err := GetExpires(tokenString)
	if err != nil {
		fmt.Println("Error extracting expiration time:", err)

		return
	}

	fmt.Println("Token expires at:", expirationTime)
}

func ExampleIsTokenValid() {
	// To demonstrate, first create a token using the BuildJWTString function
	userID := 12345
	tokenString, err := BuildJWTString(userID)
	if err != nil {
		fmt.Println("Error creating JWT:", err)

		return
	}

	// Now let's check if this token is valid
	isValid, err := IsTokenValid(tokenString)
	if err != nil {
		fmt.Println("Error checking token validity:", err)

		return
	}

	if isValid {
		fmt.Println("Token is valid!")
	} else {
		fmt.Println("Token is not valid!")
	}
	// Output: Token is valid!
}
