package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Check if it is a valid PBM magic number
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("Invalid PBM magic number")
	}

	// Read width and height
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, errors.New("Invalid dimensions")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, errors.New("Invalid width")
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, errors.New("Invalid height")
	}

	// Read image data
	data := make([][]bool, height)
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		if magicNumber == "P1" {
			data[i] = parseP1Line(line, width)
		} else {
			data[i] = parseP4Line(line, width)
		}
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Helper function to parse P1 (ASCII) line
func parseP1Line(line string, width int) []bool {
	data := make([]bool, width)
	for i, char := range line {
		data[i] = char == '1'
	}
	return data
}

// Helper function to parse P4 (binary) line
func parseP4Line(line string, width int) []bool {
	data := make([]bool, width)

	// Ensure that the line has enough bytes to cover the width
	if len(line) < (width+7)/8 {
		return nil
	}

	for i := 0; i < width; i++ {
		// Calculate the byte index and bit position within the byte
		byteIndex := i / 8
		bitPos := uint(7 - (i % 8))

		// Extract the bit from the byte
		bit := (line[byteIndex] >> bitPos) & 1
		data[i] = bit == 1
	}

	return data
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
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

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for y := 0; y < pbm.height/2; y++ {
		pbm.data[y], pbm.data[pbm.height-y-1] = pbm.data[pbm.height-y-1], pbm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	filename := "C:/Users/JENGO/Netbpm/sample_640426.pbm"
	pbm, err := ReadPBM(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("Width: %d\n", pbm.width)
	fmt.Printf("Height: %d\n", pbm.height)

	width, height := pbm.Size()
	fmt.Printf("Image Size: %d x %d\n", width, height)

	x, y := 2, 3
	fmt.Printf("Pixel at (%d, %d): %v\n", x, y, pbm.At(x, y))

	newValue := true
	pbm.Set(x, y, newValue)
	fmt.Printf("New pixel value at (%d, %d): %v\n", x, y, pbm.At(x, y))

	outputFilename := "output.pbm"
	err = pbm.Save(outputFilename)
	if err != nil {
		fmt.Println("Error saving the PBM image:", err)
		return
	}

	fmt.Println("PBM image saved successfully to", outputFilename)
}