package main

/*

Author Gaurav Sablok
Universitat Potsdam
Date 2024-9-11


doing the denovo alignment of the pacbiohifi reads against a protein database and holding those
alignment and the corresponding alignments in struct and extracting the matched regions. This is for
when you have a relatively less coverage and you want to search for the tags and the corresponding
protein alignments to generate annotation hints.  It has the following functions:
1. func readalignment : This aligns the reads to the protein datasets and prepares them for all the analysis.
2. func fastalignment : This prepares the fasta files for the alignment and prepares them for the analysis.
3. func tags: This prepares the coverage estimation for the alignment regions.
4. func sequence hints: This extracts the sequences upstream and the downstream

*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

var (
	pacbiofolder   string
	proteinfasta   string
	alignmentfile  string
	referencefasta string
	upstreamStart  int
	downstreamEnd  int
	pacbiofile     string
)

var rootCmd = &cobra.Command{
	Use:  "flag",
	Long: "This is golang application for generating the hints from the protein alignment to pacbiohifi reads",
}

var pacbioBamCmd = &cobra.Command{
	Use:  "pacbio",
	Long: "converts the pacbio reads to the fasta format for the alignment and annotations, provide the folder path",
	Run:  convertFunc,
}

var alignmentCmd = &cobra.Command{
	Use:  "proteinalignment",
	Long: "aligns the pacbio reads or the fasta reads to the proteins",
	Run:  alignCmd,
}

var hspCmd = &cobra.Command{
	Use:  "hspalignment",
	Long: "Analyzes the hsp from the diamond read to protein alignment",
	Run:  hspFunc,
}

var seqCmd = &cobra.Command{
	Use:  "seqHsp",
	Long: "Analyzes the hsp from the diamond read to protein alignment",
	Run:  getSeqFunc,
}

var upstreamCmd = &cobra.Command{
	Use:  "upStreamHSP",
	Long: "specific for the genome alignment regions upstream and the downstream of the alignments",
	Run:  upstreamFunc,
}

var analyzeCmd = &cobra.Command{
	Use:  "align and analyze",
	Long: "this option aligns and analyze all together and it requires only the proteins and the reference pacbio or the other fasta file",
	Run:  alignAnalyzeFunc,
}

func init() {
	pacbioBamCmd.Flags().
		StringVarP(&pacbiofolder, "pacbiofolder", "f", "folder containing the bam files and the pbi files", "pacbio conversion")
	alignmentCmd.Flags().
		StringVarP(&pacbiofile, "pacbio", "p", "reads for protein alignment", "read-protein alignment")
	alignmentCmd.Flags().
		StringVarP(&proteinfasta, "protein", "P", "protein datasets", "protein datasets for the alignment")
	hspCmd.Flags().
		StringVarP(&alignmentfile, "alignmentfile", "a", "alignment file to be analyzed", "alignment")
	hspCmd.Flags().
		StringVarP(&referencefasta, "referencefasta", "p", "pacbio reads file", "pacbio file")
	seqCmd.Flags().
		StringVarP(&alignmentfile, "alignmentfile", "a", "alignment file to be analyzed", "alignment")
	seqCmd.Flags().
		StringVarP(&referencefasta, "referencefasta", "p", "pacbio reads file", "pacbio file")
	upstreamCmd.Flags().
		StringVarP(&alignmentfile, "alignmentfile", "a", "alignment file to be analyzed", "alignment")
	upstreamCmd.Flags().
		StringVarP(&referencefasta, "referencefasta", "p", "pacbio reads file", "pacbio file")
	upstreamCmd.Flags().
		IntVarP(&upstreamStart, "upstream of the hsp tags", "u", 4, "upstream tags")
	upstreamCmd.Flags().
		IntVarP(&downstreamEnd, "downstream of the hsp tags", "d", 5, "downstream tags")
	analyzeCmd.Flags().
		StringVarP(&pacbiofolder, "pacbiofolder", "f", "folder containing the bam files and the pbi files", "pacbio conversion")
	analyzeCmd.Flags().
		StringVarP(&pacbiofile, "pacbio", "p", "reads for protein alignment", "read-protein alignment")
	analyzeCmd.Flags().
		StringVarP(&proteinfasta, "protein", "P", "protein datasets", "protein datasets for the alignment")

	rootCmd.AddCommand(seqCmd)
	rootCmd.AddCommand(upstreamCmd)
	rootCmd.AddCommand(pacbioBamCmd)
	rootCmd.AddCommand(alignmentCmd)
	rootCmd.AddCommand(hspCmd)
	rootCmd.AddCommand(analyzeCmd)
}

func convertFunc(cmd *cobra.Command, args []string) {
	convert, err := exec.Command("pbtk", "-o", "pacbio.fasta", "*.bam").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(convert))
}

func alignCmd(cmd *cobra.Command, args []string) {
	proteinreference := proteinfasta
	makedb, err := exec.Command("diamond", "makedb", "-in", proteinreference, "-d", "reference").
		Output()
	if err != nil {
		log.Fatal("command failed with %s\n", err)
	}
	fmt.Println(string(makedb))
	align, err := exec.Command("diamond", "blastx", "-d", "reference", "-q", "pacbio.fasta", "-o", "pacbioaligned.tsv").
		Output()
	if err != nil {
		log.Fatal("command failed with the error %s\n", err)
	}
	fmt.Println(string(align))
}

func sum(arr []float64) float64 {
	counter := float64(0)
	for i := range arr {
		counter += arr[i]
	}
	return counter
}

func pacbio() ([]string, []string, []float64) {
	readOpen, err := os.Open(referencefasta)
	if err != nil {
		log.Fatal(err)
	}

	readbuffer := bufio.NewScanner(readOpen)
	header := []string{}
	sequences := []string{}
	length := []float64{}

	for readbuffer.Scan() {
		line := readbuffer.Text()
		if string(line[0]) == "A" || string(line[0]) == "T" || string(line[0]) == "G" ||
			string(line[0]) == "C" {
			sequences = append(sequences, line)
		}
		if string(line[0]) == ">" {
			header = append(header, strings.ReplaceAll(string(line), ">", ""))
		}
	}
	for i := range sequences {
		length = append(length, float64(len(sequences[i])))
	}
	return header, sequences, length
}

func hspFunc(cmd *cobra.Command, args []string) {
	refID := []string{}
	alignID := []string{}
	refIdenStart := []float64{}
	refIdenEnd := []float64{}
	alignIdenStart := []float64{}
	alignIdenEnd := []float64{}
	fOpen, err := os.Open(alignmentfile)
	if err != nil {
		log.Fatal(err)
	}

	fRead := bufio.NewScanner(fOpen)

	for fRead.Scan() {
		line := fRead.Text()
		refID = append(refID, strings.Split(string(line), "\t")[0])
		alignID = append(alignID, strings.Split(string(line), "\t")[1])
		start1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[6], 32)
		end1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[7], 32)
		start2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[8], 32)
		end2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[9], 32)
		refIdenStart = append(refIdenStart, start1)
		refIdenEnd = append(refIdenEnd, end1)
		alignIdenStart = append(alignIdenStart, start2)
		alignIdenEnd = append(alignIdenEnd, end2)
	}
	id, _, length := pacbio()

	type cov struct {
		id  string
		cov float64
	}

	coverageSeq := []cov{}
	for i := range id {
		for j := range refID {
			if id[i] == refID[j] {
				coverageSeq = append(coverageSeq, cov{
					id:  refID[j],
					cov: (refIdenEnd[j] - refIdenStart[j]) / length[i] * 100,
				})
			}
		}
	}

	for i := range coverageSeq {
		fmt.Println(coverageSeq[i].id, coverageSeq[i].cov)
	}

	file, err := os.Create("coveragestimation.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := range coverageSeq {
		storeI := strconv.FormatFloat(coverageSeq[i].cov, 'f', -1, 64)
		_, err := file.WriteString(coverageSeq[i].id + "\t" + storeI + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getSeqFunc(cmd *cobra.Command, args []string) {
	refID := []string{}
	alignID := []string{}
	refIdenStart := []float64{}
	refIdenEnd := []float64{}
	alignIdenStart := []float64{}
	alignIdenEnd := []float64{}
	fOpen, err := os.Open(alignmentfile)
	if err != nil {
		log.Fatal(err)
	}

	fRead := bufio.NewScanner(fOpen)

	for fRead.Scan() {
		line := fRead.Text()
		refID = append(refID, strings.Split(string(line), "\t")[0])
		alignID = append(alignID, strings.Split(string(line), "\t")[1])
		start1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[6], 32)
		end1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[7], 32)
		start2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[8], 32)
		end2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[9], 32)
		refIdenStart = append(refIdenStart, start1)
		refIdenEnd = append(refIdenEnd, end1)
		alignIdenStart = append(alignIdenStart, start2)
		alignIdenEnd = append(alignIdenEnd, end2)
	}

	seqIDType, seqSeqType := readRef()

	type extractSeq struct {
		extractPartID  string
		extractPartSeq string
	}

	extractPartSeq := []extractSeq{}

	for i := range seqIDType {
		for j := range refID {
			if seqIDType[i] == refID[j] {
				extractPartSeq = append(extractPartSeq, extractSeq{
					extractPartID:  seqIDType[i],
					extractPartSeq: seqSeqType[i][int(refIdenStart[j]):int(refIdenEnd[j])],
				})
			}
		}
	}

	file, err := os.Create("sequences-annotation.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := range extractPartSeq {
		_, err := file.WriteString(
			">" + extractPartSeq[i].extractPartID + "\n" + extractPartSeq[i].extractPartSeq + "\n",
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func upstreamFunc(cmd *cobra.Command, args []string) {
	refID := []string{}
	alignID := []string{}
	refIdenStart := []float64{}
	refIdenEnd := []float64{}
	alignIdenStart := []float64{}
	alignIdenEnd := []float64{}
	fOpen, err := os.Open(alignmentfile)
	if err != nil {
		log.Fatal(err)
	}

	fRead := bufio.NewScanner(fOpen)

	for fRead.Scan() {
		line := fRead.Text()
		refID = append(refID, strings.Split(string(line), "\t")[0])
		alignID = append(alignID, strings.Split(string(line), "\t")[1])
		start1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[6], 32)
		end1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[7], 32)
		start2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[8], 32)
		end2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[9], 32)
		refIdenStart = append(refIdenStart, start1)
		refIdenEnd = append(refIdenEnd, end1)
		alignIdenStart = append(alignIdenStart, start2)
		alignIdenEnd = append(alignIdenEnd, end2)
	}

	seqIDType, seqSeqType := readRef()

	type extractStreamSeq struct {
		extractStreamPartID  string
		extractStreamPartSeq string
	}

	extractupstreamPartSeq := []extractStreamSeq{}

	for i := range seqIDType {
		for j := range refID {
			if seqIDType[i] == refID[j] {
				extractupstreamPartSeq = append(extractupstreamPartSeq, extractStreamSeq{
					extractStreamPartID:  seqIDType[i],
					extractStreamPartSeq: seqSeqType[i][int(int(refIdenStart[j])-upstreamStart):int(downstreamEnd+int(refIdenEnd[j]))],
				})
			}
		}

		file, err := os.Create("sequences-annotation-upstream-downstream.txt")
		if err != nil {
			log.Fatal(err)
		}
		for i := range extractupstreamPartSeq {
			_, err := file.WriteString(
				">" + extractupstreamPartSeq[i].extractStreamPartID + "\n" + extractupstreamPartSeq[i].extractStreamPartSeq + "\n",
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func alignAnalyzeFunc(cmd *cobra.Command, args []string) {
	readOpen, err := os.Open(referencefasta)
	if err != nil {
		log.Fatal(err)
	}

	readbuffer := bufio.NewScanner(readOpen)
	header := []string{}
	sequences := []string{}
	length := []float64{}

	for readbuffer.Scan() {
		line := readbuffer.Text()
		if string(line[0]) == "A" || string(line[0]) == "T" || string(line[0]) == "G" ||
			string(line[0]) == "C" {
			sequences = append(sequences, line)
		}
		if string(line[0]) == ">" {
			header = append(header, strings.ReplaceAll(string(line), ">", ""))
		}
	}
	for i := range sequences {
		length = append(length, float64(len(sequences[i])))
	}

	refID := []string{}
	alignID := []string{}
	refIdenStart := []float64{}
	refIdenEnd := []float64{}
	alignIdenStart := []float64{}
	alignIdenEnd := []float64{}
	fOpen, err := os.Open(alignmentfile)
	if err != nil {
		log.Fatal(err)
	}

	fRead := bufio.NewScanner(fOpen)

	for fRead.Scan() {
		line := fRead.Text()
		refID = append(refID, strings.Split(string(line), "\t")[0])
		alignID = append(alignID, strings.Split(string(line), "\t")[1])
		start1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[6], 32)
		end1, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[7], 32)
		start2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[8], 32)
		end2, _ := strconv.ParseFloat(strings.Split(string(line), "\t")[9], 32)
		refIdenStart = append(refIdenStart, start1)
		refIdenEnd = append(refIdenEnd, end1)
		alignIdenStart = append(alignIdenStart, start2)
		alignIdenEnd = append(alignIdenEnd, end2)
	}
	id, _, length := pacbio()

	type cov struct {
		id  string
		cov float64
	}

	coverageSeq := []cov{}
	for i := range id {
		for j := range refID {
			if id[i] == refID[j] {
				coverageSeq = append(coverageSeq, cov{
					id:  refID[j],
					cov: (refIdenEnd[j] - refIdenStart[j]) / length[i] * 100,
				})
			}
		}
	}

	for i := range coverageSeq {
		fmt.Println(coverageSeq[i].id, coverageSeq[i].cov)
	}

	file, err := os.Create("coveragestimation.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := range coverageSeq {
		storeI := strconv.FormatFloat(coverageSeq[i].cov, 'f', -1, 64)
		_, err := file.WriteString(coverageSeq[i].id + "\t" + storeI + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	seqIDType, seqSeqType := readRef()

	type extractSeq struct {
		extractPartID  string
		extractPartSeq string
	}

	extractPartSeq := []extractSeq{}

	for i := range seqIDType {
		for j := range refID {
			if seqIDType[i] == refID[j] {
				extractPartSeq = append(extractPartSeq, extractSeq{
					extractPartID:  seqIDType[i],
					extractPartSeq: seqSeqType[i][int(refIdenStart[j]):int(refIdenEnd[j])],
				})
			}
		}
	}

	file1, err := os.Create("sequences-annotation.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := range extractPartSeq {
		_, err := file1.WriteString(
			">" + extractPartSeq[i].extractPartID + "\n" + extractPartSeq[i].extractPartSeq + "\n",
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	type extractStreamSeq struct {
		extractStreamPartID  string
		extractStreamPartSeq string
	}

	extractupstreamPartSeq := []extractStreamSeq{}

	for i := range seqIDType {
		for j := range refID {
			if seqIDType[i] == refID[j] {
				extractupstreamPartSeq = append(extractupstreamPartSeq, extractStreamSeq{
					extractStreamPartID:  seqIDType[i],
					extractStreamPartSeq: seqSeqType[i][int(int(refIdenStart[j])-upstreamStart):int(downstreamEnd+int(refIdenEnd[j]))],
				})
			}
		}

		file, err := os.Create("sequences-annotation-upstream-downstream.txt")
		if err != nil {
			log.Fatal(err)
		}
		for i := range extractupstreamPartSeq {
			_, err := file.WriteString(
				">" + extractupstreamPartSeq[i].extractStreamPartID + "\n" + extractupstreamPartSeq[i].extractStreamPartSeq + "\n",
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func readRef() ([]string, []string) {
	extractID := []string{}
	seqID := []string{}

	readF, err := os.Open(referencefasta)
	if err != nil {
		log.Fatal(err)
	}

	openF := bufio.NewScanner(readF)

	for openF.Scan() {
		line := openF.Text()
		if string(line[0]) == "A" || string(line[0]) == "T" || string(line[0]) == "G" ||
			string(line[0]) == "C" {
			seqID = append(seqID, line)
		}
		if string(line[0]) == ">" {
			extractID = append(extractID, strings.ReplaceAll(string(line), ">", ""))
		}
	}
	return extractID, seqID
}
