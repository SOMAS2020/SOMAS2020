package roles

import (
	"testing"

	"reflect" // Used to compare two maps

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Failing test
func TestPresident(t *testing.T) {
	t.Error()
}

func mapSum(in_map map[int]int) int {
	sum := 0
	for _, val := range in_map {
		sum += val
	}
	return sum
}

type allocInput struct {
	resourceRequests   map[int]int
	commonPoolResource int
}

func TestAllocationRequests(t *testing.T) {
	cases := []struct {
		name  string
		input allocInput
		reply map[int]int
		want  error
	}{
		{
			name:  "Too many resources",
			input: allocInput{resourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}, commonPoolResource: 100},
			reply: map[int]int{1: 3, 2: 7, 3: 10, 4: 14, 5: 17, 6: 21},
			want:  nil,
		},
		{
			name:  "Enough resources",
			input: allocInput{resourceRequests: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30}, commonPoolResource: 150},
			reply: map[int]int{1: 5, 2: 10, 3: 15, 4: 20, 5: 25, 6: 30},
			want:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// mClient := &mockClientEcho{
			// 	id:   shared.Team1,
			// 	echo: tc.reply,
			// }
			// server := &SOMASServer{
			// 	clientMap: map[shared.ClientID]baseclient.Client{
			// 		shared.Team1: mClient,
			// 	},
			// }

			president := &basePresident{
				id:     int(shared.Team1),
				budget: 50,
				//resourceRequests:
			}

			val, _ := president.evaluateAllocationRequests(tc.input.resourceRequests, tc.input.commonPoolResource)
			if tc.want == nil {
				if !reflect.DeepEqual(tc.reply, val) {
					t.Errorf("%v - Failed. Expected %v, got %v", tc.name, tc.reply, val)

				}

			}
			//testutils.CompareTestErrors(tc.want, got, t)
			// president.resourceAllocation

			// Not sure how to test considering this class can be overridable?
			// Could check the sum is less than the given common pool
			// testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}
