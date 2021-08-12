package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var inputFile = flag.String("input", "", "input file in 1pif format")
var outputFile = flag.String("output", "", "output CSV file")

type onepifRow struct {
	Uuid           string `json:"uuid"`
	UpdatedAt      int    `json:"updatedAt"`
	LocationKey    string `json:"locationKey"`
	SecurityLevel  string `json:"securityLevel"`
	ContentsHash   string `json:"contentsHash"`
	Title          string `json:"title"`
	Location       string `json:"location"`
	SecureContents struct {
		Fields []struct {
			Value       string `json:"value"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			Id          string `json:"id,omitempty"`
			Designation string `json:"designation,omitempty"`
		} `json:"fields"`
		PasswordHistory []struct {
			Value string `json:"value"`
			Time  int    `json:"time"`
		} `json:"passwordHistory"`
		NotesPlain string `json:"notesPlain"`
		HtmlMethod string `json:"htmlMethod"`
		Sections   []struct {
			Fields []struct {
				K string `json:"k"`
				N string `json:"n"`
				V string `json:"v"`
				T string `json:"t"`
			} `json:"fields,omitempty"`
			Name  string `json:"name"`
			Title string `json:"title,omitempty"`
		} `json:"sections"`
		URLs []struct {
			Url string `json:"url"`
		} `json:"URLs"`
	} `json:"secureContents"`
	TxTimestamp int    `json:"txTimestamp"`
	CreatedAt   int    `json:"createdAt"`
	TypeName    string `json:"typeName"`
}

func (r *onepifRow) fieldWithDesignation(designation string) (string, error) {
	for _, field := range r.SecureContents.Fields {
		if field.Designation == designation {
			return field.Value, nil
		}
	}

	return "", fmt.Errorf("no %v found", designation)
}

func (r *onepifRow) username() (string, error) {
	return r.fieldWithDesignation("username")
}

func (r *onepifRow) password() (string, error) {
	return r.fieldWithDesignation("password")
}

func (r *onepifRow) otpAuth() (string, error) {
	for _, section := range r.SecureContents.Sections {
		for _, field := range section.Fields {
			if strings.HasPrefix(field.V, "otpauth://") {
				return field.V, nil
			}
		}
	}

	return "", errors.New("no otpauth found")
}

func main()  {
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	ifile, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer ifile.Close()

	ofile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer ofile.Close()

	csvWriter := csv.NewWriter(ofile)

	// write header
	csvWriter.Write([]string{"Title","Url","Username","Password","OTPAuth"})

	scanner := bufio.NewScanner(ifile)

	lineNr := 0
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()
		if strings.HasPrefix(line, "***") && strings.HasSuffix(line, "***") {
			continue
		}

		var row onepifRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			log.Printf("failed to parse line %d: %v\n", lineNr, err)
			continue
		}

		if len(row.SecureContents.URLs) == 0 {
			continue
		}

		username, _ := row.username()
		password, _ := row.password()
		otpauth, _ := row.otpAuth()

		for _, url := range row.SecureContents.URLs {
			record := []string{
				row.Title,
				url.Url,
				username,
				password,
				otpauth,
			}

			csvWriter.Write(record)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	csvWriter.Flush()
}
