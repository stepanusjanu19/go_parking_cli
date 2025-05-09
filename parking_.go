package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Car struct {
	RegistrationNumber string
	SlotNumber         int
	ArrivalTime        time.Time
}

type ParkingLot struct {
	Capacity       int
	ParkedCars     map[int]*Car
	AvailableSlots []int
}

var parkingLot *ParkingLot = &ParkingLot{}

func (p *ParkingLot) Create(capacity int) {
	p.Capacity = capacity
	p.ParkedCars = make(map[int]*Car)
	p.AvailableSlots = make([]int, capacity)
	for i := 0; i < capacity; i++ {
		p.AvailableSlots[i] = i + 1
	}
	fmt.Printf("Created a parking lot with %d slots\n", capacity)
}

func (p *ParkingLot) Park(regNum string) {
	if len(p.AvailableSlots) == 0 {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	sort.Ints(p.AvailableSlots)
	slot := p.AvailableSlots[0]
	p.AvailableSlots = p.AvailableSlots[1:]
	car := &Car{
		RegistrationNumber: regNum,
		SlotNumber:         slot,
		ArrivalTime:        time.Now(),
	}
	p.ParkedCars[slot] = car
	fmt.Printf("Allocated slot number: %d\n", slot)
}

func (p *ParkingLot) Leave(regNum string, hours int) {

	for slot, car := range p.ParkedCars {
		if car.RegistrationNumber == regNum {
			charge := 10
			if hours > 2 {
				charge += (hours - 2) * 10
			}
			fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n",
				regNum, slot, charge)
			p.AvailableSlots = append(p.AvailableSlots, slot)
			delete(p.ParkedCars, slot)
			return
		}
	}
	fmt.Printf("Registration number %s not found\n", regNum)
}

func (p *ParkingLot) Status() {
	if len(p.ParkedCars) == 0 {
		fmt.Println("Parking lot is empty")
		return
	}
	fmt.Println("Slot No. Registration No.")
	slots := []int{}
	for slot := range p.ParkedCars {
		slots = append(slots, slot)
	}
	sort.Ints(slots)
	for _, slot := range slots {
		car := p.ParkedCars[slot]
		fmt.Printf("%d %s\n", slot, car.RegistrationNumber)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	logFile, err := os.OpenFile("input.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()
	logger := bufio.NewWriter(logFile)

	log := func(format string, a ...interface{}) {
		output := fmt.Sprintf(format, a...)
		logger.WriteString(output)
		logger.Flush()
	}

	fmt.Println("===================Automate Parking Lot===================")
	fmt.Println("Choose a command:")
	fmt.Println("1. Create parking lot of size n : create_parking_lot {capacity}")
	fmt.Println("2. Park a car : park {car_number}")
	fmt.Println("3. Remove(Unpark) car from : leave {car_number} {hours}")
	fmt.Println("4. Print status of parking slot : status")
	fmt.Println("5. Exit : exit")
	fmt.Println("==========================================================")
	fmt.Print("> ")

	for scanner.Scan() {

		line := scanner.Text()
		log(line + "\n")
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "create_parking_lot":
			if len(args) != 1 {
				fmt.Println("Usage: create_parking_lot <capacity>")
				continue
			}
			cap, _ := strconv.Atoi(args[0])
			parkingLot.Create(cap)

		case "park":
			if len(args) != 1 {
				fmt.Println("Usage: park <registration_number>")
				continue
			}
			parkingLot.Park(args[0])

		case "leave":
			if len(args) != 2 {
				fmt.Println("Usage: leave <registration_number> <hours>")
				continue
			}
			hours, _ := strconv.Atoi(args[1])
			parkingLot.Leave(args[0], hours)

		case "status":
			parkingLot.Status()

		case "exit":
			fmt.Println("=============================================")
			fmt.Println("============Thank You Use This Apps==========")
			fmt.Println("=============================================")
			return

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
