package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"math/rand"
)

type Position int

const (
	Goalkeeper Position = iota
	Defender
	Midfielder
	Striker
)

type Player struct {
	Position  Position
	Shooting  int
	Defending int
	Tackling  int
	Passing   int
	Physical  int
	Speed     int
	Attacking int
}

type Team struct {
	Name        string
	Goalkeeper  Player
	Defenders   []Player
	Midfielders []Player
	Strikers    []Player
}

type Match struct {
	TeamA           Team
	TeamB           Team
	BallPosition    string // "defensive", "midfield", "final"
	Possession      *Team
	TeamAPossession int
	TeamAAttack     int
	TeamBPossession int
	TeamBAttack     int
	Iterations      int
	ScoreA          int
	ScoreB          int
}

func generateTeam(name string, formation string) Team {
	team := Team{Name: name}

	playersFormation := strings.Split(formation, "-") // devide formation to positions

	defs, _ := strconv.Atoi(playersFormation[0])     // get defenders count
	mids, _ := strconv.Atoi(playersFormation[1])     // get mids count
	strikers, _ := strconv.Atoi(playersFormation[2]) // get strikers count

	// generate players
	team.Goalkeeper = generatePlayer(Goalkeeper)
	for i := 0; i < defs; i++ {
		team.Defenders = append(team.Defenders, generatePlayer(Defender))
	}
	for i := 0; i < mids; i++ {
		team.Midfielders = append(team.Midfielders, generatePlayer(Midfielder))
	}
	for i := 0; i < strikers; i++ {
		team.Strikers = append(team.Strikers, generatePlayer(Striker))
	}

	return team
}

func generatePlayer(position Position) Player {
	p := Player{Position: position}

	switch position {
	case Goalkeeper:
		p.Shooting = rand.Intn(20) + 10
		p.Defending = rand.Intn(30) + 70
		p.Passing = rand.Intn(40) + 30
		p.Physical = rand.Intn(30) + 70
		p.Speed = rand.Intn(40) + 30
		p.Attacking = rand.Intn(20) + 10
	case Defender:
		p.Shooting = rand.Intn(30) + 20
		p.Defending = rand.Intn(30) + 60
		p.Passing = rand.Intn(40) + 40
		p.Physical = rand.Intn(30) + 60
		p.Speed = rand.Intn(40) + 40
		p.Attacking = rand.Intn(30) + 20
	case Midfielder:
		p.Shooting = rand.Intn(40) + 30
		p.Defending = rand.Intn(40) + 40
		p.Passing = rand.Intn(30) + 70
		p.Physical = rand.Intn(40) + 50
		p.Speed = rand.Intn(40) + 50
		p.Attacking = rand.Intn(40) + 30
	case Striker:
		p.Shooting = rand.Intn(30) + 70
		p.Defending = rand.Intn(20) + 10
		p.Passing = rand.Intn(40) + 40
		p.Physical = rand.Intn(40) + 40
		p.Speed = rand.Intn(30) + 60
		p.Attacking = rand.Intn(30) + 70
	}
	return p
}

func (m *Match) determineInitialPossession() {
	isAAttack := m.Possession != nil && m.Possession.Name == m.TeamA.Name
	isBAttack := m.Possession != nil && m.Possession.Name == m.TeamB.Name
	aStr := calculateMidfieldStrength(m.TeamA, isAAttack)
	bStr := calculateMidfieldStrength(m.TeamB, isBAttack)

	if aStr+rand.Intn(aStr/2) > bStr+rand.Intn(bStr/2) {
		fmt.Println("Team A won the possession")
		m.Possession = &m.TeamA
	} else {
		fmt.Println("Team B won the possession")
		m.Possession = &m.TeamB
	}
}

func calculateMidfieldStrength(team Team, isAttacking bool) int {
	strength := 0
	for _, m := range team.Midfielders {
		strength += m.Passing
		if isAttacking {
			strength += m.Attacking
		} else {
			strength += m.Defending
		}
		strength += m.Physical
	}
	return strength
}

func calculateAttackStrength(team Team) int {
	strength := 0
	for _, m := range team.Strikers {
		strength += m.Shooting
		strength += m.Attacking
		strength += m.Physical
		strength += m.Speed
	}
	return strength
}

func calculateDefendStrength(team Team) int {
	strength := 0
	for _, m := range team.Defenders {
		strength += m.Tackling
		strength += m.Attacking
		strength += m.Physical
		strength += m.Speed
	}
	return strength
}

func (m *Match) simulate() {
	ticker := time.NewTicker(3 * time.Second) // Tick every second
	defer ticker.Stop()

	for i := 1; i <= m.Iterations; i++ {
		fmt.Printf("== minute %d'==\n", i)
		if m.Possession != nil {
			if m.Possession.Name == m.TeamA.Name {
				m.TeamAPossession += 1
			} else {
				m.TeamBPossession += 1
			}
		}

		switch m.BallPosition {
		case "mid":
			m.midThirdAction()
		case "final":
			m.finalThirdAction()
		}

		fmt.Println()
		// Wait for the next tick
		<-ticker.C
	}
}

