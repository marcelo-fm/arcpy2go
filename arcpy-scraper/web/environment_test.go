package web

import "testing"

func TestParseEnvironmentSettingPage(t *testing.T) {
	html := `
<html>
  <body>
    <header class="trailer-1">
      <h1>Current Workspace (Environment setting)</h1>
    </header>
    <div>
      <table>
        <tbody>
          <tr>
            <td>path</td>
            <td>The default location for geoprocessing tool input and output.</td>
          </tr>
        </tbody>
      </table>
    </div>
    <pre class="arcpyclass_msig">arcpy.env.workspace = path</pre>
  </body>
</html>`

	data, err := parseHTML("https://pro.arcgis.com/en/pro-app/3.4/tool-reference/environment-settings/current-workspace.htm", []byte(html))
	if err != nil {
		t.Fatalf("error parsing environment setting page: %v", err)
	}

	if data.FunctionName != "CurrentWorkspace" {
		t.Fatalf("expected FunctionName CurrentWorkspace, got %q", data.FunctionName)
	}
	if data.Command != "arcpy.env.workspace" {
		t.Fatalf("expected Command arcpy.env.workspace, got %q", data.Command)
	}
	if !data.IsAssignment {
		t.Fatal("expected IsAssignment to be true")
	}
	if data.AssignmentValueParam != "path" {
		t.Fatalf("expected AssignmentValueParam path, got %q", data.AssignmentValueParam)
	}
	if len(data.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(data.Parameters))
	}
	if data.Parameters[0].Name != "path" {
		t.Fatalf("expected parameter name path, got %q", data.Parameters[0].Name)
	}
	if data.Parameters[0].Comment == "" {
		t.Fatal("expected parameter comment to be populated")
	}
}
