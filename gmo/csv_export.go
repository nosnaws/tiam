package gmo

import (
	"fmt"
	"log"
	"os"
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
		"%s,%d,%f,%f,%f,%f,%f,%f\n",
		item.name,
		score,
		item.geno[0],
		item.geno[1],
		item.geno[2],
		item.geno[3],
		item.geno[4],
		item.geno[5],
	)
}

func createHeaderRow() string {
	return "NAME,FITNESS,EXPLORATION,ALPHA,VORONOI,FOODA,FOODB,BIG_SNAKE\n"
}