func (m *Match) midThirdAction() {
	isAAttack := m.Possession.Name == m.TeamA.Name
	isBAttack := m.Possession.Name == m.TeamB.Name
	aStr := calculateMidfieldStrength(m.TeamA, isAAttack)
	bStr := calculateMidfieldStrength(m.TeamB, isBAttack)
	fmt.Printf("midfield battle, %s trying to attack\n", m.Possession.Name)

	// Add random factor (up to 50% of total strength)
	aTotal := aStr + rand.Intn(aStr/2)
	bTotal := bStr + rand.Intn(bStr/2)

	// Use a new random source for better randomness
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// re-iterrate from midfield
	if rng.Float64() < 0.1 {
		fmt.Println("midfield battle is a draw, init possession again")
		m.Possession = nil
		m.determineInitialPossession()
		return
	}

	if m.Possession.Name == m.TeamA.Name { // current possession is team A
		if aTotal > bTotal {
			fmt.Println("Team A won battle, proceed to final third")
			m.BallPosition = "final"
		} else {
			fmt.Println("Team B steal the ball, back to midfield battle")
			m.Possession = &m.TeamB // change possession, re-iterrate from midfield
		}
	} else { // current possession is team B
		if bTotal > aTotal {
			fmt.Println("team B won battle, proceed to final third")
			m.BallPosition = "final"
		} else {
			fmt.Println("Team A steal the ball, back to midfield battle")
			m.Possession = &m.TeamA // change possession, re-iterrate from midfield
		}
	}
	fmt.Printf("*battle score [%d vs %d]\n", aTotal, bTotal)
}

func (m *Match) finalThirdAction() {
	var attackingTeam, defendingTeam Team
	if m.Possession.Name == m.TeamA.Name {
		attackingTeam = m.TeamA
		defendingTeam = m.TeamB
		m.TeamAAttack += 1
	} else {
		attackingTeam = m.TeamB
		defendingTeam = m.TeamA
		m.TeamBAttack += 1
	}

	// Calculate attack strength
	attackStr := calculateAttackStrength(attackingTeam)

	// Calculate defense strength
	defenseStr := calculateDefendStrength(defendingTeam)
	defenseStr += defendingTeam.Goalkeeper.Defending

	// Add randomness
	attackTotal := attackStr + rand.Intn(attackStr)
	defenseTotal := defenseStr + rand.Intn(defenseStr)

	if attackTotal > defenseTotal {
		// Goal scored
		if attackingTeam.Name == m.TeamA.Name {
			m.ScoreA++
		} else {
			m.ScoreB++
		}
		fmt.Printf("GOAL! %s scores! (%d-%d)\n", attackingTeam.Name, m.ScoreA, m.ScoreB)
	} else {
		fmt.Printf("%s success defending, back to mid\n", defendingTeam.Name)
		fmt.Printf("attacking score: %d vs defending score: %d\n", attackTotal, defenseTotal)
	}

	// Reset to midfield
	m.BallPosition = "mid"
	m.Possession = &defendingTeam
	// m.determineInitialPossession()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate teams
	teamA := generateTeam("Team A", "4-5-1")
	teamB := generateTeam("Team B", "4-4-2")

	fmt.Printf("Team A strength: ATT: %d, MID: %d, DEF: %d\n", calculateAttackStrength(teamA), calculateMidfieldStrength(teamA, false), calculateDefendStrength(teamA)+teamA.Goalkeeper.Defending)
	fmt.Printf("Team B strength: ATT: %d, MID: %d, DEF: %d\n", calculateAttackStrength(teamB), calculateMidfieldStrength(teamB, false), calculateDefendStrength(teamB)+teamA.Goalkeeper.Attacking)

	// Create match
	match := Match{
		TeamA:        teamA,
		TeamB:        teamB,
		BallPosition: "mid",
		Iterations:   45,
	}

	// Determine initial possession
	match.determineInitialPossession()

	// Run match simulation
	match.simulate()

	// Show results
	fmt.Printf("\nFinal Score: %s %d - %d %s\n",
		match.TeamA.Name, match.ScoreA,
		match.ScoreB, match.TeamB.Name)
	fmt.Println("Team Stats:")
	fmt.Printf("Ball Possession: %d%% - %d%% (%d - %d)\n", (match.TeamAPossession*100)/match.Iterations, (match.TeamBPossession*100)/match.Iterations, match.TeamAPossession, match.TeamBPossession)
	fmt.Printf("Attack: %d - %d", match.TeamAAttack, match.TeamBAttack)

}
