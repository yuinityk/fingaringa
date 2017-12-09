package main

import (
 "fmt"
 "math/rand"
 "math"
 "sort"
)


// Genom   : individual in group
// geneSeq : gene sequence of the individual
// fitness : fitness of the individual (=opt(g_i)), desired to be SMALL
type Genom struct {
  geneSeq [][]int
  fitness float64
}

func (g Genom) getFitness() float64 {
  return g.fitness
}

func (g Genom) getGene() [][]int {
  return g.geneSeq
}

func (g *Genom) setFitness(fitness float64) {
  g.fitness = fitness
}

func (g *Genom) setGene(geneSeq [][]int) {
  g.geneSeq = geneSeq
  g.evaluate()
}


type GenomList []Genom

func (gl GenomList) min() float64 {
  var k float64  = 1e7
  for i:=0;i<len(gl);i++{
    if k > gl[i].getFitness(){
      k = gl[i].getFitness()
    }
  }
  return k
}

// for sorting GenomList
// begin
func (gl GenomList) Len() int {
  return len(gl)
}

func (gl GenomList) Swap(i, j int) {
  gl[i], gl[j] = gl[j], gl[i]
}

func (gl GenomList) Less(i, j int) bool {
  return gl[i].fitness < gl[j].fitness
}
// end
// for sorting GenomList

func (gl GenomList) shuffle() {
  n := len(gl)
  for i := n-1; i>=0; i-- {
    j := rand.Intn(i+1)
    gl[i], gl[j] = gl[j], gl[i]
  }
}


func chordToPos(chord []int) (candidatePositions [][]int) {
  // Input
  //    chord              : tone list

  // Return
  //    candidatePositions : possible combination of positions of each tone
  
  for i := 0; i < len(chord); i++ {
    candidatePositions = append(candidatePositions,pos[chord[i]])
  }
  return
}

func randomChooseChordPos(candidatePositions [][]int) (positions []int) {
  // Input
  //    candidatePositions    : possible combination of positions of each tone

  // Return
  //    positions             : positions of each tone of the chord

  var ind int = 0
  for i := 0; i < len(candidatePositions); i++ {
    ind = rand.Intn(len(candidatePositions[i]))
    positions = append(positions,candidatePositions[i][ind])
  }
  return
}

func positionToX(position int) float64 {
  // Input
  //    position : position of tone on the fingerboard

  // Return
  //    float64  : distance between the position and the nat

  return 650*(1-1/(math.Pow(2,float64(position)/12)))
}

func createGenom(melody [][]int) (genom Genom) {
  // Input
  //    melody   : list of chords consisting of tone(s)

  // Return
  //　　genom    : randomly generated positions of each chords

  var positions [][]int
  for i := 0; i<len(melody); i++ {
    var pos []int = randomChooseChordPos(chordToPos(melody[i]))
    positions = append(positions,pos)
  }
  genom.setGene(positions)
  return
}

func evalChord(chord []int, melodyIndex int) (ret float64){
  for i := 0; i < len(chord); i++ {
    ret = 0
  }
  return
}

func (g *Genom)evaluate() {

  gene := g.getGene()
  var sum float64 = 0
  for d := 1; d < 3 ; d++ {
    for i := 0; i < len(gene)-d; i++ {
      g1 := gene[i]
      g2 := gene[i+d]
      for j := 0; j < len(g1); j++ {
        for k := 0; k < len(g2); k++ {
          if g1[j]*g2[k] != 0{
            sum += math.Abs(positionToX(g1[j])-positionToX(g2[k]))
          }
        }
      }
    }
  }
  // sum /= 2 
  // for i := 0; i < len(gene); i++ {
  //   sum += evalChord(gene[i],i)
  // }
  g.setFitness(sum)
}


func selectGenom(gl GenomList, eliteNum int) (gl_Elite GenomList) {
  // Input
  //    gl        : genome from which elite genome are chosen
  //    eliteNum  : number of genome to be chosen as elite genom

  // Return
  //    gl_Elite : elite genome

  gl_Elite = append(gl_Elite,gl...)
  sort.Sort(gl_Elite)
  return gl_Elite[:eliteNum]
}

func crossover(genom1, genom2 Genom) (gl_Ret GenomList) {
  // Input
  //    genom1, genom2 : genome to be crossed over
 
  // Return
  //    gl_Ret         : crossed over genome

  cross_1 := rand.Intn(GENOM_LENGTH)
  cross_2 := rand.Intn(GENOM_LENGTH-cross_1)+cross_1
  gene_1 := genom1.getGene()
  gene_2 := genom2.getGene()

  var progeny_1 [][]int
  var progeny_2 [][]int
  progeny_1 = append(progeny_1, gene_1[:cross_1]...)
  progeny_2 = append(progeny_2, gene_2[:cross_1]...)
  progeny_1 = append(progeny_1,gene_2[cross_1:cross_2]...)
  progeny_2 = append(progeny_2,gene_1[cross_1:cross_2]...)
  progeny_1 = append(progeny_1,gene_1[cross_2:]...)
  progeny_2 = append(progeny_2,gene_2[cross_2:]...)

  progenom1 := Genom{progeny_1,0}
  progenom2 := Genom{progeny_2,0}
  progenom1.evaluate()
  progenom2.evaluate()
  gl_Ret = append(gl_Ret,progenom1)
  gl_Ret = append(gl_Ret,progenom2)
  return gl_Ret 
}

