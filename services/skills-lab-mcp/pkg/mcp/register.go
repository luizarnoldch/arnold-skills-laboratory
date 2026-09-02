package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/luizarnoldch/skills-lab-mcp/pkg/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Deps struct {
	Skills   *tools.Skills
	TestSets *tools.TestSets
	Evals    *tools.Evals
	Optimize *tools.Optimize
	Jobs     *tools.Jobs
}

func Register(server *mcp.Server, d Deps) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "skills_list",
		Description: "Lista todos los skills disponibles en el laboratorio con sus punteros a versiones current y test.",
	}, wrap(d.Skills.List))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "skill_get",
		Description: "Obtiene un skill por id o nombre con description/content current, versiones, e indica cuál es current y cuál la de pruebas.",
	}, wrap(d.Skills.Get))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "test_set_upload",
		Description: "Carga un set de pruebas (prompts) para un skill y aplica split automático train/validation.",
	}, wrap(d.TestSets.Upload))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "test_set_list",
		Description: "Lista los prompts de prueba de un skill, opcionalmente filtrados por split.",
	}, wrap(d.TestSets.List))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "test_set_update",
		Description: "Actualiza un prompt de prueba: query, should_trigger, split (train/validation) o runs.",
	}, wrap(d.TestSets.Update))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "baseline_eval_start",
		Description: "Encola evaluación baseline (train + validation) para una description específica o la de pruebas.",
	}, wrap(d.Evals.BaselineStart))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "trigger_eval_start",
		Description: "Encola evaluación de trigger queries para un split (train, validation o all).",
	}, wrap(d.Evals.TriggerStart))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "optimize_start",
		Description: "Encola optimización de description con parámetros configurables.",
	}, wrap(d.Optimize.Start))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "job_get",
		Description: "Consulta el estado y resultados de un job async (baseline_eval, trigger_eval u optimize).",
	}, wrap(d.Jobs.Get))
}

func wrap[In, Out any](fn func(context.Context, In) (Out, error)) func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args In) (*mcp.CallToolResult, any, error) {
		out, err := fn(ctx, args)
		if err != nil {
			return nil, nil, err
		}
		return jsonResult(out)
	}
}

func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, v, nil
}
