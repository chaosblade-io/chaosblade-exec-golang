package action

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/chaosblade-io/chaosblade-exec-golang/chaos/model"
)

type object struct {
	Index  int
	Name   string
	Fields map[string]string
}

func TestModifyAction_Execute(t *testing.T) {
	obj := object{
		Index:  0,
		Name:   "object",
		Fields: map[string]string{"a": "11", "b": "22"},
	}
	marshal, err := json.Marshal(obj)
	log.Println(string(marshal), err)
	type fields struct {
		OriginalValue   interface{}
		ExperimentFlags []model.ExperimentFlag
	}
	type args struct {
		ctx  context.Context
		rule *model.Experiment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "pointer type",
			fields: fields{
				OriginalValue: &obj,
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "{\"Index\":0,\"Name\":\"pointer-modify\"}",
					},
				},
			}},
			want: &object{
				Index:  0,
				Name:   "pointer-modify",
				Fields: map[string]string{"a": "11", "b": "22"},
			},
			wantErr: false,
		},
		{
			name: "struct type",
			fields: fields{
				OriginalValue: obj,
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "{\"Index\":0,\"Name\":\"struct-modify\"}",
					},
				},
			}},
			want: object{
				Index:  0,
				Name:   "struct-modify",
				Fields: nil,
			},
			wantErr: false,
		},
		{
			name: "int type",
			fields: fields{
				OriginalValue: 100,
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "200",
					},
				},
			}},
			want:    200,
			wantErr: false,
		},
		{
			name: "string type",
			fields: fields{
				OriginalValue: "bob",
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "caspar",
					},
				},
			}},
			want:    "caspar",
			wantErr: false,
		},
		{
			name: "map type",
			fields: fields{
				OriginalValue: map[string]int{"k1": 1, "k2": 2},
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "{\"k2\":22,\"k3\":3}",
					},
				},
			}},
			want:    map[string]int{"k2": 22, "k3": 3},
			wantErr: false,
		},
		{
			name: "slice type",
			fields: fields{
				OriginalValue: []string{"1","2","3"},
			},
			args: args{ctx: context.TODO(), rule: &model.Experiment{
				ExperimentRule: model.ExperimentRule{
					Target: "target",
					Action: "modify",
					Flags: map[string]string{
						"value": "[\"11\"]",
					},
				},
			}},
			want:    []string{"11"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &modifyAction{
				OriginalValue: tt.fields.OriginalValue,
			}
			response := m.Execute(tt.args.ctx, tt.args.rule)
			if !response.Success != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", response.Error, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(response.Result, tt.want) {
				t.Errorf("Execute() got = %v, want %v", response.Result, tt.want)
			}
		})
	}
}
