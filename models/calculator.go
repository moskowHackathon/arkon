package models

import(
	"math"
	"fmt"
	"math/rand"
	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"time"
)

const (
	_  = iota
	YES
	NO
	NA
)

/*c := bayesian.NewClassifier("Team1", "Team2")
	c.Learn([]string{"catalog", "price", "images"}, "Team1")
	c.Learn([]string{"cart", "checkout", "shipping"}, "Team2")*/

type Case struct {
	label string
	team int
	features []int
}

type SessionCalc struct {
	next int
	cases []Case
	answers []int
	isCompleted bool
}

func newSessionCalc(cases []Case) *SessionCalc {
	lf := len(cases[0].features)
	f := make([]int, lf)
	return  &SessionCalc{
		cases:cases,
		answers: f,
	}
}

//todo
// если тупик
func (s *SessionCalc) GetNextToCheck() (bool, int){
	beego.Error("START get next")
	beego.Error("cases: ", s.cases)
	beego.Error("answers: ", s.answers)

	if (!s.isLastAnswered()) {
		return true, s.next
	}
	if s.isCompleted {
		return false, 0
	}

	lf := len(s.cases[0].features)

	beego.Error("cases next", s.cases)

	variant := make([]int, lf)
	variantTrue := make([]int, lf)
	variantFalse := make([]int, lf)
	for _, icase :=range s.cases {
		for i, val := range icase.features {
			if s.answers[i]>0 {
				continue
			}
			if val==YES {
				variantTrue[i]+=1
			}
			if val==NO {
				variantFalse[i]+=1
			}
		}
	}
	l:=int(len(s.cases)/2)
	for i,_ := range variantTrue {
		variant[i] = abs(l-variantTrue[i])+abs(l-variantFalse[i])
	}

	fmt.Println(variantTrue)
	fmt.Println(variantFalse)
	fmt.Println(variant)

	res := getMinIndex(variant, s.answers)
	if res<0 {
		return false, 0
	}
	beego.Error("found: ", res)
	s.next = res

	return true, res
}

func (s *SessionCalc) ApplyAnswer(ans int) {
	beego.Error("apply: ", ans)
	s.answers[s.next] = ans
	if ans==NA  {
		return
	}
	newCases := []Case{}
	for i,icase:=range s.cases {
		if icase.features[s.next]==NA || icase.features[s.next]==ans {
			newCases = append(newCases, s.cases[i])
		}
	}
	s.cases = newCases
	beego.Error("new cases len: ", len(s.cases))
	beego.Error("new cases: ", s.cases)
}

func (s *SessionCalc) CheckStatus() (isFinish bool, result int){
	if s.isOneTeam() {
		s.isCompleted = true
		return true, s.cases[0].team
	}
	if !s.isCanContinue(){
		beego.Error("I can not continue")
		s.isCompleted = true
		return true, -1
	}
	return false, -1
}

func (s *SessionCalc) isCanContinue() bool{
	isCan := false
	for _, icase :=range s.cases {
		v := NA
		for i, val := range icase.features {
			if s.answers[i]>0 {
				continue
			}
			if val==NA {
				continue
			}
			if v==NA {
				v = val
				continue
			}
			if v!=val {
				isCan = true
			}

		}
		if isCan {
			break
		}
	}
	return isCan
}

func (s *SessionCalc) isOneTeam() bool{
	team:=s.cases[0].team
	for _,icase:=range s.cases {
		if icase.team!=team {
			return false
		}
	}
	return true
}

func (s *SessionCalc) isLastAnswered() bool{
	return s.answers[s.next]>0
}

func getMinIndex(arr []int, used []int)int{
	min := 0
	ret:= -1
	var arMins []int

	for i, value := range arr {
		if used[i]>0 {
			continue
		}
		if ret<0 || value <= min {
			min = value
			ret = i
			arMins = append(arMins, ret)
		}
	}
	if len(arMins)<=1 {
		return ret
	}
	beego.Error("arMins:", arMins)
	rand.Seed(time.Now().Unix())
	k:=rand.Int()
	beego.Error(k)
	n := k % len(arMins)
	beego.Error("n:", n)
	return arMins[n]
	//return ret
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}