func nextGeneCreate(gl, gl_Elite, gl_Progeny GenomList) (gl_Next GenomList) {
  // Input
  //    gl                   : genome at current generation
  //    gl_Elite, gl_Progeny : genome to be added to gl

  // Return
  //    gl_Next              : genome at next generation

  gl_Next = append(gl_Next, gl...)
  gl_Next.shuffle()
  gl_Next = gl_Next[:len(gl_Next)-len(gl_Elite)-len(gl_Progeny)]
  gl_Next = append(gl_Next, gl_Elite...)
  gl_Next = append(gl_Next, gl_Progeny...)
  return gl_Next
}

func mutation(gl GenomList, indivMutateRate float64, geneMutateRate float64) (gl_Ret GenomList) {
  // Input
  //    gl              : genome
  //    indivMutateRate : probability of mutation for each genom (individual)
  //    geneMutateRate  : probability of mutation for each gene of a gene sequence

  // Return
  //    gl_Ret          : mutated genome

  for i := 0; i < len(gl); i++ {
    genom := gl[i]
    if indivMutateRate > float64(rand.Intn(100))/100 {
      var newGenom [][]int
      for j := 0; j < len(genom.getGene()); j++ {
        if geneMutateRate > float64(rand.Intn(100))/100{
          p := randomChooseChordPos(chordToPos(input_melody[j]))
          newGenom = append(newGenom, p)
        } else {
          p := genom.getGene()[j]
          newGenom = append(newGenom, p)
        }
      }
      genom.setGene(newGenom)
    }
    gl_Ret = append(gl_Ret,genom)
  }
  return gl_Ret
}

// constants 
const INF = 1e7
const MAX_GENOM_LIST = 200       // number of genome at each generation
const SELECT_GENOM = 30          // number of elite genome
const INDIVIDUAL_MUTATION = 0.3  // probability of mutation for each individual
const GENOM_MUTATION = 0.4       // probability of mutation for each gene
const MAX_GENERATION = 2000       // max number of iteration for generation

// positions for each tone; one can play tones with different strings, different positions
var pos = [][]int{{0},{1},{2},{3},{4},{0,5},{1,6},{2,7},{3,8},{4,9},{0,5,10},{1,6,11},{2,7,12},{3,8,13},{4,9,14},{0,5,10},{1,6,11},{2,7,12},{3,8,13},{0,4,9,14},{1,5,10},{2,6,11},{3,7,12},{4,8,13},{0,5,9,14},{1,6,10,15},{2,7,11,16},{3,8,12,17},{4,9,13},{5,10,14},{6,11,15},{7,12,16},{8,13,17},{9,14,18},{10,15},{11,16},{12,17},{13,18},{14},{15},{16},{17},{18},{19},{20}}

// melody from Fantasy on Themes from La Traviata by J.Arcas
// https://www.youtube.com/watch?v=TDa5mUuSkbY
var input_melody = [][]int{{34},{26,22,17},{26,22,17},{33},{26,22,17},{31,26,22,17},{33,27,24,5},{17},{19},{21},{22},{24},{26},{27},{29},{31},{33},{34}}

var GENOM_LENGTH = len(input_melody)

func main() {

  var gl_Current GenomList
  for i:= 0; i<MAX_GENOM_LIST; i++ {
    gl_Current = append(gl_Current,createGenom(input_melody))
  }

  for count:=0; count< MAX_GENERATION; count++ {
    
    fmt.Println(count,gl_Current.min())
    gl_Elite := selectGenom(gl_Current, SELECT_GENOM)
    var gl_Progeny GenomList
    for i:=1; i<SELECT_GENOM; i++ {
      gl_Progeny = append(gl_Progeny,crossover(gl_Elite[i-1],gl_Elite[i])...)
    }
    gl_Next := nextGeneCreate(gl_Current, gl_Elite, gl_Progeny)
    gl_Next = mutation(gl_Next, INDIVIDUAL_MUTATION, GENOM_MUTATION)
    gl_Current = gl_Next
  }

  sort.Sort(gl_Current)
  for i:=0;i<len(gl_Current);i++{
    //fmt.Println(gl_Current[i].getFitness(),gl_Current[i].getGene())
    i=i
  }
  fmt.Println(gl_Current[0].getFitness(),gl_Current[0].getGene())

}
