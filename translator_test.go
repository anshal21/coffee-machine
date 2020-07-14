package coffeemachine

// func Test_Translate(t *testing.T) {
// 	ir := `{
// 		"id": "some_ruleset",
// 		"predicates": {
// 		  "P1": "a > b",
// 		  "P2": "b > a",
// 		  "P3": "output_2 > d"
// 		},
// 		"rules": {
// 		  "R1": {
// 			"predicate": "P1",
// 			"postEvals": {
// 			  "output_1": {
// 				"expr": "a + b",
// 				"echo": false
// 			  },
// 			  "output_2": {
// 				"expr": "constant_string"
// 			  }
// 			}
// 		  },
// 		  "R2": {
// 			"predicate": "P2",
// 			"postEvals": {
// 			  "output_1": {
// 				"expr": "some_constant_response"
// 			  }
// 			}
// 		  },
// 		  "R3": {
// 			"predicate": "Predicate:P1 && Predicate:P2",
// 			"poseEvals": {
// 			  "soem_variable_name": {
// 				"expr": "a*b"
// 			  }
// 			}
// 		  }
// 		},

// 		"eelations": {
// 		  "edges": [
// 			{
// 			  "from": "R1",
// 			  "to": "R2",
// 			  "forwardOutput": true
// 			}
// 			]
// 		}

// 	  }`
// 	translator := translator{}
// 	translator.Translate([]byte(ir))
// }
