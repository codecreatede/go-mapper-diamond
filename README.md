# go-mapper-diamond 

- a dimaond mapping of the pacbiohifi reads/fasta to the protein databases. 
- It just needs the pacbiohifi reads and the protein database and it aligns them and extracts the regions and store them as structs and interfaces.  
- diamond mapper can be found at [diamond mapper](https://github.com/bbuchfink/diamond)
- for hint generation and tag generation for AUGUSTUS and BRAKER and the sequence extraction for generating the environmental tags for the sequencing. 

```
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go -h
This is golang application for generating the hints from the protein alignment to pacbiohifi reads

Usage:
  flag [command]

Available Commands:
  align
  completion       Generate the autocompletion script for the specified shell
  help             Help about any command
  hspalignment
  pacbio
  proteinalignment
  seqHsp
  upStreamHSP

Flags:
  -h, --help   help for flag

Use "flag [command] --help" for more information about a command.
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go align -h
this option aligns and analyze all together and it requires only the proteins and the reference pacbio or the other fasta file

Usage:
  flag align and analyze [flags]

Flags:
  -h, --help                  help for align
  -p, --pacbio string         read-protein alignment (default "reads for protein alignment")
  -f, --pacbiofolder string   pacbio conversion (default "folder containing the bam files and the pbi files")
  -P, --protein string        protein datasets for the alignment (default "protein datasets")
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go hspalignment -h
Analyzes the hsp from the diamond read to protein alignment

Usage:
  flag hspalignment [flags]

Flags:
  -a, --alignmentfile string    alignment (default "alignment file to be analyzed")
  -h, --help                    help for hspalignment
  -p, --referencefasta string   pacbio file (default "pacbio reads file")
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go pacbio -h
converts the pacbio reads to the fasta format for the alignment and annotations, provide the folder path

Usage:
  flag pacbio [flags]

Flags:
  -h, --help                  help for pacbio
  -f, --pacbiofolder string   pacbio conversion (default "folder containing the bam files and the pbi files")
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go proteinalignment -h
aligns the pacbio reads or the fasta reads to the proteins

Usage:
  flag proteinalignment [flags]

Flags:
  -h, --help             help for proteinalignment
  -p, --pacbio string    read-protein alignment (default "reads for protein alignment")
  -P, --protein string   protein datasets for the alignment (default "protein datasets")
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go seqHsp -h                                                
Analyzes the hsp from the diamond read to protein alignment

Usage:
  flag seqHsp [flags]

Flags:
  -a, --alignmentfile string    alignment (default "alignment file to be analyzed")
  -h, --help                    help for seqHsp
  -p, --referencefasta string   pacbio file (default "pacbio reads file")
gauavsablok@gauravsablok ~/Desktop/codecreatede/golang/gomapper-diamond ±main⚡ » \
go run main.go upStreamHSP -h
specific for the genome alignment regions upstream and the downstream of the alignments

Usage:
  flag upStreamHSP [flags]

Flags:
  -a, --alignmentfile string             alignment (default "alignment file to be analyzed")
  -d, --downstream of the hsp tags int   downstream tags (default 5)
  -h, --help                             help for upStreamHSP
  -p, --referencefasta string            pacbio file (default "pacbio reads file")
  -u, --upstream of the hsp tags int     upstream tags (default 4)
```

- the usage of the corresponding each section is given below: 
```
go run main.go seqHsp -a ./samplefiles/matches.tsv -p ./samplefiles/fastafile.fasta 
go run main.go upStreamHSP -a ./samplefiles/matches.tsv -p ./samplefiles/fastafile.fasta -u 10 -d 10
go run main.go alignment -a matches.tsv -p ./samplefiles/pacbioreads.fasta
go run main.go analyze  -a ./samplefiles/matches.tsv -P ./samplefiles/fastafile.fasta
```

- in case of the binary use the following should be done 
```
./gomapperdiamod -h
./gomapperdiamond seqHsp -a ./samplefiles/matches.tsv -p ./samplefiles/fastafile.fasta
./gomappeddiamond upStreamHSP -a ./samplefiles/matches.tsv -p ./samplefiles/fastafile.fasta -u 10 -d 10
./gomapperdiamond alignment -a matches.tsv -p ./samplefiles/pacbioreads.fasta
./gomapperdiamond analyze  -a ./samplefiles/matches.tsv -P ./samplefiles/fastafile.fasta
```

Gaurav Sablok
