package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	subtle    = lipgloss.Color("#181818")
	highlight = lipgloss.Color("#7D56F4")
	special   = lipgloss.Color("#181818")

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(subtle).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(special).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(special).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)
)

type Character struct {
	// Basic Info
	Name       string `json:"name"`
	Race       string `json:"race"`
	Class      string `json:"class"`
	Level      int    `json:"level"`
	Background string `json:"background"`

	// Ability Scores
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`

	// Additional Info
	Alignment     string   `json:"alignment"`
	Experience    int      `json:"experience"`
	HitPoints     int      `json:"hitPoints"`
	ArmorClass    int      `json:"armorClass"`
	Initiative    int      `json:"initiative"`
	Speed         int      `json:"speed"`
	Proficiencies []string `json:"proficiencies"`
	Languages     []string `json:"languages"`
	Equipment     []string `json:"equipment"`
}

func main() {
	characters := loadCharacters()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived Ctrl+C. Saving characters and exiting...")
		saveCharacters(characters)
		os.Exit(0)
	}()

	for {
		var action string
		huh.NewSelect[string]().
			Title("What would you like to do?").
			Options(
				huh.NewOption("Create a new character", "create"),
				huh.NewOption("View all characters", "view"),
				huh.NewOption("Delete a character", "delete"),
				huh.NewOption("Exit", "exit"),
			).
			Value(&action).
			Run()

		switch action {
		case "create":
			newCharacter, created := createCharacter()
			if created {
				characters = append(characters, newCharacter)
				saveCharacters(characters)
				fmt.Println("Character added successfully!")
			} else {
				fmt.Println("Character creation cancelled.")
			}
		case "view":
			viewCharacters(characters)
		case "delete":
			characters = deleteCharacter(characters)
			saveCharacters(characters)
		case "exit":
			fmt.Println("Saving characters and exiting. Goodbye!")
			saveCharacters(characters)
			return
		}
	}
}

func createCharacter() (Character, bool) {
	var character Character
	var levelStr, expStr, hpStr, acStr, initStr, speedStr string
	var strStr, dexStr, conStr, intStr, wisStr, chaStr string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&character.Name),
			huh.NewSelect[string]().
				Title("Race").
				Options(
					huh.NewOption("Human", "Human"),
					huh.NewOption("Elf", "Elf"),
					huh.NewOption("Dwarf", "Dwarf"),
					huh.NewOption("Halfling", "Halfling"),
					huh.NewOption("Orc", "Orc"),
				).
				Value(&character.Race),
			huh.NewSelect[string]().
				Title("Class").
				Options(
					huh.NewOption("Fighter", "Fighter"),
					huh.NewOption("Wizard", "Wizard"),
					huh.NewOption("Rogue", "Rogue"),
					huh.NewOption("Cleric", "Cleric"),
					huh.NewOption("Monk", "Monk"),
				).
				Value(&character.Class),
			huh.NewInput().Title("Level").Value(&levelStr),
			huh.NewInput().Title("Background").Value(&character.Background),
			huh.NewInput().Title("Alignment").Value(&character.Alignment),
			huh.NewInput().Title("Experience Points").Value(&expStr),
			huh.NewInput().Title("Hit Points").Value(&hpStr),
			huh.NewInput().Title("Armor Class").Value(&acStr),
			huh.NewInput().Title("Initiative").Value(&initStr),
			huh.NewInput().Title("Speed").Value(&speedStr),
			huh.NewInput().Title("Strength").Value(&strStr),
			huh.NewInput().Title("Dexterity").Value(&dexStr),
			huh.NewInput().Title("Constitution").Value(&conStr),
			huh.NewInput().Title("Intelligence").Value(&intStr),
			huh.NewInput().Title("Wisdom").Value(&wisStr),
			huh.NewInput().Title("Charisma").Value(&chaStr),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Error during character creation:", err)
		return Character{}, false
	}

	// Convert string inputs to integers with error handling
	convertToInt := func(s string) int {
		value, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("Invalid input for %s: %s. Using default value 0.\n", s, err)
			return 0
		}
		return value
	}

	character.Level = convertToInt(levelStr)
	character.Experience = convertToInt(expStr)
	character.HitPoints = convertToInt(hpStr)
	character.ArmorClass = convertToInt(acStr)
	character.Initiative = convertToInt(initStr)
	character.Speed = convertToInt(speedStr)
	character.Strength = convertToInt(strStr)
	character.Dexterity = convertToInt(dexStr)
	character.Constitution = convertToInt(conStr)
	character.Intelligence = convertToInt(intStr)
	character.Wisdom = convertToInt(wisStr)
	character.Charisma = convertToInt(chaStr)

	// Function to handle repeating inputs with feedback
	addRepeatingInput := func(prompt string, slice *[]string) {
		for {
			var input string
			err := huh.NewInput().
				Title(prompt).
				Value(&input).
				Run()

			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}

			input = strings.TrimSpace(input)
			if input == "" {
				break
			}

			*slice = append(*slice, input)
			fmt.Printf("Added: %s\n", input)
		}
	}

	fmt.Println("\nAdding Proficiencies:")
	addRepeatingInput("Enter a proficiency (or leave blank to finish)", &character.Proficiencies)

	fmt.Println("\nAdding Languages:")
	addRepeatingInput("Enter a language (or leave blank to finish)", &character.Languages)

	fmt.Println("\nAdding Equipment:")
	addRepeatingInput("Enter an equipment item (or leave blank to finish)", &character.Equipment)

	// Confirm character creation
	var confirm bool
	huh.NewConfirm().
		Title("Do you want to add this character?").
		Value(&confirm).
		Run()

	if confirm {
		return character, true
	} else {
		return Character{}, false
	}
}

func viewCharacters(characters []Character) {
	if len(characters) == 0 {
		fmt.Println("No characters found.")
		return
	}

	for i, char := range characters {
		// Title
		fmt.Printf("%s\n", titleStyle.Render(fmt.Sprintf("Character %d: %s", i+1, char.Name)))

		// Basic Info
		basicInfo := lipgloss.JoinHorizontal(lipgloss.Top,
			subtitleStyle.Render("Race"),
			normalStyle.Render(char.Race),
			subtitleStyle.Render("Class"),
			normalStyle.Render(char.Class),
			subtitleStyle.Render("Level"),
			normalStyle.Render(fmt.Sprintf("%d", char.Level)),
		)
		fmt.Println(basicInfo)

		// Background and Alignment
		bgAlign := lipgloss.JoinHorizontal(lipgloss.Top,
			subtitleStyle.Render("Background"),
			normalStyle.Render(char.Background),
			subtitleStyle.Render("Alignment"),
			normalStyle.Render(char.Alignment),
		)
		fmt.Println(bgAlign)

		// Abilities
		abilities := lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinHorizontal(lipgloss.Top,
				infoStyle.Render("STR"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Strength)),
				infoStyle.Render("DEX"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Dexterity)),
				infoStyle.Render("CON"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Constitution)),
			),
			lipgloss.JoinHorizontal(lipgloss.Bottom,
				infoStyle.Render("INT"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Intelligence)),
				infoStyle.Render("WIS"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Wisdom)),
				infoStyle.Render("CHA"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Charisma)),
			),
		)
		fmt.Println(abilities)

		// Combat Stats
		combatStats := lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.JoinHorizontal(lipgloss.Top,
				infoStyle.Render("HP"),
				normalStyle.Render(fmt.Sprintf("%3d", char.HitPoints)),
				infoStyle.Render("Armor Class"),
				normalStyle.Render(fmt.Sprintf("%2d", char.ArmorClass)),
				infoStyle.Render("Initiative"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Initiative)),
				infoStyle.Render("Speed"),
				normalStyle.Render(fmt.Sprintf("%2d", char.Speed)),
			),
		)
		fmt.Println(combatStats)

		// Proficiencies, Languages, and Equipment
		fmt.Printf("%s %s\n", subtitleStyle.Render("Proficiencies:"), normalStyle.Render(strings.Join(char.Proficiencies, ", ")))
		fmt.Printf("%s %s\n", subtitleStyle.Render("Languages:"), normalStyle.Render(strings.Join(char.Languages, ", ")))
		fmt.Printf("%s %s\n", subtitleStyle.Render("Equipment:"), normalStyle.Render(strings.Join(char.Equipment, ", ")))

		fmt.Println() // Add a blank line between characters
	}
}

func deleteCharacter(characters []Character) []Character {
	if len(characters) == 0 {
		fmt.Println("No characters to delete.")
		return characters
	}

	viewCharacters(characters)

	var indexStr string
	huh.NewInput().
		Title("Enter the number of the character to delete (or 0 to cancel)").
		Value(&indexStr).
		Run()

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index > len(characters) {
		fmt.Println("Invalid input. No character deleted.")
		return characters
	}

	if index == 0 {
		fmt.Println("Deletion cancelled.")
		return characters
	}

	deletedCharacter := characters[index-1]
	characters = append(characters[:index-1], characters[index:]...)
	fmt.Printf("Character '%s' has been deleted.\n", deletedCharacter.Name)

	return characters
}

func loadCharacters() []Character {
	data, err := os.ReadFile("characters.json")
	if err != nil {
		return []Character{}
	}

	var characters []Character
	err = json.Unmarshal(data, &characters)
	if err != nil {
		fmt.Println("Error loading characters:", err)
		return []Character{}
	}

	return characters
}

func saveCharacters(characters []Character) {
	data, err := json.MarshalIndent(characters, "", "  ")
	if err != nil {
		fmt.Println("Error saving characters:", err)
		return
	}

	err = os.WriteFile("characters.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
