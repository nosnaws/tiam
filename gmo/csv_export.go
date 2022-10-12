package gmo

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func writeCSV(gen int, candidates []canditate, score map[string]int) {
	outDir := os.Getenv("CSV_OUTPUT_DIR")

	if outDir != "" {
		fmt.Println("Writing generation to file")
		file, err := os.Create(fmt.Sprintf("%s/generation_%d.csv", outDir, gen))
		if err != nil {
			log.Println("ERROR: Unable to create file", err)
			return
		}
		defer file.Close()

		output := createHeaderRow()
		for _, i := range candidates {
			output = output + createRow(i, score[i.name])
		}

		_, err = file.WriteString(output)
		if err != nil {
			log.Println("ERROR: Unable to save generation data set to disc", err)
		}
	}
}

func createRow(item canditate, score int) string {
	return fmt.Sprintf(
		"%s,%d,%f,%f,%f,%f,%f\n",
		item.name,
		score,
		item.geno[0],
		item.geno[1],
		item.geno[2],
		item.geno[3],
		item.geno[4],
	)
}

func createHeaderRow() string {
	return "NAME,FITNESS,FOOD_A,FOOD_B,VORONOI_A,VORONOI_B,LENGTH_WEIGHT\n"
}

func CombineCSVs() {
	genRegex := regexp.MustCompile(`generation_(.+).csv`)
	baseDir := "/Users/aswanson/dev/innovation/tiam/evo_out"
	dir, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatalln(err)
	}

	allRows := []string{}
	allRows = append(allRows, "NAME,FITNESS,FOOD_A,FOOD_B,VORONOI_A,VORONOI_B,LENGTH_WEIGHT,GENERATION\n")
	for _, file := range dir {
		genNum := genRegex.FindStringSubmatch(file.Name())
		fileDir := baseDir + "/" + file.Name()

		raw, err := os.ReadFile(fileDir)
		if err != nil {
			log.Fatalln(err)
		}

		contents := string(raw)

		rows := strings.Split(contents, "\n")
		rows = rows[1:]           //remove header row
		rows = rows[:len(rows)-1] // remove empty last line

		for _, row := range rows {
			newRow := row + fmt.Sprintf(",%s\n", genNum[1])
			allRows = append(allRows, newRow)
		}
	}

	newFile, err := os.Create(baseDir + "/combined.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer newFile.Close()

	output := ""
	for _, row := range allRows {
		output = output + row
	}

	newFile.WriteString(output)
}
