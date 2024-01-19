package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM struct to represent a PGM image
type PGM struct {
	data        [][]uint8
	width       int
	height      int
	magicNumber string
	max         int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("Unsupported PGM format")
	}

	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	scanner.Scan()
	maxVal, _ := strconv.Atoi(scanner.Text())

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
		for j := range data[i] {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(val)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxVal,
	}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintln(writer, pgm.data[i][j])
		}
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90Â° clockwise.
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
		for j := range newData[i] {
			newData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number, width, and height
	_, err = fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)
	if err != nil {
		return err
	}

	// Write image data
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pbm.magicNumber == "P1" {
				if pixel {
					_, err = file.WriteString("1 ")
				} else {
					_, err = file.WriteString("0 ")
				}
			} else {
				// For P4 format, write binary data
				if pixel {
					_, err = file.Write([]byte{0xFF})
				} else {
					_, err = file.Write([]byte{0x00})
				}
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbmData := make([][]bool, pgm.height)
	for i := 0; i < pgm.height; i++ {
		pbmData[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			pbmData[i][j] = pgm.data[i][j] > uint8(pgm.max)/2
		}
	}

	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1", // Assuming "P1" is the magic number for PBM
	}
}

func main() {
	// Read a PGM image from a file
	pgmFilename := "C:/Users/JENGO/Netbpm/paiement-passporr-timbre.pgm" // Change this to the actual path of your PGM file
	pgm, err := ReadPGM(pgmFilename)
	if err != nil {
		fmt.Println("Error reading PGM:", err)
		return
	}

	// Display PGM information
	fmt.Printf("PGM Magic Number: %s\n", pgm.magicNumber)
	fmt.Printf("PGM Width: %d\n", pgm.width)
	fmt.Printf("PGM Height: %d\n", pgm.height)
	fmt.Printf("PGM Max Value: %d\n", pgm.max)

	// Invert the colors of the PGM image
	pgm.Invert()

	// Flip the PGM image horizontally
	pgm.Flip()

	// Save the modified PGM image
	modifiedPGMFilename := "C:/Users/JENGO/Netbpm/pgmfile.pgm" // Change this to the desired output path
	err = pgm.Save(modifiedPGMFilename)
	if err != nil {
		fmt.Println("Error saving modified PGM:", err)
		return
	}
	fmt.Println("Modified PGM image saved successfully to", modifiedPGMFilename)

	// Convert the modified PGM image to PBM
	pbm := pgm.ToPBM()

	// Display PBM information
	fmt.Printf("\nPBM Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("PBM Width: %d\n", pbm.width)
	fmt.Printf("PBM Height: %d\n", pbm.height)

	// Save the PBM image
	pbmFilename := "C:/Users/JENGO/Netbpm/pgmfile.pgm" // Change this to the desired output path
	err = pbm.Save(pbmFilename)
	if err != nil {
		fmt.Println("Error saving PBM:", err)
		return
	}
	fmt.Println("PBM image saved successfully to", pbmFilename)
}
