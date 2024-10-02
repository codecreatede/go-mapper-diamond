package main

/*

Author Gaurav Sablok
Universitat Potsdam
Date 2024-9-11


doing the denovo alignment of the pacbiohifi reads against a protein database and holding those
alignment and the corresponding alignments in struct and extracting the matched regions. This is for
when you have a relatively less coverage and you want to search for the tags and the corresponding
protein alignments to generate annotation hints.

*/


import (
	"fmt"
	"flag"
	"os"
	"os/exec"
	"log"
	"strings"
	"github.com/spf13/cobra"
)

func main () {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

var pacbio string
var protein string
var alignment string

var rootCmd = &cobra.Command {
	Use: "flag",
	Long: "This is golang application for generating the hints from the protein alignment to pacbiohifi reads",
}

var alignCmd = &cobra.Command {
	Use: "align",
	Long: "aligns the reads to the protein database and then generates the hints",
	Run: alignFunc,

}

var analyzeCmd = &cobra.Command {
	Use: "analyze",
	Long: "analyze the already aligned reads to protein alignment",
	Run: analyzeFunc,
}


func init() {

	alignCmd.Flags().StringVarP(&pacbio, "pacbio", "p", "reads for protein alignment", "read-protein alignment")
	alignCmd.Flags().StringVarP(&protein, "protein", "P", "protein datasets", "protein datasets for the alignment")
	analyzeCmd.Flags().StringVarP(&alignment, "alignment", "a", "analyze the already given alignment", "post alignment analyzer")

  rootCmd.AddCommand(alignCmd)
  rootCmd.AddCommand(analyzeCmd)

}


func alignCmd (cmd *cobra.Command, args [] string){

	type mapperID struct {
		id string
	}
	type mapperSeq struct {
		seq string
	}

	pacbioOpen, err := os.Open(pacbio)
	if err != nil {
		log.Fatal (err)
	}

	pacbioRead := bufio.NewScanner(pacbioOpen)

  seqID = []mapperID{}
  seqSeq = []mapperSeq{}

  for pacbioRead.Scan() {
		line := pacBioRead.Text()
		if strings.HasPrefix(string(line), "@") {
			seqID = append(seqID, mapperID{
        id : string.Split(string(line), "\t")[0],
			})
		if strings.HasPrefix(string(line), "A") || strings.HasPrefix(string(line), "T") || strings.HasPrefix(string(line), "C") || strings.HasPrefix(string(line), "C")
			seqSeq = append(seqSeq, mapperSeq{
				seq : string(line)
			})
			combinedSeqID = []{seqID, seqSeq}
		}
			writeFasta, err := os.Create("pacbio.fasta", os.O_CREATE|os.O_WRONLY, 0644)
		  if err ! = nil {
			log.Fatal(err)
		 }
		 defer writeFasta.Close()
     for i := range seqID {
			for j := range seqSeq {
				write, err := writeFasta.write(>seqID[i]\nseqSeq[j])
				if err ! = nil {
					log.Fatal(err)
				}
			}

	  out, err := exec.Command("diamond", "makedb", "-in", "pacbio.fasta", "-d", "reference").Output()
		if err != nil {
					log.Fatal("command failed with %s\n", err)
				}
		fmt.Println(out)

		out, err := exec.Command("diamond", "blastx", "-d", "reference", "-q", "*readsFasta", "-o", "aligned.tsv").Output()
	  if err != nil {
					log.Fatal("command failed with the error %s\n", err)
				}
		fmt.Println(out)
		 }

	type hspStruct struct {
       refseq string
			 aligned string
			 alignmentMatch string
			 start string
			 end string
			 evalue string
			}
	alignmentOpen, err := os.Open("aligned.tsv")
			if err != nil {
				log.Fatal(err)
			}
	alignmentRead, err := bufio.NewScanner(alignmentOpen)
	holdAlignment := []hspStruct{}
	for alignmentRead.Scan() {
				line := alignmentRead.Text()
				holdAlignment := append(holadalignment, hspStruct{
					refseq : strings.Split(line, "\t")[0]
					aligned : strings.Split(line, "\t")[1]
					alignedMatch : strings.Split(line, "\t")[2]
					start : strings.Split(line, "\t")[6]
					end : strings.Split(line, "\t")[7]
					evalue: strings.Split(line, "\t")[11]
				})
			}

  type nucleotideIDStruct struct {
				id string
			}

	type nucleotideSeqstruct struct {
				seq string
			}

	type alignedFetch struct {
				seqID string
				seqSeq string
			}

  hspSeq := []alignedFetch{}
	for i := range holdAlignment {
		for j := range seqIDWrite {
			if seqIDwrite[i] == holdAlignment[i].refseq {
					hspSeq = append(hspSeq, alignedFetch{
						seqID : seqIDwrite[i]
						seqSeq : SeqSeqWrite[i][int(holdAlignment[i].start):int(holdAlignment[i].end)]
					})
				}
			}
				}

	writeHsp, err := os.Create("writeHSP", os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
	defer writeHsp.Close()
				for i:= range hspSeq {
					write, err := writeHsp.write(>hspSeq[i].seqID\nhspseqSeq[i]\n)
				}
}

	}
}
