package main

import (
	"fmt"
	"log"
	"math/rand"
	
	"time"

	"github.com/joho/godotenv"
	"notes-api/config"
	"notes-api/models"
)

// Sample data for seeding
var (
	sampleUsers = []models.User{
		{Name: "John Doe", Email: "john@example.com", Password: "password123"},
		{Name: "Jane Smith", Email: "jane@example.com", Password: "password123"},
		{Name: "Bob Johnson", Email: "bob@example.com", Password: "password123"},
		{Name: "Alice Brown", Email: "alice@example.com", Password: "password123"},
		{Name: "Charlie Wilson", Email: "charlie@example.com", Password: "password123"},
	}

	sampleNoteTitles = []string{
		"Meeting Notes",
		"Project Ideas",
		"Shopping List",
		"Book Recommendations",
		"Travel Plans",
		"Recipe Collection",
		"Workout Routine",
		"Learning Goals",
		"Daily Reflections",
		"Code Snippets",
		"Business Ideas",
		"Movie Watchlist",
		"Gift Ideas",
		"Home Improvement",
		"Financial Planning",
	}

	sampleNoteContents = []string{
		"This is a sample note content about various topics and ideas.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Important points to remember:\n1. First point\n2. Second point\n3. Third point",
		"Meeting agenda:\n- Review last week's progress\n- Discuss new features\n- Plan next sprint",
		"Ideas for the weekend:\n- Visit the museum\n- Try a new restaurant\n- Go for a hike",
		"Technical notes:\n- Use proper error handling\n- Implement logging\n- Add unit tests",
		"Personal goals:\n- Read more books\n- Exercise regularly\n- Learn a new skill",
		"Shopping items:\n- Groceries\n- Office supplies\n- Birthday gift for mom",
		"Travel checklist:\n- Book flights\n- Reserve hotel\n- Pack essentials\n- Check weather",
		"Daily thoughts and reflections on life, work, and personal growth.",
	}
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	config.ConnectDB()

	// Seed data
	if err := seedDatabase(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	fmt.Println("Database seeded successfully!")
}

func seedDatabase() error {
	db := config.GetDB()

	// Clear existing data (optional - uncomment if you want to reset)
	// if err := db.Exec("DELETE FROM notes").Error; err != nil {
	// 	return fmt.Errorf("failed to clear notes: %w", err)
	// }
	// if err := db.Exec("DELETE FROM users").Error; err != nil {
	// 	return fmt.Errorf("failed to clear users: %w", err)
	// }

	// Seed users
	fmt.Println("Seeding users...")
	var users []models.User
	for _, user := range sampleUsers {
		// Check if user already exists
		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			fmt.Printf("User %s already exists, skipping...\n", user.Email)
			users = append(users, existingUser)
			continue
		}

		// Create new user
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
		users = append(users, user)
		fmt.Printf("Created user: %s (%s)\n", user.Name, user.Email)
	}

	// Seed notes
	fmt.Println("\nSeeding notes...")
	rand.Seed(time.Now().UnixNano())

	for _, user := range users {
		// Create 3-7 notes per user
		numNotes := rand.Intn(5) + 3
		for i := 0; i < numNotes; i++ {
			note := models.Note{
				Title:   sampleNoteTitles[rand.Intn(len(sampleNoteTitles))],
				Content: sampleNoteContents[rand.Intn(len(sampleNoteContents))],
				UserID:  user.ID,
			}

			// Add some variation to titles to avoid duplicates
			if rand.Float32() < 0.3 {
				note.Title = fmt.Sprintf("%s %d", note.Title, rand.Intn(100))
			}

			if err := db.Create(&note).Error; err != nil {
				return fmt.Errorf("failed to create note for user %s: %w", user.Email, err)
			}
		}
		fmt.Printf("Created %d notes for user: %s\n", numNotes, user.Name)
	}

	// Print summary
	var userCount, noteCount int64
	db.Model(&models.User{}).Count(&userCount)
	db.Model(&models.Note{}).Count(&noteCount)

	fmt.Printf("\nSeeding completed!\n")
	fmt.Printf("Total users: %d\n", userCount)
	fmt.Printf("Total notes: %d\n", noteCount)

	// Print login credentials
	fmt.Printf("\nSample login credentials:\n")
	for _, user := range sampleUsers {
		fmt.Printf("Email: %s, Password: password123\n", user.Email)
	}

	return nil
}