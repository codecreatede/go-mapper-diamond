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
)




func main() {

    readPacBio := flag.String("pacbiohifi", "path to the fastq file", "file")
	readFasta := flag.String("proteinfasta", "path to the genome fasta", "file")

	 flag.Parse()

	// pacbio Structs
	type mapperID struct {
		id string
	}
	type mapperSeq struct {
		seq string
	}

	// pacbio opening and converting into structs
	pacbioOpen, err := os.Open(*readPacbio)
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

	 // pacbio fasta writing
			writeFasta, err := os.Create("writeFasta", os.O_CREATE|os.O_WRONLY, 0644)
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

    // running the diamond mapper
		cmd := exec.Command("diamond", "makedb", "--in", "*readFasta", "-d", "reference")
	  err := cmd.Run()
		if err != nil {
					log.Fatal("blast database not formatted %s\n",err)
				}
	  out, err := exec.Command("diamond", "makedb", "-in", "*readFasta", "-d", "reference").Output()
		if err != nil {
					log.Fatal("command failed with %s\n", err)
				}
		fmt.Println(out)

		cmd := exec.Command("diamond", "blastx", "-d", "reference", "-q", "*readFasta", "-o", "aligned.tsv")
		err := cmd.Run()
		if err != nil {
					log.Fatal("blastx didnt ran successfully %s\n", err)
				}
		out, err := exec.Command("diamond", "blastx", "-d", "reference", "-q", "*readsFasta", "-o", "aligned.tsv").Output()
	  if err != nil {
					log.Fatal("command failed with the error %s\n", err)
				}
		fmt.Println(out)
		 }

	// extracting the hsp struct for the sequence extraction

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
    writeFastaOpen , err := os.Open("writeFasta")
			if err != nil {
				log.Fatal(err)
			}
		writeFastRead := bufio.Newscanner(writeFastaOpen)
		seqIDWrite := []nucleotideIDStruct{}
		seqSeqWrite := []nucleotideSeqstruct{}
			for writeFastaRead.Scan() {
				line := writeFastaRead.Text()
				if strings.HasPrefix(string(line), ">") {
					seqIDWrite = append(seqIDWrite, nucleotideIDstruct{
						id : line
					})
				if strings.HasPrefix(string(line), "A") || strings.HasPrefix(string(line), "T") || strings.HasPrefix(string(line), "G") || strings.HasPrefix(string(line), "C") {
						seqSeqWrite = append(seqSeqWrite, nucleotiSeqStruct{
							seq : line
						})
					}
				}

 // extracting the hsp sequences

    hspSeq := []alignedFetch{}
	for i := range holdAlignment {
					for j := range seqIDWrite {
				if seqIDwrite[i] == holdAlignment.refseq[i] {
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
