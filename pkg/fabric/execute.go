package fabric

import (
	"context"
	"fmt"
	"os"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func (m *Mixin) loadAction(ctx context.Context) (*Action, error) {
	var action Action
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute(ctx context.Context) error {
	action, err := m.loadAction(ctx)
	if err != nil {
		return err
	}
	var output string
	uuid := uuid.New()
	var outFilePath = "/cnab/app/" + uuid.String()

	action.Steps[0].Flags = append(action.Steps[0].Flags, builder.NewFlag("filePath", outFilePath))
	output, err = builder.ExecuteSingleStepAction(ctx, m.RuntimeConfig, action)
	if err != nil {
		return err
	}

	if _, err := os.Stat(outFilePath); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return err
	}

	fmt.Println("File exists")
	fmt.Println("ExecuteSingleStepAction OUTPUT", output)

	// executedStep := action.Steps[0]

	// outputData, err := os.ReadFile(outFilePath)

	// fmt.Println("OUTPUT EVIDENCE", string(outputData), len(outputData))

	// if len(executedStep.Instruction.Outputs) > 0 {

	// 	var instructionOutput = InstructionOutput{Name: executedStep.Instruction.Name, Outputs: executedStep.Instruction.Outputs}

	// 	//read from file

	// 	builder.ProcessJsonPathOutputs(ctx, m.RuntimeConfig, instructionOutput, string(outputData))

	// }

	return err
}

type InstructionOutput struct {
	Name    string   `yaml:"name"`
	Outputs []Output `yaml:"outputs,omitempty"`
}

func (s InstructionOutput) GetOutputs() []builder.Output {
	//	Go doesn't have generics, nothing to see here...
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}
